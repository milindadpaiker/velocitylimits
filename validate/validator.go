package validate

import (
	"context"
	"fmt"

	"github.com/milind/velocitylimit/config"
	dal "github.com/milind/velocitylimit/validate/internal"
)

//Validator ...
type Validator struct {
	rulechain []Rule
	dal       DataStore
}

func BuildRulesChain(cfg config.VelocityLimitConfig, ds DataStore) []Rule {
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
		//return &Validator{dal: dal.NewMemoryDataStore(), rulechain: rc}, nil
	case "sqlite":
		return nil, fmt.Errorf("ErrNotImplemented")
	default:
		return nil, fmt.Errorf("ErrInvalidDataStore")
	}

	//load config file and build rule chain
	var cfg config.VelocityLimitConfig
	cfg.DayLimit = 5000.00
	cfg.WeekLimit = 20000.00
	cfg.MaxAttemptsPerDay = 3
	rc := BuildRulesChain(cfg, ds)
	return &Validator{dal: dal.NewMemoryDataStore(), rulechain: rc}, nil
}

//Process is the maain validation functions that calls inidividual rules on incoming transaction.
func (v *Validator) Process(inFund *Deposit) (*DepositStatus, error) {
	//create context so that already fetched values can be passed around through rules
	//This way each rule need not access DB if previous rule has already fetched data.
	ctx := context.Background()

	isDuplicate, err := v.txnDuplicate(ctx, inFund)
	if isDuplicate {
		return nil, fmt.Errorf("ErrDuplicateTxn: %s", err.Error())
	}
	//if any other fault like DBAccessError treat it as error and do not process record
	if err != nil {
		return nil, err
	}

	//chain through various rules and validate
	for _, rule := range v.rulechain {

		valid, err := rule.Do(ctx, inFund)
		if err != nil {
			return &DepositStatus{}, err
		}
		if !valid {
			return &DepositStatus{
				ID:         inFund.ID,
				CustomerID: inFund.CustomerID,
				Accepted:   false,
				Amount:     inFund.LoadAmount,
				Time:       inFund.Time,
			}, nil
		}
	}
	//This is neither dupicate nor an invalid txn. Save and emit

	//transform to db model for saving to database
	txn := &dal.Transaction{
		CustomerID: inFund.CustomerID,
		ID:         inFund.ID,
		LoadAmount: inFund.LoadAmount,
		Time:       inFund.Time,
	}
	err = v.dal.SaveCustomerTxn(txn)
	if err != nil {

		return nil, err
	}
	//return status as accepted
	return &DepositStatus{
		ID:         inFund.ID,
		CustomerID: inFund.CustomerID,
		Accepted:   true,
		Amount:     inFund.LoadAmount,
		Time:       inFund.Time,
	}, nil
}

func (v *Validator) txnDuplicate(ctx context.Context, inFund *Deposit) (bool, error) {
	txn, err := v.dal.RetrieveCustomerTxns(inFund.CustomerID)
	if txn != nil {
		for _, t := range txn {
			if t.ID == inFund.ID {
				return true, fmt.Errorf("Original transaction on %v", t.Time)
			}
		}
	}
	return false, err

}
