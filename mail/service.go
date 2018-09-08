package mail

import (
	"context"

	"github.com/sirupsen/logrus"
)

// Service is capable of sending mail messages using many providers.
// It contains the failover functionality.
type Service struct {
	providers []Provider

	l logrus.FieldLogger
}

// NewService constructs the Service capable of sending mail messages.
func NewService(providers []Provider, l logrus.FieldLogger) *Service {
	return &Service{
		providers: providers,
		l:         l,
	}
}

// Send sends the mail message using given mail providers.
func (s *Service) Send(ctx context.Context, args SendArgs) (sent bool) {
	for _, provider := range s.providers {
		err := provider.Send(ctx, args)
		if err == nil {
			sent = true
			break
		}
		s.l.Errorf("couldn't send mail, provider=%s, err: %s", provider.Name(), err)
	}
	return sent
}
