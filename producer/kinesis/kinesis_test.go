package kinesis

import (
	"context"
	"testing"
)

func TestNewKinesisProducer(t *testing.T) {
	ctx := context.Background()

	producer := NewKinesisProducer()

	err := producer.PutRecord(ctx, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
}
