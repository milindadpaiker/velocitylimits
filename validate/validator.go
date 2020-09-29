package validate

import (
	"context"
	"fmt"
	"log"

	"github.com/milind/velocitylimit/config"
	dal "github.com/milind/velocitylimit/validate/internal"
)

//Validator struct holds the chain of rules it validates agaist and the backing datastore
type Validator struct {
	rulechain []Rule
	dal       DataStore
}

//BuildRulesChain creates chain of rules
func BuildRulesChain(cfg config.Config, ds DataStore) []Rule {
	var rc []Rule
	rc = append(rc, NewDailyAttemptsRule(cfg, ds))
	rc = append(rc, NewDailyLimitsRule(cfg, ds))
	rc = append(rc, NewWeeklyLimitsRule(cfg, ds))
	return rc
}

//NewValidator ...
func NewValidator(backend string) (*Validator, error) {
	//create backing datastore first
	var ds DataStore
	switch backend {
	case "memory":
		ds = dal.NewMemoryDataStore()
	case "sqlite":
		return nil, fmt.Errorf("ErrSqliteNotImplemented")
	default:
		return nil, fmt.Errorf("ErrInvalidBackend")
	}

	rc := BuildRulesChain(config.Configuration, ds)
	return &Validator{dal: dal.NewMemoryDataStore(), rulechain: rc}, nil
}

//Process is the maain validation functions that calls inidividual rules on incoming transaction.
func (v *Validator) Process(inFund *Deposit) (*DepositStatus, error) {
	//create context so that already fetched values can be passed around into rules
	//This way each rule need not access DB if previous rule has already fetched required data.
	ctx := context.Background()

	//two-fold responsibility
	//1. To identify and ignore repeat transaction (repeating loadId for a customer)
	//2. If programm is in recovery mode/completing last incomplete run, it avoids re-processing already processed funds
	isDuplicate, err := v.txnDuplicate(ctx, inFund)
	if isDuplicate {
		return nil, fmt.Errorf("ErrDuplicateTxn: %s", err.Error())
	}
	//if any other faults like DBAccessError treat it as error and do not process record
	if err != nil {
		return nil, err
	}

	var validTxn = true
	var accepted bool

	//chain through various rules and validate
	for i := 0; i < len(v.rulechain) && validTxn; i++ {
		validTxn, err = v.rulechain[i].Do(ctx, inFund)
		if !validTxn {
			log.Printf("valiation %s failed for transaction %+v", v.rulechain[i].String(), inFund)
		}
		//if error do not process this txn. Neither treat as valid or invalid but return
		if err != nil {
			return &DepositStatus{}, err
		}
	}

	//create the db model instance of txn and save
	txn := &dal.Transaction{
		CustomerID: inFund.CustomerID,
		ID:         inFund.ID,
		LoadAmount: inFund.LoadAmount,
		Time:       inFund.Time,
	}

	if !validTxn {
		txn.Status = dal.Invalid
	} else {
		txn.Status = dal.Valid
		accepted = true
	}
	err = v.dal.SaveCustomerTxn(txn)
	if err != nil {
		//if failed to the txn return without emitting
		log.Printf("Failed to save customer transaction. %+v", txn)
		return nil, err
	}
	//if successfully processed and saved by validator module, only then emit the result to output
	return &DepositStatus{
		ID:         inFund.ID,
		CustomerID: inFund.CustomerID,
		Accepted:   accepted,
		Amount:     inFund.LoadAmount,
		Time:       inFund.Time,
	}, nil
}

func (v *Validator) txnDuplicate(ctx context.Context, inFund *Deposit) (bool, error) {
	txn, err := v.dal.GetAllTxns(inFund.CustomerID)
	if txn != nil {
		for _, t := range txn {
			if t.ID == inFund.ID {
				return true, fmt.Errorf("Original transaction on %v", t.Time)
			}
		}
	}
	return false, err
}
