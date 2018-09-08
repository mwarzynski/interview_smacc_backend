package mail

import "context"

// SendArgs define required values for the message to be sent.
type SendArgs struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Provider is an abstraction for mail provider.
type Provider interface {
	Name() string
	Send(ctx context.Context, args SendArgs) error
}
