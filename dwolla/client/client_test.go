// Copyright 2019 Ahmed Abouzied. All rights reserved.

// Package dwolla is a client library for Dwolla v2 rest api.
package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockToken = `
{
  "access_token": "SF8Vxx6H644lekdVKAAHFnqRCFy8WGqltzitpii6w2MVaZp1Nw",
  "token_type": "bearer",
  "expires_in": 3600
}
`
var mockRoot = `
{
  "_links": {
    "account": {
      "href": "https://api-sandbox.dwolla.com/accounts/ad5f2162-404a-4c4c-994e-6ab6c3a13254"
    },
    "customers": {
      "href": "https://api-sandbox.dwolla.com/customers"
    }
  }
}
`

func TestSetAuthToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockToken)
	}))
	defer ts.Close()

	mock := &Client{
		Env:          "Test",
		ClientID:     "123456789",
		ClientSecret: "123456789",
		rootURL:      "http://localhost:8080",
	}
	mock.SetRootURL(ts.URL)
	token, err := mock.AuthToken()
	if err != nil {
		t.Error(err)
	}
	t.Log("token : ", token)
}

func TestRoot(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockRoot)
	}))
	defer ts.Close()
	mock := &Client{
		Env:          "Test",
		ClientID:     "123456789",
		ClientSecret: "123456789",
		rootURL:      "http://localhost:8080",
	}
	mock.SetRootURL(ts.URL)
	_, err := mock.Root()
	if err != nil {
		t.Error(err)
	}
}
