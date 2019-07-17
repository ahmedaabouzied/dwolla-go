package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/ahmedaabouzied/dwolla/client"
	"github.com/ahmedaabouzied/dwolla/funding"
	"github.com/pkg/errors"
)

// Account represents a Dwolla master account that was estabslished on dwolla.com
type Account struct {
	Links map[string]map[string]string `json:"_links"`
	ID    string                       `json:"id"`   // Dwolla account ID
	Name  string                       `json:"name"` // Dwolla account holder name
}

// RetrieveAccount returns the Dwolla master account
func RetrieveAccount(c *client.Client) (*Account, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", c.Links["account"]["href"], nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating get root request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	acc := &Account{}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(acc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse json response")
	}
	return acc, nil
}

// CreateFundingSource adds a funding resource to the master dwolla account
func (a *Account) CreateFundingSource(c *client.Client, fundingResource *funding.Resource) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("POST", c.RootURL()+"/funding-sources", nil)
	if err != nil {
		return errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Add("Content-Type", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		return errors.New(res.Status)
	}
	return nil
}
