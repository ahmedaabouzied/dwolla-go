// Copyright 2019 Ahmed Abouzied. All rights reserved.

// Package dwolla is a client library for Dwolla v2 rest api.
package dwolla

import (
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

func TestAuthToken(t *testing.T) {
	gotenv.Load(".env")
	client := CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	token, err := client.AuthToken()
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("No token was returned")
	}
	t.Log(token)
}

func TestRoot(t *testing.T) {
	client := CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	resources, err := client.Root()
	if err != nil {
		t.Error(err)
	}
	t.Log(resources["_links"])
}
