package validate

import (
	"fmt"

	dal "github.com/milind/velocitylimit/validate/internal"
)

type Validator struct {
	dal DataStore
}

//NewValidator ...
func NewValidator(backend string) (*Validator, error) {
	switch backend {
	case "memory":
		return &Validator{dal: dal.NewMemoryDataStore()}, nil
	case "sqlite":
		return nil, fmt.Errorf("ErrSqliteValidator")
	}
	return nil, fmt.Errorf("ErrValidator")

}

//Process is the maain validation functions that calls inidividual rules on incoming transaction.
func (v *Validator) Process(inFund *Deposit) (*DepositStatus, error) {

	isDuplicate, err := v.txnDuplicate(inFund)
	if isDuplicate {
		return nil, fmt.Errorf("ErrDuplicateTxn: %s", err.Error())
	}
	//if any other fault like DBAccessError treat it as error and do not process record
	if err != nil {
		return nil, err
	}

	//chain through various rules and validate

	//This is neither dupicate nor an invalid txn. Save and emit
	txn := &dal.Transaction{
		CustomerID: inFund.CustomerID,
		ID:         inFund.ID,
		LoadAmount: inFund.LoadAmount,
		Time:       inFund.Time,
	}

	err = v.dal.SaveCustomerTxn(txn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &DepositStatus{
		ID:         inFund.ID,
		CustomerID: inFund.CustomerID,
		Accepted:   true,
		Amount:     inFund.LoadAmount,
		Time:       inFund.Time,
	}, nil
}

func (v *Validator) txnDuplicate(inFund *Deposit) (bool, error) {
	txn, err := v.dal.RetrieveCustomerTxn(inFund.CustomerID, inFund.ID)
	if txn != nil && txn.CustomerID == inFund.CustomerID && txn.ID == inFund.ID {
		return true, fmt.Errorf("Original transaction on %v", txn.Time)
	}
	return false, err

}
