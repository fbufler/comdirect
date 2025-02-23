package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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

	challengeID, err := client.ValidateSession(authResponse, sessions[0].Identifier)
	if err != nil {
		panic(err)
	}
	fmt.Println(challengeID)

	// sleep for 10 seconds
	fmt.Println("Please accept the challenge in the comdirect app")
	time.Sleep(10 * time.Second)

	newSession, err := client.ActivateSession(authResponse, sessions[0].Identifier, challengeID)
	if err != nil {
		panic(err)
	}
	fmt.Println(newSession)

	// get account balances TODO: Broken
	accountBalances, err := client.AccountBalance(authResponse)
	if err != nil {
		panic(err)
	}
	fmt.Println(accountBalances)

	// Revoke token
	err = client.RevokeToken(authResponse)
	if err != nil {
		panic(err)
	}
}
