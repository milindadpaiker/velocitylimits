package validate

import (
	"context"
	"time"

	dal "github.com/milind/velocitylimit/validate/internal"
)

//This file descries all structs/interfaces that are part of business domain

// Deposit is the individual fund load request
type Deposit struct {
	ID         int       `json:"id"`
	CustomerID uint      `json:"customer_id"`
	LoadAmount float64   `json:"load_amount"`
	Time       time.Time `json:"time"`
}

// DepositStatus shows status of individual fund load request
type DepositStatus struct {
	ID         int  `json:"id"`
	CustomerID uint `json:"customer_id"`
	Accepted   bool `json:"accepted"`
	//to be removed
	Amount float64
	Time   time.Time
}

// VLRules indicate velocity limit rules
type VLRules struct {
	DayLimit          float64
	WeekLimit         float64
	MaxAttemptsPerDay uint
}

//Rules ...
type Rule interface {
	Do(context.Context, *Deposit) (bool, error)
	String() string
}

//DataStore ...
type DataStore interface {
	SaveCustomerTxn(txn *dal.Transaction) error
	GetAllTxns(custID uint) ([]*dal.Transaction, error)
	GetLastNValidTxns(custID uint, numberOfRecentTxn uint) ([]*dal.Transaction, error)
	GetLastNTxns(custID uint, numberOfRecentTxn uint) ([]*dal.Transaction, error)
}
