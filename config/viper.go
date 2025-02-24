package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ClientConfig struct {
	APIURL         string `mapstructure:"api-url"`
	TokenURL       string `mapstructure:"token-url"`
	RevokeTokenURL string `mapstructure:"revoke-token-url"`
	ClientID       string `mapstructure:"client-id"`
	ClientSecret   string `mapstructure:"client-secret"`
	Zugangsnummer  string `mapstructure:"zugangsnummer"`
	Pin            string `mapstructure:"pin"`
}

type CliConfig struct {
	EnableCache   bool   `mapstructure:"enable-cache"`
	EncryptionKey string `mapstructure:"encryption-key"`
	StoragePath   string `mapstructure:"storage-path"`
}

type Config struct {
	Client  ClientConfig `mapstructure:"client"`
	Cli     CliConfig    `mapstructure:"cli"`
	Verbose bool         `mapstructure:"verbose"`
}

func setDefaults() {
	viper.SetDefault("client.api-url", "https://api.comdirect.de/api")
	viper.SetDefault("client.token-url", "https://api.comdirect.de/oauth/token")
	viper.SetDefault("client.revoke-token-url", "https://api.comdirect.de/oauth/revoke")
	viper.SetDefault("cli.enable-cache", false)
	viper.SetDefault("cli.storage-path", cliStoragePath())
}

func cliStoragePath() string {
	if cfg.Cli.StoragePath != "" {
		return cfg.Cli.StoragePath
	}
	if os.Getenv("OS") == "Windows_NT" {
		return windowsCliStoragePath()
	}
	return unixCliStoragePath()
}

func windowsCliStoragePath() string {
	return os.TempDir() + "\\token-cache"
}

func unixCliStoragePath() string {
	return os.TempDir() + "/token-cache"
}

var cfg Config = Config{}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	configPaths := configPaths()
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
	viper.AutomaticEnv()
	setDefaults()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Cli.EnableCache {
		if cfg.Cli.EncryptionKey == "" {
			panic("encryption key must be set if cache is enabled")
		}
		// AES encryption key must be 128, 192, or 256 bits long
		if len(cfg.Cli.EncryptionKey) != 16 && len(cfg.Cli.EncryptionKey) != 24 && len(cfg.Cli.EncryptionKey) != 32 {
			panic(fmt.Sprintf("encryption key must be 128, 192, or 256 bits long, got %d", len(cfg.Cli.EncryptionKey)))
		}
	}
}

func configPaths() []string {
	if os.Getenv("OS") == "Windows_NT" {
		return windowsConfigPaths()
	}
	return unixConfigPaths()
}

func windowsConfigPaths() []string {
	return []string{
		"C:\\ProgramData\\comdirect",
		"C:\\Users\\%USERNAME%\\AppData\\Roaming\\comdirect",
		".",
	}
}

func unixConfigPaths() []string {
	return []string{
		"/etc/comdirect",
		"$HOME/.comdirect",
		".",
	}
}
func Get() *Config {
	return &cfg
}
