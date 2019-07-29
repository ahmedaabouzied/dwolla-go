package funding

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

var mockBalance = `
{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/c2eb3f03-1b0e-4d18-a4a2-e552cc111418/balance",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "balance"
    },
    "funding-source": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/c2eb3f03-1b0e-4d18-a4a2-e552cc111418",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "funding-source"
    }
  },
  "balance": {
    "value": "4616.87",
    "currency": "USD"
  },
  "total": {
      "value": "4616.87",
      "currency": "USD"
  },
  "lastUpdated": "2017-04-18T15:20:25.880Z"
}
`
var mockMicroDepositsDetails = `

{
  "_links": {
    "self": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/dfe59fdd-7467-44cf-a339-2020dab5e98a/micro-deposits",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "micro-deposits"
    },
    "verify-micro-deposits": {
      "href": "https://api-sandbox.dwolla.com/funding-sources/dfe59fdd-7467-44cf-a339-2020dab5e98a/micro-deposits",
      "type": "application/vnd.dwolla.v1.hal+json",
      "resource-type": "micro-deposits"
    }
  },
  "created": "2016-12-30T20:56:53.000Z",
  "status": "failed",
  "failure": {
    "code": "R03",
    "description": "No Account/Unable to Locate Account"
  }
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
func TestGetFundingSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	t.Log("Funding Source ID = ", fundingSource.ID)
}

func TestUpdateFundingSource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	fundingSource.Name = "Bank of Brgle"
	err = fundingSource.Update()
	if err != nil {
		t.Error(err)
	}
}

func TestInitiateMicroDeposits(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	fundingSource.Client.SetRootURL(ts.URL)
	err = fundingSource.IntiateMicroDeposits()
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyMicroDeposits(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}

	vr := &VerifyMicroDepositsRequest{
		Amount1: &Amount{
			Value:    "0.03",
			Currency: "USD",
		},
		Amount2: &Amount{
			Value:    "0.03",
			Currency: "USD",
		},
	}
	err = fundingSource.VerifyMicroDeposits(vr)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBalance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockBalance)
	}))
	defer ts.Close()
	fundingSource.Client.SetRootURL(ts.URL)
	balance, err := fundingSource.GetBalance()
	if err != nil {
		t.Error(err)
	}
	t.Log("Balance = ", balance.Total.Value+balance.Total.Currency)
}
func TestGetMicroDepositsDetails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockMicroDepositsDetails)
	}))
	defer ts.Close()
	fundingSource.Client.SetRootURL(ts.URL)
	details, err := fundingSource.GetMicroDepositsDetails()
	if err != nil {
		t.Error(err)
	}
	t.Log("Micro-deposit status = ", details.Status)
}

func TestRemove(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, mockFundingSource)
	}))
	defer ts.Close()
	mock := stubClient()
	mock.SetRootURL(ts.URL)
	fundingSource, err := GetFundingSource(mock, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
	err = fundingSource.Remove()
	if err != nil {
		t.Error(err)
	}
}
