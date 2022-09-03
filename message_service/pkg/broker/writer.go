package broker

import "context"

type BrokerWriter interface {
	Push(parent context.Context, key, value []byte) (err error)
}
