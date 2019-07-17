// Copyright 2019 Ahmed Abouzied. All rights reserved.

// Package dwolla is a client library for Dwolla v2 rest api.
package dwolla

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Client represents a client for the dwolla REST API.
type Client struct {
	Client       *http.Client           //http client
	Env          string                 // either sandbox or production
	ClientID     string                 // Dwolla client ID
	ClientSecret string                 // Dwolla Client Secret
	authToken    string                 // Dowlla Auth token that expires in 1 hour
	rootURL      string                 // Root url of dwolla api. Differs according to Env
	Links        map[string]interface{} // Links to account resources
}

// CreateClient creates a new Dwolla Client
func CreateClient(env string, clientID string, clientSecret string) *Client {
	c := &Client{
		Client:       &http.Client{},
		Env:          env,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	switch env {
	case "sandbox":
		c.SetRootURL("https://api-sandbox.dwolla.com")
	case "production":
		c.SetRootURL("https://api.dwolla.com")
	default:
		c.SetRootURL("https://api-sandbox.dwolla.com")
	}
	return c
}

// SetRootURL sets the rootURL of the client to the given value
func (c *Client) SetRootURL(URL string) {
	c.rootURL = URL
}

// RootURL returns the root url of the clinet
func (c *Client) RootURL() string {
	return c.rootURL
}

// AuthToken returns the current auth token
func (c *Client) AuthToken() (string, error) {
	err := c.SetAccessToken()
	if err != nil {
		return "", errors.Wrap(err, "failed to refresh access token")
	}
	return c.authToken, nil
}

// SetAccessToken makes a request to dwolla to get an access token. Then sets this token into the current client.
func (c *Client) SetAccessToken() error {
	// Request access token
	req, err := http.NewRequest("POST", c.RootURL()+"/token", bytes.NewReader([]byte("grant_type=client_credentials")))
	if err != nil {
		return errors.Wrap(err, "error creating get token request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(c.ClientID+":"+c.ClientSecret))))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.Wrap(err, "error making request to dowlla api")
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "failed to get access token from dwolla api")
	}
	token, err := decodeAuthTokenResp(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to decode access token")
	}
	c.authToken = token
	return nil
}

// Root returns the resources avaliable by dwolla api
func (c *Client) Root() (map[string]interface{}, error) {
	// Request access token
	req, err := http.NewRequest("GET", c.RootURL(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating get root request")
	}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "error refreshing access token")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := c.Client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error making request to root endpoint")
	}
	defer res.Body.Close()
	resources := make(map[string]interface{})
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&resources)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json response")
	}
	return resources, nil
}

func decodeAuthTokenResp(r io.Reader) (string, error) {
	type authResponse struct {
		TokenType string `json:"token_type"`
		Token     string `json:"access_token"`
		ExpiresIn int    `json:"expires_in"`
	}
	d := json.NewDecoder(r)
	tokenResp := &authResponse{}
	err := d.Decode(tokenResp)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode json response")
	}
	return tokenResp.Token, nil
}
