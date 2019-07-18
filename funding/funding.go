package funding

import "github.com/ahmedaabouzied/dwolla/client"

// Resource represents bank account connected to dwolla account
type Resource struct {
	ID              string        `json:"id"`
	Status          string        `json:"status"`
	AccountNumber   string        `json:"accountNumber"`
	RoutingNumber   string        `json:"routingNumber"`
	BankAccountType string        `json:"bankAccountType"`
	Name            string        `json:"name"`
	BankName        string        `json:"bankName"`
	Type            string        `json:"type"`
	Created         string        `json:"created"`
	Removed         bool          `json:"removed"`
	Channels        []string      `json:"channels"`
	Links           []client.Link `json:"_links"`
}

// ListResourcesResponse is the response that is returned by dwolla
// to list funding resources request
type ListResourcesResponse struct {
	Links   map[string]client.Link `json:"_links"`
	Embeded map[string][]Resource  `json:"_embedded"`
}
