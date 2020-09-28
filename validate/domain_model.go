package validate

import (
	"time"

	dal "github.com/milind/velocitylimit/validate/internal"
)

//This file descries all structs/interfaces that are part of business domain

// Deposit is the individual fund load request
type Deposit struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	LoadAmount float64   `json:"load_amount"`
	Time       time.Time `json:"time"`
}

// DepositStatus shows status of individual fund load request
type DepositStatus struct {
	ID         int  `json:"id"`
	CustomerID int  `json:"customer_id"`
	Accepted   bool `json:"accepted"`
	//to be removed
	Amount float64
	Time   time.Time
}

// Rules indicate velocity limit rules
type Rules struct {
	DayLimit          float64
	WeekLimit         float64
	MaxAttemptsPerDay uint
}

//DataStore ...
type DataStore interface {
	RetrieveCustomerTxn(custID, loadID int) (*dal.Transaction, error)
	RetrieveCustomerTxns(custID int) ([]*dal.Transaction, error)
	SaveCustomerTxn(txn *dal.Transaction) error
}
