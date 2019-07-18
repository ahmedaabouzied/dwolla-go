package customer

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

type listCustomersResponse struct {
	Links    map[string]client.Link `json:"_links"`
	Embedded map[string][]Customer  `json:"_embedded"`
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

// List retrieves a list of created customers
func List(c *client.Client) ([]Customer, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", c.Links["customers"]["href"], nil)
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
		body := &listCustomersResponse{}
		err = d.Decode(body)
		return body.Embedded["customers"], nil
	case 403:
		return nil, errors.New("not authorized to list customers")
	case 404:
		return nil, errors.New("account not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// GetCustomer retrieves a customer belonging to the authorized Dwolla Master Account by it's ID
func GetCustomer(c *client.Client, customerID string) (*Customer, error) {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get auth token")
	}
	req, err := http.NewRequest("GET", c.Links["customers"]["href"]+"/"+customerID, nil)
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
		body := &Customer{}
		err = d.Decode(body)
		return body, nil
	case 403:
		return nil, errors.New("not authorized to retrieve the customer")
	case 404:
		return nil, errors.New("account not found")
	default:
		return nil, errors.New(res.Status)
	}
}

// Update can be used to achieve the following :
//
//
// Update Customer information,
// upgrade an unverified Customer to a verified Customer,
// suspend a Customer, deactivate a Customer,
// reactivate a Customer,
// and update a verified Customer’s information to retry verification
//
func (cu *Customer) Update(c *client.Client) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	body, err := json.Marshal(cu)
	if err != nil {
		return errors.Wrap(err, "error marshalling customer into req body")
	}
	req, err := http.NewRequest("POST", cu.Links["self"].Href, bytes.NewReader(body))
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
	case 403:
		return errors.New("not authorized to update the customer")
	case 404:
		return errors.New("account not found")
	default:
		return errors.New(res.Status)
	}

}

// TODO : Add ListBusinessClassification Method

// TODO : Add RetrieveBusinessClassification Method

// AddDocument uploads a document to a customer for verification
func (cu *Customer) AddDocument(c *client.Client, file *os.File, documentType string) error {
	hc := &http.Client{}
	token, err := c.AuthToken()
	if err != nil {
		return errors.Wrap(err, "failed to get auth token")
	}
	if err != nil {
		return errors.Wrap(err, "error parsing file")
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// file, err := fileHeader.Open()
	// if err != nil {
	// 	return errors.Wrap(err, "error opening file")
	// }
	part, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		return errors.Wrap(err, "error uploading file")
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return errors.Wrap(err, "error uploading file")
	}
	err = writer.WriteField("documentType", documentType)
	if err != nil {
		return errors.Wrap(err, "error uploading file")
	}
	writer.Close()
	req, err := http.NewRequest("POST", cu.Links["self"].Href+"/documents", body)
	if err != nil {
		return errors.Wrap(err, "error creating the request")
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
	req.Header.Add("Cache-Control", "no-cache")
	res, err := hc.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request to dwolla api")
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 201:
		return nil
	case 403:
		return errors.New("not authorized to uplaod document to customer")
	case 404:
		return errors.New("account not found")
	default:
		return errors.New(res.Status)
	}

}
