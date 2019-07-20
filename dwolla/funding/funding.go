// Package funding provides methods to use funding resources via dwolla api.
package funding

import (
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
