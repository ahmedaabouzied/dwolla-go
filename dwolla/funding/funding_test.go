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
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	t.Log("sourceID", source.ID)
}

func TestUpdateFundingSource(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	source.Name = "Bank of Brgle"
	err = source.Update(client)
	if err != nil {
		t.Error(err)
	}
}

func TestInitiateMicroDeposits(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	err = source.IntiateMicroDeposits(client)
	if err != nil {
		t.Error(err)
	}
}

func TestVerifyMicroDeposits(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
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
	err = source.VerifyMicroDeposits(client, vr)
	if err != nil {
		t.Error(err)
	}
}

func TestGetBalance(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	_, err = source.GetBalance(client)
	if err != nil {
		t.Error(err)
	}
}
func TestGetMicroDepositsDetails(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	_, err = source.GetMicroDepositsDetails(client)
	if err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	gotenv.Load("../../.env")
	client, err := client.CreateClient("sandbox", os.Getenv("DWOLLA_PUBLIC_KEY"), os.Getenv("DWOLLA_SECRET_KEY"))
	if err != nil {
		t.Error(err)
	}
	source, err := GetFundingSource(client, "5d2776f0-ba12-45ea-a7da-b4b01fa95ace")
	if err != nil {
		t.Error(err)
	}
	err = source.Remove(client)
	if err != nil {
		t.Error(err)
	}
}
