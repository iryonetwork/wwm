package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/log"
	"github.com/iryonetwork/wwm/metrics/api"
	metricsServer "github.com/iryonetwork/wwm/metrics/server"
	"github.com/iryonetwork/wwm/service/statusReporter"
	"github.com/iryonetwork/wwm/service/statusReporter/polling"
	"github.com/iryonetwork/wwm/service/tracing"
)

func main() {
	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "localStatusReporter").
		Logger()

	// create context with cancel func
	ctx, cancelContext := context.WithCancel(context.Background())
	defer cancelContext()

	// get config
	cfg, err := GetConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get config")
	}
	logger.Print(cfg)

	traceCloser := tracing.New("localStatusReporter", cfg.TracerAddr)
	defer traceCloser.Close()

	// initialize status reporter
	r := statusReporter.New(logger)

	// prepare default config for URL status polling components
	defaultPollingCfg := polling.Cfg{
		Interval:       &cfg.Interval,
		CountThreshold: &cfg.CountThreshold,
		StatusValidity: &cfg.StatusValidity,
	}

	// add URL status polling components
	for id, c := range cfg.Components.Local {
		pollingCfg, timeout := c.getConfig(cfg.DefaultTimeout, &defaultPollingCfg)

		var url polling.URLStatusEndpoint
		switch c.URLType {
		case polling.TypeInternalURL:
			url = polling.NewInternalURL(c.URL, timeout)
		default:
			url = polling.NewExternalURL(c.URL, timeout)
		}
		p := polling.New(url, &pollingCfg, logger)
		p.Start(ctx)
		r.AddComponent(statusReporter.Local, id, p)
	}
	for id, c := range cfg.Components.Cloud {
		pollingCfg, timeout := c.getConfig(cfg.DefaultTimeout, &defaultPollingCfg)

		var url polling.URLStatusEndpoint
		switch c.URLType {
		case polling.TypeInternalURL:
			url = polling.NewInternalURL(c.URL, timeout)
		default:
			url = polling.NewExternalURL(c.URL, timeout)
		}
		p := polling.New(url, &pollingCfg, logger)
		p.Start(ctx)
		r.AddComponent(statusReporter.Cloud, id, p)
	}
	for id, c := range cfg.Components.External {
		pollingCfg, timeout := c.getConfig(cfg.DefaultTimeout, &defaultPollingCfg)

		var url polling.URLStatusEndpoint
		switch c.URLType {
		case polling.TypeInternalURL:
			url = polling.NewInternalURL(c.URL, timeout)
		default:
			url = polling.NewExternalURL(c.URL, timeout)
		}
		p := polling.New(url, &pollingCfg, logger)
		p.Start(ctx)
		r.AddComponent(statusReporter.External, id, p)
	}

	// initialize metrics middleware
	m := api.NewMetrics("localStatusReporter", "")

	// setup handler
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	}).Handler(m.Middleware(log.APILogMiddleware(r.Handler("status"), logger)))

	// add tracer middleware
	handler = tracing.Middleware(handler)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPortHTTP),
		Handler: handler,
	}

	httpsServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPortHTTPS),
		Handler: handler,
	}

	// Start servers
	// create exit channel that is used to wait for all servers goroutines to exit orderly and carry the errors
	exitCh := make(chan error, 3)

	// start serving metrics
	go func() {
		exitCh <- metricsServer.ServePrometheusMetrics(ctx, fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.MetricsPort), cfg.MetricsNamespace, logger)
	}()
	go func() {
		defer httpServer.Close()
		defer httpsServer.Close()

		errHttpCh := make(chan error)
		go func() {
			logger.Info().Msgf("Starting HTTP status reporter server at %s", httpServer.Addr)
			errHttpCh <- httpServer.ListenAndServe()
		}()

		errHttpsCh := make(chan error)
		go func() {
			logger.Info().Msgf("Starting HTTPS status reporter server at %s", httpsServer.Addr)
			errHttpsCh <- httpsServer.ListenAndServeTLS(cfg.CertPath, cfg.KeyPath)
		}()

		for {
			select {
			case err := <-errHttpCh:
				exitCh <- err
				return
			case err := <-errHttpsCh:
				exitCh <- err
				return
			case <-ctx.Done():
				exitCh <- fmt.Errorf("StatusReporter server exiting because of cancelled context")
				//do nothing except deferred cleanup
				return
			}
		}
	}()

	// run cleanup when sigint or sigterm is received or error on starting server happened
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer cancelContext()

		for {
			select {
			case err := <-exitCh:
				logger.Info().Msg("exiting application because of exiting server goroutine")
				// pass error back to channel satisfy exit condition
				exitCh <- err
				return
			case <-signalChan:
				logger.Info().Msg("received interrupt")
				return
			}
		}
	}()

	<-ctx.Done()
	for i := 0; i < 3; i++ {
		err := <-exitCh
		if err != nil {
			logger.Debug().Err(err).Msg("gouroutine exit message")
		}
	}
}
