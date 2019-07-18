// Package masspayment provides methods to use mass payments via the dwolla api.
package masspayment

import "github.com/ahmedaabouzied/dwolla-go/dwolla/client"

// MassPayment represents a mass payment on dwolla api
type MassPayment struct {
	Links         map[string]client.Link `json:"_links"`
	ID            string                 `json:"id"`
	Status        string                 `json:"status"`
	CreatedAt     string                 `json:"created"`
	Metadata      map[string]string      `json:"metadata"`
	CorrelationID string                 `json:"correlationId"`
}

// ListMassPaymentsResponse is the response that is returned by dwolla
// to list mass payments
type ListMassPaymentsResponse struct {
	Links    map[string]client.Link   `json:"_links"`
	Embedded map[string][]MassPayment `json:"_embedded"`
}
