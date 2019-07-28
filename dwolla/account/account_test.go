package account

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
)

var mockAccount string = `{
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

var mockFundingSources string = `
{
    "_links": {
        "self": {
            "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/funding-sources",
            "resource-type": "funding-source"
        }
    },
    "_embedded": {
        "funding-sources": [
            {
                "_links": {
                    "self": {
                        "href": "https://api-sandbox.dwolla.com/funding-sources/04173e17-6398-4d36-a167-9d98c4b1f1c3",
                        "type": "application/vnd.dwolla.v1.hal+json",
                        "resource-type": "funding-source"
                    },
                    "account": {
                        "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b",
                        "type": "application/vnd.dwolla.v1.hal+json",
                        "resource-type": "account"
                    }
                },
                "id": "04173e17-6398-4d36-a167-9d98c4b1f1c3",
                "status": "verified",
                "type": "bank",
                "bankAccountType": "checking",
                "name": "My Account - Checking",
                "created": "2017-09-25T20:03:41.000Z",
                "removed": false,
                "channels": [
                    "ach"
                ],
                "bankName": "First Midwestern Bank"
            },
            {
                "_links": {
                    "self": {
                        "href": "https://api-sandbox.dwolla.com/funding-sources/b268f6b9-db3b-4ecc-83a2-8823a53ec8b7"
                    },
                    "account": {
                        "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b"
                    },
                    "with-available-balance": {
                        "href": "https://api-sandbox.dwolla.com/funding-sources/b268f6b9-db3b-4ecc-83a2-8823a53ec8b7"
                    },
                    "balance": {
                        "href": "https://api-sandbox.dwolla.com/funding-sources/b268f6b9-db3b-4ecc-83a2-8823a53ec8b7/balance"
                    }
                },
                "id": "b268f6b9-db3b-4ecc-83a2-8823a53ec8b7",
                "status": "verified",
                "type": "balance",
                "name": "Balance",
                "created": "2017-08-22T18:21:51.000Z",
                "removed": false,
                "channels": []
            }
        ]
    }
}
`

var mockMassPayments string = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/mass-payments"
    },
    "first": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/mass-payments?limit=25&offset=0"
    },
    "last": {
      "href": "https://api-sandbox.dwolla.com/accounts/ca32853c-48fa-40be-ae75-77b37504581b/mass-payments?limit=25&offset=0"
    }
  },
  "_embedded": {
    "mass-payments": [
      {
        "_links": {
          "self": {
            "href": "https://api-sandbox.dwolla.com/mass-payments/b4b5a699-5278-4727-9f81-a50800ea9abc"
          },
          "source": {
            "href": "https://api-sandbox.dwolla.com/funding-sources/84c77e52-d1df-4a33-a444-51911a9623e9"
          },
          "items": {
            "href": "https://api-sandbox.dwolla.com/mass-payments/b4b5a699-5278-4727-9f81-a50800ea9abc/items"
          }
        },
        "id": "b4b5a699-5278-4727-9f81-a50800ea9abc",
        "status": "complete",
        "created": "2015-09-03T14:14:10.000Z",
        "metadata": {
          "UserJobId": "some ID"
        },
        "correlationId": "8a2cdc8d-629d-4a24-98ac-40b735229fe2"
      }
    ]
  },
  "total": 1
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

func stubAccount() *Account {
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
	mockAccount := &Account{
		Client: mock,
		Links:  mockLinks,
		ID:     "mock-account",
		Name:   "Master Mock",
	}
	return mockAccount
}
func TestRetrieveAccount(t *testing.T) {
	stubAcc := stubAccount()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockAccount)
	}))
	defer ts.Close()
	stubAcc.Client.SetRootURL(ts.URL)
	account, err := RetrieveAccount(stubAcc.Client)
	if err != nil {
		t.Error(err)
	}
	t.Log(account.ID)
}

func TestCreateFundingSource(t *testing.T) {
	stubAcc := stubAccount()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockAccount)
	}))
	defer ts.Close()
	stubAcc.Client.SetRootURL(ts.URL)
	fundingSource := &funding.Resource{
		Name:            "My Bank",
		RoutingNumber:   "222222226",
		AccountNumber:   "123456789",
		BankAccountType: "checking",
	}
	err := stubAcc.CreateFundingSource(fundingSource)
	if err != nil {
		t.Error(err)
	}
}

func TestListFundingSources(t *testing.T) {
	stubAcc := stubAccount()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w)
	}))
	defer ts.Close()

	stubAcc.Client.SetRootURL(ts.URL)
	_, err := stubAcc.ListFundingResources()
	if err != nil {
		t.Error(err)
	}
}

func TestListMassPayments(t *testing.T) {
	stubAcc := stubAccount()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockMassPayments)
	}))
	defer ts.Close()

	stubAcc.Client.SetRootURL(ts.URL)
	_, err := stubAcc.ListMassPayments()
	if err != nil {
		t.Error(err)
	}
}
