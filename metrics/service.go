package metrics

import (
	metrics "github.com/armon/go-metrics"
	"github.com/pkg/errors"
)

// New creates metrics service.
func New(host, port, prefix string) (*metrics.Metrics, error) {
	sink, err := metrics.NewStatsdSink(host + ":" + port)
	if err != nil {
		return nil, errors.Wrap(err, "NewStatsdSink")
	}

	conf := metrics.DefaultConfig(prefix)
	conf.EnableRuntimeMetrics = true
	conf.EnableHostname = false
	m, err := metrics.New(conf, sink)
	if err != nil {
		return nil, errors.Wrap(err, "New")
	}

	return m, nil
}

// NewDummy creates dummy metrics service for testing.
func NewDummy() *metrics.Metrics {
	m, _ := metrics.New(&metrics.Config{}, &metrics.BlackholeSink{})
	return m
}
