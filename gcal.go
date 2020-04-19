package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

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
	conf := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  "http://127.0.0.1:" + port,
		Scopes: []string{
			"https://www.googleapis.com/auth/calendar.readonly",
		},
		Endpoint: google.Endpoint,
	}
	authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	openBrowser(authURL)

	// Now we wait for the user to authenticate
	// TODO: Handle error state
	authCode := <-c
	// var authCode string
	// if _, err := fmt.Scan(&authCode);} err != nil {
	// 	log.Fatalf("Unable to read authorization code: %v", err)
	// }

	tok, err := conf.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func openBrowser(url string) {
	exec.Command("open", url).Run()
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func authorizeCalendar() *calendar.Service {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Fetch the token and create a new calendar client
	ctx := context.Background()
	token := getToken(ctx, config)
	srv, err := calendar.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))

	return srv
}
