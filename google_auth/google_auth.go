package google_auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-instaloader/config"
	"go-instaloader/utils/rlog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/http"
	"os"
)

var FileNotExistErr = errors.New("file not exist")

func GetHttpClient(jsonCredentialPath string) (client *http.Client, err error) {
	// read json credential file
	b, err := os.ReadFile(jsonCredentialPath)
	if err != nil {
		rlog.Error("Unable to read client secret file: %v", err)
		return nil, err
	}

	// get config
	config1, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		rlog.Error(fmt.Sprintf("Unable to parse client secret file to config: %v", err))
		return nil, err
	}

	// get token
	tokFile := config.Instance.TokenPath
	tok, tokErr := tokenFromFile(tokFile)
	if tokErr != nil {
		if errors.Is(tokErr, FileNotExistErr) {
			// if not exist, create a new one
			tok = getTokenFromWeb(config1)

			// save the token
			rlog.Info("Saving credential file to: %s\n", tokFile)
			f, err := os.OpenFile(tokFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				rlog.Fatal(fmt.Sprintf("Unable to cache oauth token: %v", err))
			}
			defer f.Close()
			json.NewEncoder(f).Encode(tok)
		} else {
			rlog.Error("get token from %s err: %s", tokFile, err.Error())
			return nil, tokErr
		}
	}

	// assign to client
	client = config1.Client(context.Background(), tok)
	return
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, FileNotExistErr
		} else {
			return nil, err
		}
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}
