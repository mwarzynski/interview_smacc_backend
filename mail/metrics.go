package mail

import (
	"context"
	"fmt"
	"time"

	metrics "github.com/armon/go-metrics"
)

type providerWithMetrics struct {
	Provider

	m *metrics.Metrics
}

// ProviderWithMetrics wraps Provider with metrics functionality.
// Time of request execution as well as errors are metriced.
func ProviderWithMetrics(p Provider, m *metrics.Metrics) Provider {
	return &providerWithMetrics{
		Provider: p,
		m:        m,
	}
}

func (pm *providerWithMetrics) Send(ctx context.Context, message Message) error {
	defer pm.m.MeasureSince([]string{"mail", "provider", pm.Provider.Name()}, time.Now())
	err := pm.Provider.Send(ctx, message)
	if err != nil {
		pm.m.IncrCounter([]string{"errors", fmt.Sprintf("mail-provider-%s", pm.Provider.Name())}, 1)
	}
	return err
}
