package customer

import (
	"os"
	"testing"

	"github.com/ahmedaabouzied/dwolla-go/dwolla/client"
	"github.com/ahmedaabouzied/dwolla-go/dwolla/funding"
	"github.com/subosito/gotenv"
)

func TestCreate(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	customer := &Customer{FirstName: "Jane", LastName: "Merchant", Email: "jmerchantere13@nomailer.com", Type: "receive-only", BusinessName: "Jane corp llc", IPAddress: "99.99.99.99"}
	id, err := Create(client, customer)
	if err != nil {
		t.Error(err)
	}
	t.Log(id)
}

func TestList(t *testing.T) {
	gotenv.Load("../../.env")
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
	gotenv.Load("../../.env")
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
	gotenv.Load("../../.env")
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
	err = customer.Update()
	if err != nil {
		t.Error(err)
	}
}

// func TestAddDocument(t *testing.T) {
// 	gotenv.Load("../../.env")
// 	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	customers, err := List(client)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	customer, err := GetCustomer(client, customers[0].ID)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// Create a file called test.png before running tests
// 	file, err := os.Open("test.png")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer file.Close()
// 	err = customer.AddDocument(client, file, "passport")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestListDocuments(t *testing.T) {
	gotenv.Load("../../.env")
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
	documnets, err := customer.ListDocuments()
	if err != nil {
		t.Error(err)
	}
	t.Log("number of docs : ", len(documnets))
}

func TestGetDocument(t *testing.T) {
	gotenv.Load("../../.env")
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
	documents, err := customer.ListDocuments()
	if err != nil {
		t.Error(err)
	}
	if len(documents) > 0 {
		document, err := GetDocument(client, documents[0].ID)
		if err != nil {
			t.Error(err)
		}
		t.Log(document.ID)
	}
}

func TestCreateFundingSource(t *testing.T) {
	gotenv.Load("../../.env")
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
	fr := &funding.Resource{
		RoutingNumber:   "222222226",
		AccountNumber:   "123456786",
		BankAccountType: "checking",
		Name:            "Jane Doe's checking",
	}
	err = customer.CreateFundingSource(fr)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateFundingSourceToken(t *testing.T) {
	gotenv.Load("../../.env")
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
	token, err := customer.CreateFundingSourceToken()
	if err != nil {
		t.Error(err)
	}
	t.Log("token : ", token)
}
func TestCreateIAVFundingSourceToken(t *testing.T) {
	gotenv.Load("../../.env")
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
	token, err := customer.CreateIAVFundingSourceToken()
	if err != nil {
		t.Error(err)
	}
	t.Log("token : ", token)
}
func TestListFundingSources(t *testing.T) {
	gotenv.Load("../../.env")
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
	sources, err := customer.ListFundingSources()
	if err != nil {
		t.Error(err)
	}
	t.Log(len(sources))
}
