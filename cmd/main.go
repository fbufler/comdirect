package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/fbufler/comdirect/pkg/comdirect"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	// load from env
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	zugangsnummer := os.Getenv("ZUGANGSNUMMER")
	pin := os.Getenv("PIN")
	apiURL := os.Getenv("API_URL")
	tokenURL := os.Getenv("TOKEN_URL")
	revokeTokenURL := os.Getenv("REVOKE_TOKEN_URL")

	// create config
	config := comdirect.Config{
		APIURL:         apiURL,
		TokenURL:       tokenURL,
		RevokeTokenURL: revokeTokenURL,
		ClientID:       clientID,
		ClientSecret:   clientSecret,
		Zugangsnummer:  zugangsnummer,
		Pin:            pin,
	}

	fmt.Println(config)

	// create client
	client := comdirect.NewClient(config)

	// authenticate
	authResponse, err := client.NewToken()
	if err != nil {
		panic(err)
	}

	sessions, err := client.Sessions(authResponse)
	if err != nil {
		panic(err)
	}
	fmt.Println(sessions)

	validated, err := client.ValidateSession(authResponse, &sessions[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(validated)

	// Revoke token
	err = client.RevokeToken(authResponse)
	if err != nil {
		panic(err)
	}
}
