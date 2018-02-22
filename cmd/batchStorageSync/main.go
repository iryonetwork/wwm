package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/rs/zerolog"

	"github.com/iryonetwork/wwm/gen/storage/client"
	storageSync "github.com/iryonetwork/wwm/sync/storage"
	"github.com/iryonetwork/wwm/sync/storage/batch"
)

type clientAuthInfoWriter struct {
	apiKey string
}

func (a *clientAuthInfoWriter) AuthenticateRequest(r runtime.ClientRequest, f strfmt.Registry) error {
	r.SetHeaderParam("Authorization", a.apiKey)
	return nil
}

func main() {
	lastSuccessfulRun, _ := strfmt.ParseDateTime("2018-01-18T15:22:46.123Z")

	// initialize logger
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "batchStorageSync").
		Logger()

	// initialize local storage API client
	localTransportCfg := &client.TransportConfig{
		Host:     "iryo.local",
		BasePath: "storage",
		Schemes:  []string{"https"},
	}
	localClient := client.NewHTTPClientWithConfig(strfmt.NewFormats(), localTransportCfg)

	// initialize cloud storage API client
	cloudTransportCfg := &client.TransportConfig{
		Host:     "iryo.cloud",
		BasePath: "storage",
		Schemes:  []string{"https"},
	}
	cloudClient := client.NewHTTPClientWithConfig(strfmt.NewFormats(), cloudTransportCfg)

	// initizalize mock request authenticator
	auth := &clientAuthInfoWriter{"SECRETSECRETSECRETSECRETSECRETSE"}

	// initialize handlers
	handlers := storageSync.NewHandlers(localClient.Operations, auth, cloudClient.Operations, auth, logger.With().Str("component", "sync/storage/handlers").Logger())

	// Create context with cancel func
	ctx, shutdown := context.WithCancel(context.Background())

	// Initialize batchStorageSync
	s := batch.New(handlers, logger)

	// Run sync
	errCh := make(chan error)
	go func() {
		err := s.Sync(ctx, time.Time(lastSuccessfulRun))
		errCh <- err
	}()

	// Run cleanup when sigint or sigterm is received
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil {
			logger.Fatal().Err(err).Msg("batch sync failed")
		}
		logger.Info().Msg("batch sync successfull")
		// Save lastSuccesfulRun
	case <-signalChan:
		logger.Info().Msg("stopping batch sync due to interrupt")
		shutdown()
	}
}
