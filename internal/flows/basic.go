package flows

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"

	"github.com/fbufler/comdirect/config"
	"github.com/fbufler/comdirect/internal/cache"
	"github.com/fbufler/comdirect/pkg/comdirect"
)

func Bootstrap(cfg *config.Config) (*comdirect.Client, *comdirect.AuthToken, error) {
	config := comdirect.Config{
		APIURL:         cfg.Client.APIURL,
		TokenURL:       cfg.Client.TokenURL,
		RevokeTokenURL: cfg.Client.RevokeTokenURL,
		ClientID:       cfg.Client.ClientID,
		ClientSecret:   cfg.Client.ClientSecret,
		Zugangsnummer:  cfg.Client.Zugangsnummer,
		Pin:            cfg.Client.Pin,
	}

	client := comdirect.NewClient(config)
	cache := cache.NewCache(cfg.Cli.StoragePath, cfg.Cli.EncryptionKey)

	token, err := loadCache(cfg)
	if token != nil && !token.IsExpired() {
		return client, token, nil
	}

	if err != nil {
		slog.Warn(fmt.Sprintf("unable to load token from cache: %s", err))
		slog.Info("proceeding with authentication flow")
	}

	token, err = client.Authenticate(twoFaHandler)
	if err != nil {
		return client, token, err
	}

	if cfg.Cli.EnableCache {
		err = cache.Save(token)
		if err != nil {
			slog.Warn("unable to store token in cache")
		}
	}

	return client, token, err
}

func loadCache(cfg *config.Config) (*comdirect.AuthToken, error) {
	if cfg.Cli.EnableCache {
		c := cache.NewCache(cfg.Cli.StoragePath, cfg.Cli.EncryptionKey)
		token, err := c.Load()
		if err != nil {
			return nil, err
		}
		if token != nil {
			return token, nil
		}
	}
	return nil, nil
}

func twoFaHandler(tanHeader comdirect.TANHeader) error {
	slog.Info("Please verify the TAN")
	slog.Debug(fmt.Sprintf("TAN - id: %s - typ: %s", tanHeader.Id, tanHeader.Typ))

	slog.Info("Press enter to continue")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	slog.Info("Continuing")
	return nil
}
