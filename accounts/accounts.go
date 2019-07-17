package accounts

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Account represents a Dwolla master account that was estabslished on dwolla.com
type Account struct {
	selfURL             string // URL to the main account
	receiveRUL          string // URL to receive money to the account
	fundingResourcesURL string // URL to retrieve funding resources for the account
	transfersURL        string // URL to list transfers of the account
	customersURL        string // URL to list customers of the account
	sendURL             string // URL to send money from the account
	ID                  string // Dwolla account ID
	Name                string // Dwolla account holder name
}

type accountResponse struct {
	Links map[string]map[string]string `json:"_links"`
	ID    string                       `json:"id"`
	Name  string                       `json:"name"`
}

// RetrieveAccount returns the Dwolla master account
func RetrieveAccount(accountURL string, token string) (*Account, error) {
	hc := &http.Client{}
	req, err := http.NewRequest("GET", accountURL, nil)
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
