package customer

import (
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla/client"
	"github.com/subosito/gotenv"
)

func TestCreate(t *testing.T) {
	gotenv.Load("../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customer := &Customer{FirstName: "Jane", LastName: "Merchant", Email: "jmerchante@nomail.com", Type: "receive-only", BusinessName: "Jane corp llc", IPAddress: "99.99.99.99"}
	err = Create(client, customer)
	if err != nil {
		t.Error(err)
	}
}
