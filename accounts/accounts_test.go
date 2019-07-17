package accounts

import (
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

var clientLink = "https://api-sandbox.dwolla.com/accounts/bdb5a377-8c99-443c-b05d-0597bb656a83"

func TestRetrieveAccount(t *testing.T) {
	gotenv.Load("../.env")
	token := os.Getenv("DWOLLA_TOKEN")
	account, err := RetrieveAccount(clientLink, token)
	if err != nil {
		t.Error(err)
	}
	t.Log("account ID : ", account.ID)
}
