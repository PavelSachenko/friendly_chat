package broker

import (
	"context"
)

type BrokerReader interface {
	Read(ctx context.Context, message chan []byte) error
	Close() error
}
