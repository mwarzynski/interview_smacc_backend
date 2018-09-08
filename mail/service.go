package mail

import "context"

// Service is capable of sending mail messages using many providers.
// It contains the failover functionality.
type Service struct {
	providers []Provider
}

// NewService constructs
func NewService(providers []Provider) *Service {
	return &Service{
		providers: providers,
	}
}

// Send sends the mail message using given mail providers.
func (s *Service) Send(ctx context.Context, params SendArgs) error {
	return nil
}
