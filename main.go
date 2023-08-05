package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var UpbitWebsocketUrl = "wss://api.upbit.com/websocket/v1"

var interrupt = make(chan os.Signal, 1)

var EnvConfigs *envConfigs

type envConfigs struct {
	AccessKey string `mapstructure:"ACCESS_KEY"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func main() {
	signal.Notify(interrupt, os.Interrupt)

	EnvConfigs = loadEnvVariables()

	tokenString, err := createToken(EnvConfigs.SecretKey, EnvConfigs.AccessKey)
	if err != nil {
		log.Fatal(err)
	}

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+tokenString)

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	c, _, err := dialer.Dial(UpbitWebsocketUrl, headers)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer c.Close()

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(`[{"ticket":"test"},{"type":"ticker","codes":["KRW-BTC"]}]`))
			if err != nil {
				return
			}

			log.Println("tick:", t)

			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %s", message)

		case <-interrupt:
			log.Println("interrupt")

			// WebSocket 연결 종료
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}

			log.Println("successfully closed websocket connection")

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return
		}
	}
}

func createToken(secretKey string, accessKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"access_key": accessKey,
		"nonce":      uuid.New().String(),
	})
	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return ss, nil
}

func loadEnvVariables() (config *envConfigs) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	return
}
