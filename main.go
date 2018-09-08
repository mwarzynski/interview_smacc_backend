package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caarlos0/env"
	"github.com/mwarzynski/smacc-backend/metrics"
	"github.com/mwarzynski/smacc-backend/transport"
	"github.com/sirupsen/logrus"
)

func main() {
	// parse config
	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatalf("failed to parse config: %v\n", err)
	}

	// application context
	appCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// logger setup
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %v\n", err)
	}
	l := logrus.New()
	l.SetLevel(logLevel)

	// metrics
	metricsService := metrics.NewDummy()
	if config.MetricsEnabled {
		metricsService, err = metrics.New(config.MetricsHost, config.MetricsPort, config.MetricsPrefix)
	}
	if err != nil {
		l.WithField("tags", []string{"metrics"}).Fatalf("creating service: %v", err)
	}

	// default internal HTTP client
	defaultHTTPTimeout := time.Duration(10) * time.Second
	defaultInternalClient := http.Client{
		Timeout: defaultHTTPTimeout,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: (&net.Dialer{
				Timeout:   defaultHTTPTimeout,
				KeepAlive: defaultHTTPTimeout,
			}).DialContext,
			MaxIdleConns:          100,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			TLSHandshakeTimeout:   10 * time.Second,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// router
	r := transport.Init(
		appCtx,
		&defaultInternalClient,
		metricsService,
		l,
	)

	// setup server
	server := &http.Server{
		Addr:         config.Listen,
		Handler:      r,
		ReadTimeout:  time.Second * time.Duration(config.HTTPServerTimeoutSeconds),
		WriteTimeout: time.Second * time.Duration(config.HTTPServerTimeoutSeconds),
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	// graceful shutdown setup
	go func() {
		<-signals
		cancelCtx()

		if err := server.Shutdown(appCtx); err != nil {
			log.Fatalf("%v", err)
			l.WithField("tags", []string{"error", "shutdown"}).Infof("Shutdown error")
		}
		l.WithField("tags", []string{"stop"}).Infof("server stopped")
	}()

	// run server
	l.WithField("tags", []string{"start"}).Infof("started listening on %s...", config.Listen)
	server.ListenAndServe()
}
