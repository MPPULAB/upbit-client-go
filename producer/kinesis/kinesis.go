package kinesis

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"log"
)

type Producer interface {
	PutRecord(ctx context.Context, data []byte) error
}

type kinesisProducer struct {
	awsKinesis *kinesis.Client
}

func NewKinesisProducer() Producer {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	return &kinesisProducer{
		awsKinesis: kinesis.NewFromConfig(cfg),
	}
}

func (k *kinesisProducer) PutRecord(ctx context.Context, data []byte) error {
	_, err := k.awsKinesis.PutRecord(ctx, &kinesis.PutRecordInput{
		Data:         data,
		StreamName:   aws.String("test"),
		PartitionKey: aws.String("test"),
	})

	if err != nil {
		return err
	}

	return nil
}
