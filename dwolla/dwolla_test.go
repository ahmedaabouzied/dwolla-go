package dwolla

import (
	"os"
	"testing"

	"github.com/subosito/gotenv"
)

func TestCreateClient(t *testing.T) {
	gotenv.Load("../.env")
	client, err := CreateClient(Sandbox, os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	t.Log(client)
}
