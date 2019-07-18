package customer

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/ahmedaabouzied/dwolla/client"
	"github.com/pkg/errors"
)

// Customer represents an individual or business with whom you intend to transact with
type Customer struct {
	ID           string                 `json:"id"`
	FirstName    string                 `json:"firstName"`
	LastName     string                 `json:"lastName"`
	Email        string                 `json:"email"`
	Type         string                 `json:"type"`
	Status       string                 `json:"status"`
	BusinessName string                 `json:"businessName"`
	IPAddress    string                 `json:"ipAddress"`
	CreatedAt    string                 `json:"created"`
	Links        map[string]client.Link `json:"_links"`
}

// Create a new customer
func Create(c *client.Client, cu *Customer) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(cu)
	if err != nil {
		return errors.Wrap(err, "error marshalling customer into req body")
	}
	req, err := http.NewRequest("POST", c.Links["customers"]["href"], bytes.NewReader(body))
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
	case 403:
		return errors.New("not authorized to create customers")
	case 404:
		return errors.New("account not found")
	default:
		return errors.New(res.Status)
	}
}
