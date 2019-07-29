package customer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
)

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

var mockDocuments = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/customers/176878b8-ecdb-469b-a82b-43ba5e8704b2/documents"
    }
  },
  "_embedded": {
    "documents": [
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
      },
      {
        "_links": {
          "self": {
            "href": "https://api-sandbox.dwolla.com/documents/11fe0bab-39bd-42ee-bb39-275afcc050d0"
          }
        },
        "id": "11fe0bab-39bd-42ee-bb39-275afcc050d0",
        "status": "pending",
        "type": "passport",
        "created": "2015-09-29T21:45:37.000Z"
      }
    ]
  },
  "total": 2
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
var mockFundingSources = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733/funding-sources"
    },
    "customer": {
      "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733"
    }
  },
  "_embedded": {
    "funding-sources": [
      {
        "_links": {
          "self": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/ab9cd5de-9435-47af-96fb-8d2fa5db51e8"
          },
          "customer": {
            "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733"
          },
          "with-available-balance": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/ab9cd5de-9435-47af-96fb-8d2fa5db51e8"
          }
        },
        "id": "ab9cd5de-9435-47af-96fb-8d2fa5db51e8",
        "status": "verified",
        "type": "balance",
        "name": "Balance",
        "created": "2015-10-02T21:00:28.153Z",
        "removed": false,
        "channels": []
      },
      {
        "_links": {
          "self": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/98c209d3-02d6-4bee-bc0f-61e18acf0e33"
          },
          "customer": {
            "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733"
          }
        },
        "id": "98c209d3-02d6-4bee-bc0f-61e18acf0e33",
        "status": "verified",
        "type": "bank",
        "bankAccountType": "checking",
        "name": "Jane Doe’s Checking",
        "created": "2015-10-02T22:03:45.537Z",
        "removed": false,
        "channels": [
            "ach"
        ],
        "fingerprint": "4cf31392f678cb26c62b75096e1a09d4465a801798b3d5c3729de44a4f54c794"
      }
    ]
  }
}
`
var mockFundingSourceToken = `
{
 "_links": {
   "self": {
     "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733/funding-sources-token"
   }
 },
 "token": "4adF858jPeQ9RnojMHdqSD2KwsvmhO7Ti7cI5woOiBGCpH5krY"
}
`

var mockIAVToken = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/customers/5b29279d-6359-4c87-a318-e09095532733/iav-token"
    }
  },
  "token": "4adF858jPeQ9RnojMHdqSD2KwsvmhO7Ti7cI5woOiBGCpH5krY"
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

func stubClient() *mockClient {
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
	return mock
}

func TestCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer := &Customer{FirstName: "Jane", LastName: "Merchant", Email: "jmerchantere13@nomailer.com", Type: "receive-only", BusinessName: "Jane corp llc", IPAddress: "99.99.99.99"}
	id, err := Create(mock, customer)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestList(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomers)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customers, err := List(mock)
	if err != nil {
		t.Error(err)
	}
	t.Log(customers[0].ID)
}

func TestGetCustomer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	t.Log(customer.LastName)
}

func TestUpdate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	customer.SSN = "1234"
	customer.Address = "12 baker street , london"
	customer.LastName = "Tester"
	customer.Status = "verified"
	err = customer.Update()
	if err != nil {
		t.Error(err)
	}
}

func TestAddDocument(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	file, err := ioutil.TempFile(".", "*.png")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(file.Name())
	err = customer.AddDocument(file, "passport")
	if err != nil {
		t.Error(err)
	}
}

func TestListDocuments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCustomer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocuments)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	_, err = customer.ListDocuments()
	if err != nil {
		t.Error(err)
	}
}

func TestGetDocument(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	doc, err := GetDocument(mock, "56502f7a-fa59-4a2f-8579-0f8bc9d7b9cc")
	if err != nil {
		t.Error(err)
	}
	t.Log("Document ID = ", doc.ID)
}

func TestCreateFundingSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	fr := &funding.Resource{
		RoutingNumber:   "222222226",
		AccountNumber:   "123456786",
		BankAccountType: "checking",
		Name:            "Jane Doe's checking",
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	err = customer.CreateFundingSource(fr)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateFundingSourceToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSourceToken)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	token, err := customer.CreateFundingSourceToken()
	if err != nil {
		t.Error(err)
	}
	t.Log("Token = ", token)
}
func TestCreateIAVFundingSourceToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockIAVToken)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	token, err := customer.CreateIAVFundingSourceToken()
	if err != nil {
		t.Error(err)
	}
	t.Log("Token = ", token)
}
func TestListFundingSources(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockDocument)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	customer, err := GetCustomer(mock, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSources)
	}))
	defer ts.Close()
	customer.Client.SetRootURL(ts.URL)
	sources, err := customer.ListFundingSources()
	if err != nil {
		t.Error(err)
	}
	t.Log("Count of sources = ", len(sources))
}
