package kafka

import (
	"context"
	"fmt"
	"github.com/pavel/user_service/config"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

type kafkaBrokerReader struct {
	reader   *kafka.Reader
	messages chan string
}

func InitKafkaBrokerReader(cfg config.Config, c chan string) *kafkaBrokerReader {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:         []string{"localhost:9092"},
		GroupID:         "pusher",
		Topic:           "test",
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         3 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: 1,
	}

	return &kafkaBrokerReader{
		reader:   kafka.NewReader(kafkaConfig),
		messages: c,
	}
}

func (k kafkaBrokerReader) Read(ctx context.Context, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()
	for {
		m, err := k.reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("error while receiving message: %s", err.Error())
			continue
		}
		err = k.reader.CommitMessages(ctx, m)
		if err != nil {
			log.Printf("error while commiting message: %s", err.Error())
			continue
		}
		k.messages <- string(m.Value)
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}

func (k kafkaBrokerReader) Close() error {
	return k.reader.Close()
}
