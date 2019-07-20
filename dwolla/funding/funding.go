// Package funding provides methods to use funding resources via dwolla api.
package funding

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/pkg/errors"
)

// Resource represents a bank account connected to dwolla account.
type Resource struct {
	ID              string                 `json:"id"`
	Status          string                 `json:"status"`
	AccountNumber   string                 `json:"accountNumber"`
	RoutingNumber   string                 `json:"routingNumber"`
	BankAccountType string                 `json:"bankAccountType"`
	Name            string                 `json:"name"`
	BankName        string                 `json:"bankName"`
	Type            string                 `json:"type"`
	CreatedAt       string                 `json:"created"`
	Removed         bool                   `json:"removed"`
	PlaidToken      string                 `json:"plaidToken"`
	Channels        []string               `json:"channels"`
	Links           map[string]client.Link `json:"_links"`
}

// ListResourcesResponse is the response that is returned by dwolla
// to list funding resources request.
type ListResourcesResponse struct {
	Links    map[string]client.Link `json:"_links"`
	Embedded map[string][]Resource  `json:"_embedded"`
}

// VerifyMicroDepositsRequest is the request to verify microdeposits
type VerifyMicroDepositsRequest struct {
	Amount1 *Amount `json:"amount1"`
	Amount2 *Amount `json:"amount2"`
}

// Amount is the amount part of a balance.
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// BalanceResponse has the fields the describe the balance in a funding source.
type BalanceResponse struct {
	Links   map[string]map[string]client.Link `json:"_links"`
	Total   Amount                            `json:"total"`
	Balance Amount                            `json:"balance"`
}

// MicroDepositsDetails has the details for a microdeposits and their status.
type MicroDepositsDetails struct {
	Links     map[string]map[string]client.Link `json:"_links"`
	CreatedAt string                            `json:"created"`
	Status    string                            `json:"status"`
	Failure   map[string]string                 `json:"failure"`
}

// GetFundingSource retrieves a funding source by id.
func GetFundingSource(c *client.Client, sourceID string) (*Resource, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", c.RootURL()+"/funding-sources/"+sourceID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &Resource{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("customer not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// Update a funding source.
func (f *Resource) Update(c *client.Client) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(f)
	if err != nil {
		return errors.Wrap(err, "error marshalling json body for request")
	}
	req, err := http.NewRequest("POST", c.RootURL()+"/funding-sources/"+f.ID, bytes.NewReader(body))
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
	switch res.StatusCode {
	case 200:
		return nil
	case 404:
		return errors.New("funding source not found")
	case 400:
		return errors.New("only funding sources of type bank can be updated")
	case 403:
		return errors.New("a removed bank account cannot be updated")
	default:
		return errors.New(res.Status)
	}
}

// IntiateMicroDeposits for bank account verification.
func (f *Resource) IntiateMicroDeposits(c *client.Client) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("POST", f.Links["self"].Href+"/micro-deposits", nil)
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
	switch res.StatusCode {
	case 201:
		return nil
	case 404:
		return errors.New("funding source not found")
	default:
		return errors.New(res.Status)
	}
}

// VerifyMicroDeposits bank verification.
func (f *Resource) VerifyMicroDeposits(c *client.Client, vr *VerifyMicroDepositsRequest) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(vr)
	if err != nil {
		return errors.Wrap(err, "error marshalling verify micro deposits")
	}
	req, err := http.NewRequest("POST", f.Links["self"].Href+"/micro-deposits", bytes.NewReader(body))
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
	switch res.StatusCode {
	case 200:
		return nil
	case 202:
		return errors.New("Micro-deposits have not have not settled to destination bank. A Customer can verify these amounts after micro-deposits have processed to their bank")
	case 403:
		return errors.New("too many attempts, bank already verified")
	case 404:
		return errors.New("micro deposits not initiated or funding source not found")
	case 500:
		return errors.New("verify micro-deposits returned an unknown error")
	default:
		return errors.New(res.Status)
	}
}

// GetBalance retrieves balance for the funding source.
func (f *Resource) GetBalance(c *client.Client) (*BalanceResponse, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", f.Links["self"].Href+"/balance", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &BalanceResponse{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("funding source not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// GetMicroDepositsDetails retrieves the status of micro-deposits
// and checks if they are eligible for verification.
func (f *Resource) GetMicroDepositsDetails(c *client.Client) (*MicroDepositsDetails, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", f.Links["self"].Href+"/micro-deposits", nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &MicroDepositsDetails{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("micro-depostis not found or have already been verified")
	default:
		return nil, errors.New(res.Status)
	}
}

// Remove a funding resource.
func (f *Resource) Remove(c *client.Client) error {
	f.Removed = true
	err := f.Update(c)
	if err != nil {
		return errors.Wrap(err, "error removing funding source")
	}
	return nil
}
