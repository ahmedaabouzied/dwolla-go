// Package account provides methods to use the dwolla api to manage master dwolla account.
package account

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/masspayment"
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
	body, err := json.Marshal(fundingResource)
	if err != nil {
		return errors.Wrap(err, "error marshalling funding resource into req body")
	}
	req, err := http.NewRequest("POST", c.RootURL()+"/funding-sources", bytes.NewReader(body))
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
		switch res.StatusCode {
		case 400:
			return errors.Wrap(errors.New("Bad Request"), "duplicate funding resource of validation error")
		case 403:
			return errors.Wrap(errors.New("Unauthorized"), "not authorized to create funding resource")
		default:
			return errors.New(res.Status)
		}
	}
	return nil
}

// ListFundingResources retrieves a list of funding sources that belong to an Account
func (a *Account) ListFundingResources(c *client.Client) ([]funding.Resource, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", a.Links["funding-sources"]["href"], nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to dwolla server")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &funding.ListResourcesResponse{}
		err = d.Decode(body)
		return body.Embeded["funding-sources"], nil
	case 403:
		return nil, errors.New("not authorized to list funding sources")
	case 404:
		return nil, errors.New("account not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// TODO : Add ListAndSearchTransfers method

// ListMassPayments retrieves an Account’s list of previously created mass payments
func (a *Account) ListMassPayments(c *client.Client) ([]masspayment.MassPayment, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", a.Links["self"]["href"]+"/mass-payments", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to dwolla server")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &masspayment.ListMassPaymentsResponse{}
		err = d.Decode(body)
		return body.Embedded["mass-payments"], nil
	case 403:
		return nil, errors.New("not authorized to list mass payments")
	case 404:
		return nil, errors.New("account not found")
	default:
		return nil, errors.New(res.Status)
	}
}
