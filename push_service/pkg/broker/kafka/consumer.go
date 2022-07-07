package kafka

import (
	"context"
	"fmt"
	"github.com/pavel/push_service/config"
	"github.com/pavel/push_service/pkg/service/socket"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"time"
)

type kafkaBrokerReader struct {
	reader *kafka.Reader
}

func InitKafkaBrokerReader(cfg *config.Config) *kafkaBrokerReader {
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
		reader: kafka.NewReader(kafkaConfig),
	}
}

func (k kafkaBrokerReader) Read(ctx context.Context, hub *socket.Hub) error {
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
		args := strings.Split(string(m.Value), " ")
		hub.Broadcast <- socket.Broadcast{Broadcast: []byte(args[0]), Username: args[1]}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}

func (k kafkaBrokerReader) Close() error {
	return k.reader.Close()
}
