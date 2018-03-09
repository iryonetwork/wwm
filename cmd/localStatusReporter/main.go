package main

import (
	"context"
	"net/http"
	"os"
	"sync"
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
	r := statusReporter.New(logger.With().Str("component", "statusReporter").Logger())

	// add URL Status Polling components
	pollingCfg := &polling.Cfg{
		Interval:       &interval,
		CountThreshold: &countThreshold,
		StatusValidity: &statusValidity,
	}
	g := polling.New(
		polling.NewExternalURL("https://www.google.com", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "https://www.google.com").Logger(),
	)
	g.Start(ctx)
	r.AddComponent(statusReporter.External, "google", g)
	n := polling.New(
		polling.NewExternalURL("http://nna-leb.gov.lb", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "http://nna-leb.gov.lb").Logger(),
	)
	n.Start(ctx)
	r.AddComponent(statusReporter.External, "Lebanese National News Agency", n)
	localStorage := polling.New(
		polling.NewInternalURL("https://localStorage:4433/status", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "https://localStorage:4433/status").Logger(),
	)
	localStorage.Start(ctx)
	r.AddComponent(statusReporter.Local, "storage", localStorage)
	localAuth := polling.New(
		polling.NewInternalURL("https://localAuth:4433/status", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "https://localAuth:4433/status").Logger(),
	)
	localAuth.Start(ctx)
	r.AddComponent(statusReporter.Local, "auth", localAuth)
	cloudStorage := polling.New(
		polling.NewInternalURL("https://cloudStorage:4433/status", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "https://cloudStorage:4433/status").Logger(),
	)
	cloudStorage.Start(ctx)
	r.AddComponent(statusReporter.Cloud, "storage", cloudStorage)
	cloudAuth := polling.New(
		polling.NewInternalURL("https://cloudAuth:4433/status", defaultTimeout),
		pollingCfg,
		logger.With().Str("component", "URLPolling").Str("url", "https://cloudAuth:4433/status").Logger(),
	)
	cloudAuth.Start(ctx)
	r.AddComponent(statusReporter.Cloud, "auth", cloudAuth)

	// initialize metrics middleware
	m := api.NewMetrics("localStatusReporter", "")

	// setup handler
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(m.Middleware(log.APILogMiddleware(r.Handler("status"), logger.With().Str("component", "logMW").Logger())))

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// Start servers
	errCh := make(chan error)
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(errCh)
	}()

	go func() {
		wg.Add(1)
		defer server.Close()
		defer wg.Done()
		logger.Info().Msgf("Starting status reporter server at %s", addr)
		errCh <- server.ListenAndServeTLS(certFile, keyFile)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		errCh <- metricsServer.ServePrometheusMetrics(context.Background(), metricsAddr, "status", logger)
	}()

	for err := range errCh {
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}
}
