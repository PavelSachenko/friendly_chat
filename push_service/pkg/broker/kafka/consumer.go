package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pavel/push_service/config"
	"github.com/pavel/push_service/pkg/service/socket"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type kafkaBrokerReader struct {
	reader *kafka.Reader
}

func InitKafkaBrokerReader(cfg *config.Config) *kafkaBrokerReader {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:         []string{cfg.KafkaHost},
		GroupID:         "pusher",
		Topic:           "sms",
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         3 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: 1,
	}

	return &kafkaBrokerReader{
		reader: kafka.NewReader(kafkaConfig),
	}
}

type userMessage struct {
	UserIds []uint64 `json:"user_ids"`
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
		var userMessage userMessage
		err = json.Unmarshal(m.Value, &userMessage)
		if err != nil {
			log.Printf("error while decoding message: %s", err.Error())
			continue
		}

		hub.Broadcast <- socket.Broadcast{Broadcast: k.deleteUserIds(m.Value), UserIds: userMessage.UserIds}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Value))
	}
}

func (k kafkaBrokerReader) deleteUserIds(value []byte) []byte {
	var i interface{}
	if err := json.Unmarshal([]byte(value), &i); err != nil {
		log.Println(err)
	}
	if m, ok := i.(map[string]interface{}); ok {
		delete(m, "user_ids") // No problem if "foo" isn't in the map
	}

	output, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}

	return output
}
func (k kafkaBrokerReader) Close() error {
	return k.reader.Close()
}
