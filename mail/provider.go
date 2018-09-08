package mail

import "context"

// Provider is an abstraction for mail provider.
type Provider interface {
	Name() string
	Send(ctx context.Context, message Message) error
}
