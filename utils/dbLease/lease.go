package dbLease

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Config holds all the data required to establish and maintain a database
// connection.
type Config struct {
	VaultAddress string
	VaultToken   string
	VaultDBRole  string
}

// Lease holds and manages DB credentials
type lease struct {
	ctx    context.Context
	vault  *api.Client
	cfg    *Config
	logger zerolog.Logger
	lease  *api.Secret
}

// DBCreds contains credentials for connecting to database
type DBCreds struct {
	Username string
	Password string
}

// GetCreds fetches database credentials from Vault. Is also spawns a background
// worker to periodically renew the lease.
func GetCreds(ctx context.Context, cfg *Config, logger zerolog.Logger) (*DBCreds, error) {
	// create vault client
	vaultCfg := api.DefaultConfig()
	vaultCfg.Address = cfg.VaultAddress
	client, err := api.NewClient(vaultCfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create vault client")
	}

	client.SetToken(cfg.VaultToken)

	// fetch the secret
	sec, err := client.Logical().Read(fmt.Sprintf("database/creds/%s", cfg.VaultDBRole))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read credentials from vault")
	}

	l := &lease{
		cfg:    cfg,
		ctx:    ctx,
		vault:  client,
		logger: logger.With().Str("component", "utils/dbLease").Logger(),
		lease:  sec,
	}

	// spawn renew cycles
	go l.renewal()

	username, ok := sec.Data["username"].(string)
	if !ok {
		return nil, errors.New("Failed to convert security username to string")
	}
	password, ok := sec.Data["password"].(string)
	if !ok {
		return nil, errors.New("Failed to convert security password to string")
	}
	return &DBCreds{Username: username, Password: password}, nil
}

func (l *lease) renewal() {
	renewer, err := l.vault.NewRenewer(&api.RenewerInput{
		Secret: l.lease,
		Grace:  5 * time.Second,
	})
	if err != nil {
		l.logger.Warn().Err(err).Msg("Failed to setup renewer")
		return
	}

	go renewer.Renew()
	defer renewer.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-renewer.RenewCh():
			l.logger.Debug().Msg("Lease renewed")
		case err := <-renewer.DoneCh():
			l.logger.Warn().Err(err).Msg("Failed to renew lease")
			return
		}
	}
}
