package app

import "time"

// Receipt represents generic Receipt object
type Receipt struct {
	PostNum  int32
	PostAddr string

	Price    float64
	Currency int32

	IsBankCard bool
	OFD        string
	IsFiscal   bool

	IsService bool

	OperationTime time.Time
}

// NewReceipt constructs a Receipt object
func NewReceipt() *Receipt {
	return &Receipt{}
}
