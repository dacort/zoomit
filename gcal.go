package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Google oAuth parameters - despite the "Secret" name, this is not required to be ... secret
// Reference: https://developers.google.com/identity/protocols/oauth2/native-app
const clientID = "529787216143-e3i82arkvtgdcc1sp32p5jafgu63k79o.apps.googleusercontent.com"
const clientSecret = "D_hV0OftBH6aOmo9GCIfWwrD"

// Retrieve a token, saves the token, then returns the generated client.
func getToken(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(ctx, config)
		saveToken(tokFile, tok)
	}
	return tok
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	port, c := getPortAndWait()
	config.RedirectURL = "http://127.0.0.1:" + port
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	openBrowser(authURL)

	// Now we wait for the user to authenticate
	// TODO: Handle error state
	var authCode string

	select {
	case authCode = <-c:
	case <-time.After(time.Minute):
		log.Fatalln("Did not receive authentication callback in 1 minute, quitting...")
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func openBrowser(url string) {
	exec.Command("open", url).Run()
}

func getOrCreateConfigPath(file string) string {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Unable to determine OS config directory: %v", err)
	}

	configDir := path.Join(userConfigDir, "ZoomIT")

	// Ensure the config dir exists
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		log.Fatalf("Unable to create OS config directory: %v", err)
	}

	return path.Join(configDir, file)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(getOrCreateConfigPath(file))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(file string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(getOrCreateConfigPath(file), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func authorizeCalendar() *calendar.Service {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar.readonly",
		},
		Endpoint: google.Endpoint,
	}

	// Fetch the token and create a new calendar client
	ctx := context.Background()
	token := getToken(ctx, config)
	srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		log.Fatalf("Unable to create Google Calendar service: %v", err)
	}

	return srv
}
