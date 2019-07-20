package funding

import (
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/subosito/gotenv"
)

func TestGetFundingSource(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	_, err = GetFundingSource(client, "49dbaa24-1580-4b1c-8b58-24e26656fa31")
	if err != nil {
		t.Error(err)
	}
}
