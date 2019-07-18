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

func TestList(t *testing.T) {
	gotenv.Load("../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customers, err := List(client)
	if err != nil {
		t.Error(err)
	}
	t.Log(customers[0].ID)
}

func TestGetCustomer(t *testing.T) {
	gotenv.Load("../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customers, err := List(client)
	if err != nil {
		t.Error(err)
	}
	customer, err := GetCustomer(client, customers[0].ID)
	if err != nil {
		t.Error(err)
	}
	t.Log(customer.LastName)
}

func TestUpdate(t *testing.T) {
	gotenv.Load("../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customers, err := List(client)
	if err != nil {
		t.Error(err)
	}
	customer, err := GetCustomer(client, customers[0].ID)
	if err != nil {
		t.Error(err)
	}
	customer.LastName = "Doe"
	err = customer.Update(client)
	if err != nil {
		t.Error(err)
	}
}

func TestAddDocument(t *testing.T) {
	gotenv.Load("../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customers, err := List(client)
	if err != nil {
		t.Error(err)
	}
	customer, err := GetCustomer(client, customers[0].ID)
	if err != nil {
		t.Error(err)
	}
	// Create a file called test.png before running tests
	file, err := os.Open("test.png")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	err = customer.AddDocument(client, file, "passport")
	if err != nil {
		t.Error(err)
	}
}
