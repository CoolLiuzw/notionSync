package todoapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	oauth "golang.org/x/oauth2"
)

// Client is used for HTTP requests to the Notion API.
type Client struct {
	httpClient *http.Client
}

func NewClient(clientID, clientSecret string) (*Client, error) {
	authEndpoint := oauth.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}

	authConfig := oauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     authEndpoint,
		RedirectURL:  "https://login.microsoftonline.com/common/oauth2/nativeclient",
		Scopes: []string{
			"offline_access",
			"Tasks.ReadWrite",
		},
	}

	ctx := context.WithValue(context.TODO(), oauth.HTTPClient, &http.Client{})

	token, err := GetSavedToken()
	if err != nil {
		return nil, err
	}

	return &Client{authConfig.Client(ctx, token)}, nil
}

func NewRequest(method, url string, header http.Header, param url.Values, body []byte) (*http.Request, error) {
	reader := bytes.NewReader(body)
	if method == http.MethodGet {
		url = urlPrefix + url
		if param != nil {
			url = url + "?" + param.Encode()
		}
	} else {
		url = urlPrefix + url
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	if header != nil {
		req.Header = header
	}

	return req, nil
}

func NewJSONRequest(method, url string, header http.Header, reqJo interface{}) (*http.Request, error) {
	body, err := json.Marshal(reqJo)
	if err != nil {
		return nil, err
	}

	req, err := NewRequest(method, url, header, nil, body)
	if err != nil {
		return nil, err
	}

	const headerContentType = "Content-Type"
	const contentTypeJSON = "application/json"
	req.Header.Set(headerContentType, contentTypeJSON)

	return req, nil
}
