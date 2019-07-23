// Package transfer has transfer related methods and structs.
package transfer

import (
	"bytes"
	"io"
	"encoding/json"
	"net/http"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
	"github.com/pkg/errors"
)

// Transfer has the fields to make a transfer between two funding sources.
type Transfer struct {
	Client        *client.Client
	Links         map[string]client.Link `json:"_links"`
	Amount        *funding.Amount        `json:"amount"`
	Metadata      map[string]string      `json:"metadata"`
	Fees          *Transfer              `json:"fees"`
	CorrelationID string                 `json:"correlationId"`
	Status        string                 `json:"status"`
	CreatedAt     string                 `json:"created"`
	Clearing      map[string]string      `json:"clearing"`
}

// ListTransferResponse is the response for list transfers end point.
type ListTransferResponse struct {
	Links    map[string]client.Link `json:"_links"`
	Embedded map[string][]Transfer  `json:"_embedded"`
}

// OnDemandAuthResponse is the response for the on-demand-authorization endpoint.
type OnDemandAuthResponse struct {
	Links      map[string]client.Link `json:"_links"`
	BodyText   string                 `json:"bodyText"`
	ButtonText string                 `json:"buttonText"`
}

// CreateTransfer initiates a new transfer between two funding sources.
func CreateTransfer(c *client.Client, transfer *Transfer) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(transfer)
	if err != nil {
		return errors.Wrap(err, "error marshalling the json body")
	}
	req, err := http.NewRequest("POST", c.RootURL()+"/transfers", bytes.NewReader(body))
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
	case 401:
		return errors.New("invalid access token")
	case 400:
		io.Copy(os.Stdout,res.Body)
		return errors.New(res.Status)
	case 403:
		return errors.New("Not authorized to create a transfer")
	case 404:
		return errors.New("not found")
	default:
		return errors.New(res.Status)
	}
}

// GetTransfer retrieves a transaction
func GetTransfer(c *client.Client, transferID string) (*Transfer, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", c.RootURL()+"/transfers/"+transferID, nil)
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
	case 201:
		d := json.NewDecoder(res.Body)
		body := &Transfer{}
		err = d.Decode(body)
		body.Client = c
		return body, nil
	case 404:
		return nil, errors.New("transfer not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// ListFees retrieves a list of the fees of the transfer
func (t *Transfer) ListFees() (*Transfer, error) {
	var c = t.Client
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", t.Links["self"].Href+"/fees", nil)
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
		body := &Transfer{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("transfer not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// Failure retrieves the failure reassons of a transfer.
func (t *Transfer) Failure() (*map[string]string, error) {
	var c = t.Client
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", t.Links["self"].Href+"/failure", nil)
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
		body := &map[string]string{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("transfer not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// Cancel a transfer.
func (t *Transfer) Cancel() (*Transfer, error) {
	var c = t.Client
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(`{"status" : "canceled"`)
	if err != nil {
		return nil, errors.Wrap(err, "error marshalling json requset body")
	}
	req, err := http.NewRequest("POST", t.Links["self"].Href, bytes.NewReader(body))
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
		body := &Transfer{}
		err = d.Decode(body)
		return body, nil
	case 404:
		return nil, errors.New("transfer not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// CreateOnDemandAuth create an on-demand bank transfer authorization for your Customer.
// On-demand authorization allows Customers to authorize Dwolla to transfer variable amounts
// from their bank account using ACH at a later point in time for products or services delivered.
// This on-demand authorization is supplied along with the Customer’s bank details when creating
// a new Customer funding source.
func CreateOnDemandAuth(c *client.Client) (string, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return "", errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("POST", c.RootURL()+"/on-demand-authorization", nil)
	if err != nil {
		return "", errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Add("Content-Type", "application/vnd.dwolla.v1.hal+json")
	res, err := hc.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		d := json.NewDecoder(res.Body)
		body := &OnDemandAuthResponse{}
		err = d.Decode(body)
		return body.Links["self"].Href, nil
	case 404:
		return "", errors.New("transfer not found")
	default:
		return "", errors.New(res.Status)
	}
}
