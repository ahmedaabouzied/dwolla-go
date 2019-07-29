package transfer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
)

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

var mockFees = `
{
  "transactions": [
    {
      "_links": {
        "self": {
          "href": "https://api-sandbox.dwolla.com/transfers/416a2857-c887-4cca-bd02-8c3f75c4bb0e"
        },
        "source": {
          "href": "https://api-sandbox.dwolla.com/funding-sources/AB443D36-3757-44C1-A1B4-29727FB3111C"
        },
        "destination": {
          "href": "https://api-sandbox.dwolla.com/funding-sources/707177c3-bf15-4e7e-b37c-55c3898d9bf4"
        },
        "created-from-transfer": {
          "href": "https://api-sandbox.dwolla.com/transfers/83eb4b5e-a5d9-e511-80de-0aa34a9b2388"
        }
      },
      "id": "416a2857-c887-4cca-bd02-8c3f75c4bb0e",
      "status": "pending",
      "amount": {
        "value": "2.00",
        "currency": "usd"
      },
      "created": "2016-02-22T20:46:38.777Z"
    },
    {
      "_links": {
        "self": {
          "href": "https://api-sandbox.dwolla.com/transfers/e58ae1f1-7007-47d3-a308-7e9aa6266d53"
        },
        "source": {
          "href": "https://api-sandbox.dwolla.com/funding-sources/AB443D36-3757-44C1-A1B4-29727FB3111C"
        },
        "destination": {
          "href": "https://api-sandbox.dwolla.com/funding-sources/ac6d4c2a-fda8-49f6-805d-468066dd474c"
        },
        "created-from-transfer": {
          "href": "https://api-sandbox.dwolla.com/transfers/83eb4b5e-a5d9-e511-80de-0aa34a9b2388"
        }
      },
      "id": "e58ae1f1-7007-47d3-a308-7e9aa6266d53",
      "status": "pending",
      "amount": {
        "value": "1.00",
        "currency": "usd"
      },
      "created": "2016-02-22T20:46:38.860Z"
    }
  ],
  "total": 2
}
`
var mockFailure = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/transfers/E6D9A950-AC9E-E511-80DC-0AA34A9B2388/failure"
    }
  },
  "code": "R01",
  "description": "Insufficient Funds"
}
`

var mockCanceled = `
{
  "_links": {
    "cancel": {
      "href": "https://api-sandbox.dwolla.com/transfers/3d48c13a-0fc6-e511-80de-0aa34a9b2388",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "self": {
      "href": "https://api-sandbox.dwolla.com/transfers/3d48c13a-0fc6-e511-80de-0aa34a9b2388",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "source": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "account"
    },
    "source-funding-source": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/73ce02cb-8857-4f01-83fc-b6640b24f9f7",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "funding-source"
    },
    "funding-transfer": {
      "href": "https://api-sandbox.dwolla.com/transfers/3c48c13a-0fc6-e511-80de-0aa34a9b2388",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "transfer"
    },
    "destination": {
      "href": "https://api-sandbox.dwolla.com/customers/33e56307-6754-41cb-81e2-23a7f1072295",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "customer"
    },
    "destination-funding-source": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/ac6d4c2a-fda8-49f6-805d-468066dd474c",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "funding-source"
    }
  },
  "id": "3d48c13a-0fc6-e511-80de-0aa34a9b2388",
  "status": "cancelled",
  "amount": {
    "value": "22.00",
    "currency": "USD"
  },
  "created": "2016-01-28T22:34:02.663Z",
  "metadata": {
    "foo": "bar",
    "baz": "boo"
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
func TestCreateTransfer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, "Created")
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
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
	tr := &Transfer{
		Amount: amount,
		Links:  links,
	}
	err := CreateTransfer(mock, tr)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTransfer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockTransfer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	transfer, err := GetTransfer(mock, "15c6bcce-46f7-e811-8112-e8dd3bececa8")
	if err != nil {
		t.Error(err)
	}
	t.Log("Transfer fees = ", transfer.ID)
}

func TestListFees(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockTransfer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	transfer, err := GetTransfer(mock, "15c6bcce-46f7-e811-8112-e8dd3bececa8")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFees)
	}))
	defer ts.Close()
	transfer.Client.SetRootURL(ts.URL)
	fees, err := transfer.ListFees()
	if err != nil {
		t.Error(err)
	}
	t.Log("Count of fees = ", fees.Total)
}

func TestFailure(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockTransfer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	transfer, err := GetTransfer(mock, "15c6bcce-46f7-e811-8112-e8dd3bececa8")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFailure)
	}))
	defer ts.Close()
	transfer.Client.SetRootURL(ts.URL)
	fail, err := transfer.Failure()
	if err != nil {
		t.Error(err)
	}
	t.Log("Failure code = ", fail.Code)
}

func TestCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockTransfer)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	transfer, err := GetTransfer(mock, "15c6bcce-46f7-e811-8112-e8dd3bececa8")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockCanceled)
	}))
	defer ts.Close()
	transfer.Client.SetRootURL(ts.URL)
	canceled, err := transfer.Cancel()
	if err != nil {
		t.Error(err)
	}
	t.Log("Transfer status = ", canceled.Status)
}

func TestCreateOnDemand(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockOnDemandAuth)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	link, err := CreateOnDemandAuth(mock)
	if err != nil {
		t.Error(err)
	}
	t.Log("On Demand Link = ", link)
}
