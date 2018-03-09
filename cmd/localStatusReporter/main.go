package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log"
	"github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/statusReporter"
	"github.com/iryonetwork/wwm/service/statusReporter/polling"
)

func main() {
	addr := ":443"
	metricsAddr := ":9090"
	keyFile := "/certs/localStatusReporter-key.pem"
	certFile := "/certs/localStatusReporter.pem"
	defaultTimeout := time.Duration(time.Second)
	countThreshold := 3
	interval := time.Duration(3 * time.Second)
	statusValidity := time.Duration(20 * time.Second)

	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "localStatusReporter").
		Logger()

	// create context with cancel func
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initialize status reporter
	r := statusReporter.New(logger)

	// add URL Status Polling components
	pollingCfg := &polling.Cfg{
		Interval:       &interval,
		CountThreshold: &countThreshold,
		StatusValidity: &statusValidity,
	}

	g := polling.New(polling.NewExternalURL("https://www.google.com", defaultTimeout), pollingCfg, logger)
	g.Start(ctx)
	r.AddComponent(statusReporter.External, "google", g)

	n := polling.New(polling.NewExternalURL("http://nna-leb.gov.lb", defaultTimeout), pollingCfg, logger)
	n.Start(ctx)
	r.AddComponent(statusReporter.External, "Lebanese National News Agency", n)

	localStorage := polling.New(polling.NewInternalURL("https://localStorage:4433/status", defaultTimeout), pollingCfg, logger)
	localStorage.Start(ctx)
	r.AddComponent(statusReporter.Local, "storage", localStorage)

	localAuth := polling.New(polling.NewInternalURL("https://localAuth:4433/status", defaultTimeout), pollingCfg, logger)
	localAuth.Start(ctx)
	r.AddComponent(statusReporter.Local, "auth", localAuth)

	cloudStorage := polling.New(polling.NewInternalURL("https://cloudStorage:4433/status", defaultTimeout), pollingCfg, logger)
	cloudStorage.Start(ctx)
	r.AddComponent(statusReporter.Cloud, "storage", cloudStorage)

	cloudAuth := polling.New(polling.NewInternalURL("https://cloudAuth:4433/status", defaultTimeout), pollingCfg, logger)
	cloudAuth.Start(ctx)
	r.AddComponent(statusReporter.Cloud, "auth", cloudAuth)

	// initialize metrics middleware
	m := api.NewMetrics("localStatusReporter", "")

	// setup handler
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(m.Middleware(log.APILogMiddleware(r.Handler("status"), logger)))

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start servers
	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	// waitGroup for all main go routines
	var wg sync.WaitGroup
	// create error channel
	errCh := make(chan error)

	// start serving metrics
	go func() {
		wg.Add(1)
		defer wg.Done()

		errCh <- metricsServer.ServePrometheusMetrics(context.Background(), metricsAddr, "status", logger)
	}()
	go func() {
		wg.Add(1)
		defer server.Close()
		defer wg.Done()

		localErrCh := make(chan error)
		go func() {
			logger.Info().Msgf("Starting status reporter server at %s", addr)
			localErrCh <- server.ListenAndServeTLS(certFile, keyFile)
		}()

		select {
		case err := <-localErrCh:
			errCh <- err
		case <-ctx.Done():
			//do nothing except deferred cleanup
		}
	}()

	// run cleanup when sigint or sigterm is received or error on starting server happened
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		wg.Add(1)
		defer cancelContext()
		defer wg.Done()

		select {
		case err := <-errCh:
			logger.Error().Err(err).Msg("failed to start server")
		case <-signalChan:
			logger.Error().Msg("received interrupt")
		}
	}()

	wg.Wait()
}
