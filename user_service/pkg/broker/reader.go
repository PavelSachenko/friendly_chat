package broker

import (
	"context"
	"sync"
)

type BrokerReader interface {
	Read(ctx context.Context, wg *sync.WaitGroup) ([]byte, error)
	Close() error
}
