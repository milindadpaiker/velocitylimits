package validate

import (
	"context"

	"github.com/milind/velocitylimit/config"
)

type dailyLimitsRule struct {
	dayLimit float64
	attempts uint
	dal      DataStore
}

func NewDailyLimitsRule(cfg config.VelocityLimitConfig, ds DataStore) *dailyLimitsRule {
	return &dailyLimitsRule{dayLimit: cfg.DayLimit, attempts: cfg.MaxAttemptsPerDay, dal: ds}
}

func (da *dailyLimitsRule) Do(ctx context.Context, deposit *Deposit) (bool, error) {
	//low hanging fruit
	if deposit.LoadAmount > da.dayLimit {
		return false, nil
	}
	txn, err := da.dal.GetLastNValidTxns(deposit.CustomerID, da.attempts)
	if err != nil {
		return false, err
	}
	var currentDailyTotal float64
	if txn != nil {
		for _, t := range txn {
			if deposit.Time.Day() == t.Time.Day() && deposit.Time.Month() == t.Time.Month() && deposit.Time.Year() == t.Time.Year() {
				currentDailyTotal = currentDailyTotal + t.LoadAmount
			}
		}
		if deposit.LoadAmount > (da.dayLimit - currentDailyTotal) {
			return false, nil
		}
	}
	return true, nil
}
