package funding

// Resource represents bank account connected to dwolla account
type Resource struct {
	AccountNumber   uint
	RoutingNumber   uint
	BankAccountType string
	Name            string
}
