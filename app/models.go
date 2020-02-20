package app

// Receipt represents generic Receipt object
type Receipt struct {
	PostNum  string
	PostAddr string

	Price    string
	Currency string

	IsBankCard string
	OFD        string
	IsFiscal   string

	IsService string

	OperationTime string
}

// NewReceipt constructs a Receipt object
func NewReceipt() *Receipt {
	return &Receipt{}
}
