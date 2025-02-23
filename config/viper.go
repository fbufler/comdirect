package config

import "github.com/spf13/viper"

type ClientConfig struct {
	APIURL         string `mapstructure:"api-url"`
	TokenURL       string `mapstructure:"token-url"`
	RevokeTokenURL string `mapstructure:"revoke-token-url"`
	ClientID       string `mapstructure:"client-id"`
	ClientSecret   string `mapstructure:"client-secret"`
	Zugangsnummer  string `mapstructure:"zugangsnummer"`
	Pin            string `mapstructure:"pin"`
}

func setDefaults() {
	viper.SetDefault("api-url", "https://api.comdirect.de/api")
	viper.SetDefault("token-url", "https://api.comdirect.de/oauth/token")
	viper.SetDefault("revoke-token-url", "https://api.comdirect.de/oauth/revoke")
}

var cfg ClientConfig = ClientConfig{}

func init() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.comdirect")
	viper.AddConfigPath("/etc/comdirect")
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
}

func Get() *ClientConfig {
	return &cfg
}
