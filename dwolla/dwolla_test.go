package dwolla

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/customer"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/transfer"
)

var mockAccount = `{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b"
    },
    "receive": {
      "href": "https://api-sandbox.dwolla.com/transfers"
    },
    "funding-sources": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/funding-sources"
    },
    "transfers": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/transfers"
    },
    "customers": {
      "href": "https://api-sandbox.dwolla.com/customers"
    },
    "send": {
      "href": "https://api-sandbox.dwolla.com/transfers"
    }
  },
  "id": "ca32853c-48fa-40be-ae75-77b37504581b",
  "name": "Jane Doe"
}`

var mockCustomers = `
{
  "_links": {
    "first": {
      "href": "https://api-sandbox.dwolla.com/customers?limit=25&offset=0"
    },
    "last": {
      "href": "https://api-sandbox.dwolla.com/customers?limit=25&offset=0"
    },
    "self": {
      "href": "https://api-sandbox.dwolla.com/customers?limit=25&offset=0"
    }
  },
  "_embedded": {
    "customers": [
      {
        "_links": {
          "self": {
            "href": "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"
          }
        },
        "id": "FC451A7A-AE30-4404-AB95-E3553FCD733F",
        "firstName": "Jane",
        "lastName": "Doe",
        "email": "janedoe@nomail.com",
        "type": "unverified",
        "status": "unverified",
        "created": "2015-09-03T23:56:10.023Z"
      }
    ]
  },
  "total": 1
}
`
var mockFundingSource = `
{
    "_links": {
        "self": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/49dbaa24-1580-4b1c-8b58-24e26656fa31",
            "type": "application/vnd.dwolla.v1.hal+json",
            "resource-type": "funding-source"
        },
        "customer": {
            "href": "https://api-sandbox.dwolla.com/customers/4594a375-ca4c-4220-a36a-fa7ce556449d",
            "type": "application/vnd.dwolla.v1.hal+json",
            "resource-type": "customer"
        },
        "initiate-micro-deposits": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/49dbaa24-1580-4b1c-8b58-24e26656fa31/micro-deposits",
            "type": "application/vnd.dwolla.v1.hal+json",
            "resource-type": "micro-deposits"
        }
    },
    "id": "49dbaa24-1580-4b1c-8b58-24e26656fa31",
    "status": "unverified",
    "type": "bank",
    "bankAccountType": "checking",
    "name": "Test checking account",
    "created": "2017-09-26T14:14:08.000Z",
    "removed": false,
    "channels": [
        "ach"
    ],
    "bankName": "SANDBOX TEST BANK",
    "fingerprint": "5012989b55af15400e8102f95d2ec5e7ce3aef45c01613280d80a236dd8d6c3a"
}
`

var mockCustomer = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"
    }
  },
  "id": "FC451A7A-AE30-4404-AB95-E3553FCD733F",
  "firstName": "Jane",
  "lastName": "Doe",
  "email": "janedoe@nomail.com",
  "type": "unverified",
  "status": "unverified",
  "created": "2015-09-03T23:56:10.023Z"
}
`

var mockDocument = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/documents/56502f7a-fa59-4a2f-8579-0f8bc9d7b9cc"
    }
  },
  "id": "56502f7a-fa59-4a2f-8579-0f8bc9d7b9cc",
  "status": "pending",
  "type": "passport",
  "created": "2015-09-29T21:42:16.000Z"
}
`
var mockTransfer = `
{
  "_links": {
    "cancel": {
      "href": "https://api-sandbox.dwolla.com/transfers/15c6bcce-46f7-e811-8112-e8dd3bececa8",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "self": {
      "href": "https://api-sandbox.dwolla.com/transfers/15c6bcce-46f7-e811-8112-e8dd3bececa8",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "source": {
      "href": "https://api-sandbox.dwolla.com/accounts/62e88a41-f5d0-4a79-90b3-188cf11a3966",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "account"
    },
    "source-funding-source": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/12a0eaf9-9561-468d-bdeb-186b536aa2ed",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "funding-source"
    },
    "funding-transfer": {
      "href": "https://api-sandbox.dwolla.com/transfers/14c6bcce-46f7-e811-8112-e8dd3bececa8",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "destination": {
      "href": "https://api-sandbox.dwolla.com/customers/d295106b-ca20-41ad-9774-286e34fd3c2d",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "customer"
    },
    "destination-funding-source": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/500f8e0e-dfd5-431b-83e0-cd6632e63fcb",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "funding-source"
    }
  },
  "id": "15c6bcce-46f7-e811-8112-e8dd3bececa8",
  "status": "pending",
  "amount": {
    "value": "42.00",
    "currency": "USD"
  },
  "created": "2018-12-03T22:00:22.970Z",
  "clearing": {
    "source": "standard"
  }
}
`

var mockOnDemandAuth = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/on-demand-authorizations/30e7c028-0bdf-e511-80de-0aa34a9b2388"
    }
  },
  "bodyText": "I agree that future payments to Company ABC inc. will be processed by the Dwolla payment system from the selected account above. In order to cancel this authorization, I will change my payment settings within my Company ABC inc. account.",
  "buttonText": "Agree & Continue"
}
`

type mockClient struct {
	Env          string
	ClientID     string
	ClientSecret string
	authToken    string
	rootURL      string
	links        map[string]map[string]string
}

func (m *mockClient) RootURL() string {
	return m.rootURL
}
func (m *mockClient) Root() (map[string]map[string]string, error) {
	mockLinks := make(map[string]map[string]string)
	account := make(map[string]string)
	self := make(map[string]string)
	account["href"] = m.rootURL + "/account"
	self["href"] = m.rootURL
	mockLinks["self"] = self
	mockLinks["account"] = account
	return mockLinks, nil
}

func (m *mockClient) AuthToken() (string, error) {
	return m.authToken, nil
}
func (m *mockClient) Links() map[string]map[string]string {
	mockLinks := make(map[string]map[string]string)
	self := make(map[string]string)
	account := make(map[string]string)
	account["href"] = m.rootURL + "/account"
	self["href"] = m.rootURL
	mockLinks["self"] = self
	mockLinks["account"] = account
	return mockLinks
}
func (m *mockClient) SetAccessToken() error {
	return nil
}

func (m *mockClient) SetRootURL(url string) {
	m.rootURL = url
}

func stubClient() *Client {
	mock := &mockClient{
		Env:          "Test",
		ClientID:     "123456789",
		ClientSecret: "123456789",
		authToken:    "abcdefghijklmn",
		rootURL:      "http://localhost:8080",
	}
	mockLinks := make(map[string]map[string]string)
	self := make(map[string]string)
	account := make(map[string]string)
	account["href"] = mock.rootURL + "/account/"
	fundingSources := make(map[string]string)
	fundingSources["href"] = mock.rootURL + "/funding-sources/"
	self["href"] = mock.rootURL
	mockLinks["self"] = self
	mockLinks["account"] = account
	mockLinks["funding-sources"] = fundingSources
	mock.links = mockLinks
	mc := &Client{
		Client: mock,
	}
	return mc
}
func TestRetrieveAccount(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockAccount)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	account, err := mock.RetrieveAccount()
	if err != nil {
		t.Error(err)
	}
	t.Log("Account ID = ", account.ID)
}

func TestCreateCustomer(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	customer := &customer.Customer{FirstName: "Jane", LastName: "Merchant", Email: "jmerchantere13@nomailer.com", Type: "receive-only", BusinessName: "Jane corp llc", IPAddress: "99.99.99.99"}
	_, err := mock.CreateCustomer(customer)
	if err != nil {
		t.Error(err)
	}
}

func TestListCustomers(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomers)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	customers, err := mock.ListCustomers()
	if err != nil {
		t.Error(err)
	}
	t.Log("Count of customers = ", len(customers))
}

func TestGetCustomer(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	customer, err := mock.GetCustomer("ca32853c-48fa-40be-ae75-77b37504581b")
	if err != nil {
		t.Error(err)
	}
	t.Log("Customer ID = ", customer.ID)
}

func TestGetDocument(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	doc, err := mock.GetDocument("56502f7a-fa59-4a2f-8579-0f8bc9d7b9cc")
	if err != nil {
		t.Error(err)
	}
	t.Log("Document ID = ", doc.ID)
}

func TestGetFundingSource(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	doc, err := mock.GetFundingSource("49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	t.Log("Funding Source ID = ", doc.ID)
}

func TestCreateTransfer(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, "Created")
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	amount := &funding.Amount{
		Value:    "300",
		Currency: "USD",
	}
	links := make(map[string]client.Link)
	links["source"] = client.Link{
		Href: "https://api-sandbox.dwolla.com/funding-sources/707177c3-bf15-4e7e-b37c-55c3898d9bf4",
	}
	links["destination"] = client.Link{
		Href: "https://api-sandbox.dwolla.com/funding-sources/AB443D36-3757-44C1-A1B4-29727FB3111C",
	}
	tr := &transfer.Transfer{
		Amount: amount,
		Links:  links,
	}
	err := mock.CreateTransfer(tr)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransfer(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockTransfer)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	transfer, err := mock.GetTransfer("15c6bcce-46f7-e811-8112-e8dd3bececa8")
	if err != nil {
		t.Error(err)
	}
	t.Log("Transfer ID = ", transfer.ID)
}

func TestCreateOnDemandAuth(t *testing.T) {
	mock := stubClient()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockOnDemandAuth)
	}))
	defer ts.Close()
	mock.Client.SetRootURL(ts.URL)
	link, err := mock.CreateOnDemandAuth()
	if err != nil {
		t.Error(err)
	}
	t.Log("On Demand Link = ", link)
}
