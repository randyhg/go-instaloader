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

//func GetHttpClient() (client *http.Client, err error) {
//	// read json credential file
//	b, err := os.ReadFile(config.Instance.CredentialPath)
//	if err != nil {
//		rlog.Errorf("unable to read client secret file: %v", err)
//		return nil, err
//	}
//
//	// get config
//	config1, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
//	if err != nil {
//		rlog.Errorf("unable to parse client secret file to config: %v", err)
//		return nil, err
//	}
//
//	// get token
//	tokFile := config.Instance.TokenPath
//	tok, tokErr := tokenFromFile(tokFile)
//	if tokErr != nil {
//		if errors.Is(tokErr, FileNotExistErr) {
//			// if not exist, create a new one
//			tok = getTokenFromWeb(config1)
//
//			// save the token
//			if err = saveTokenToFile(tokFile, tok); err != nil {
//				return nil, err
//			}
//		} else {
//			rlog.Errorf("get token from %s err: %s", tokFile, tokErr.Error())
//			return nil, tokErr
//		}
//	}
//
//	if tok.Expiry.Before(time.Now()) {
//		rlog.Info("token expired. refreshing the token....")
//		// do refresh token
//		refTok, err := refreshToken(tok, config1)
//		if err != nil {
//			rlog.Fatal(err)
//		}
//		// assign to client
//		client = config1.Client(context.Background(), refTok)
//		return client, nil
//	}
//
//	// assign to client
//	client = config1.Client(context.Background(), tok)
//	return
//}

func GetHttpClient() (client *http.Client, err error) {
	jsonKey, err := os.ReadFile(config.Instance.ServiceKeyPath)
	if err != nil {
		rlog.Errorf("unable to read client secret file: %v", err)
		return nil, err
	}

	config1, err := google.JWTConfigFromJSON(jsonKey, sheets.SpreadsheetsScope)
	if err != nil {
		rlog.Errorf("unable to parse client secret file to config: %v", err)
		return nil, err
	}

	// assign to client
	client = config1.Client(context.Background())
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
		rlog.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		rlog.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func refreshToken(tok *oauth2.Token, config1 *oauth2.Config) (*oauth2.Token, error) {
	// Create a token source using the current token
	tokenSource := config1.TokenSource(context.Background(), tok)

	// Use the token source to get a new token
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	if err = saveTokenToFile(config.Instance.TokenPath, tok); err != nil {
		return nil, err
	}

	return newToken, nil
}

func saveTokenToFile(filePath string, token *oauth2.Token) error {
	rlog.Infof("saving credential file to: %s", filePath)
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		rlog.Error(err)
		return err
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}
