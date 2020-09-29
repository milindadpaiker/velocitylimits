package validate

import (
	"context"

	"github.com/milind/velocitylimit/config"
)

type weeklyLimitsRule struct {
	weekLimit float64
	attempts  uint
	dal       DataStore
}

func NewWeeklyLimitsRule(cfg config.Config, ds DataStore) *weeklyLimitsRule {
	return &weeklyLimitsRule{weekLimit: cfg.WeekLimit, dal: ds, attempts: cfg.MaxAttemptsPerDay}
}

//Do for weeklyLimitsRule validates if total amount of transactions per week is within limits
func (da *weeklyLimitsRule) Do(ctx context.Context, deposit *Deposit) (bool, error) {
	//low hanging fruit
	if deposit.LoadAmount > da.weekLimit {
		return false, nil
	}
	txn, err := da.dal.GetLastNValidTxns(deposit.CustomerID, 7*da.attempts)
	if err != nil {
		return false, err
	}
	var currentWeeklyTotal float64
	txnYr, txnWeek := deposit.Time.ISOWeek()
	if txn != nil {
		for _, t := range txn {
			//minor optimization possible which avoids full array traversal
			ty, tw := t.Time.ISOWeek()
			if txnYr == ty && txnWeek == tw {
				currentWeeklyTotal = currentWeeklyTotal + t.LoadAmount
			}
		}
		if deposit.LoadAmount > (da.weekLimit - currentWeeklyTotal) {
			return false, nil
		}
	}
	return true, nil
}

func (da *weeklyLimitsRule) String() string {
	return "WeeklyLimitsRule"
}
