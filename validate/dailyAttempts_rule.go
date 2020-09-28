package validate

import (
	"context"

	"github.com/milind/velocitylimit/config"
)

type dailyAttemptsRule struct {
	attempts uint
	dal      DataStore
}

func NewDailyAttemptsRule(cfg config.VelocityLimitConfig, ds DataStore) *dailyAttemptsRule {
	return &dailyAttemptsRule{attempts: cfg.MaxAttemptsPerDay, dal: ds}
}

func (da *dailyAttemptsRule) Do(ctx context.Context, deposit *Deposit) (bool, error) {
	//Get recent customer transactions.Not all. As many as max daily limit
	txn, err := da.dal.RetrieveRecentCustomerTxns(deposit.CustomerID, da.attempts)
	if err != nil {
		return false, err
	}
	if len(txn) < int(da.attempts) {
		return true, nil
	}
	if txn != nil {
		for _, t := range txn {
			if !(deposit.Time.Day() == t.Time.Day() && deposit.Time.Month() == t.Time.Month() && deposit.Time.Year() == t.Time.Year()) {
				return true, nil
			}
		}
		return false, nil
	}
	return true, nil
}
