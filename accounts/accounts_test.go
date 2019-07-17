package accounts

import (
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla/client"
	"github.com/subosito/gotenv"
)

func TestRetrieveAccount(t *testing.T) {
	gotenv.Load("../.env")
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
