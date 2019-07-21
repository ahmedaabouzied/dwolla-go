package account

import (
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
	"github.com/subosito/gotenv"
)

func TestRetrieveAccount(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	account, err := RetrieveAccount(client)
	if err != nil {
		t.Error(err)
	}
	t.Log("account ID : ", account.ID)
	t.Log("account link :", account.Links["self"]["href"])
}

func TestCreateFundingSource(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	account, err := RetrieveAccount(client)
	if err != nil {
		t.Error(err)
	}
	fundingSource := &funding.Resource{
		Name:            "My Bank",
		RoutingNumber:   "222222226",
		AccountNumber:   "123456789",
		BankAccountType: "checking",
	}
	err = account.CreateFundingSource(fundingSource)
	if err != nil {
		t.Error(err)
	}
}

func TestListFundingSources(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	account, err := RetrieveAccount(client)
	if err != nil {
		t.Error(err)
	}
	sources, err := account.ListFundingResources()
	if err != nil {
		t.Error(err)
	}
	t.Log(sources[0].ID)
}

func TestListMassPayments(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	account, err := RetrieveAccount(client)
	if err != nil {
		t.Error(err)
	}
	_, err = account.ListMassPayments()
	if err != nil {
		t.Error(err)
	}
}
