package todoapi

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	oauth "golang.org/x/oauth2"
)

func GetToken(clientID, clientSecret string) (*oauth.Token, error) {
	const authURL = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize"
	const tokenURL = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	const redirectURL = "https://login.microsoftonline.com/common/oauth2/nativeclient"

	authConfig := oauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     oauth.Endpoint{AuthURL: authURL, TokenURL: tokenURL},
		RedirectURL:  redirectURL,
		Scopes: []string{
			"offline_access",
			"Tasks.ReadWrite",
		},
	}
	url := authConfig.AuthCodeURL("state", oauth.AccessTypeOffline)
	log.Println("go to the next link : ")
	log.Println(url)
	log.Println()
	log.Println("Put the code below:")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Println("Can`t get code!")
		return nil, err
	}

	log.Println(code)

	ctx := context.WithValue(context.TODO(), oauth.HTTPClient, &http.Client{})
	token, err := authConfig.Exchange(ctx, code)
	log.Printf("token: %v", token)

	return token, err
}

func GetSavedToken() (*oauth.Token, error) {
	savedRefreshToken, err := ioutil.ReadFile("token.txt")
	if err != nil {
		return nil, err
	}

	newToken := oauth.Token{
		RefreshToken: string(savedRefreshToken),
		TokenType:    "Bearer",
	}

	return &newToken, nil
}

func SaveToken(refreshToken []byte) error {
	file, err := os.Create("token.txt")
	if err != nil {
		return err
	}

	_, err = file.Write(refreshToken)
	if err != nil {
		return err
	}
	_ = file.Close()

	return nil
}
