package dal

import "time"

type Status int

const (
	Valid Status = 1 + iota
	Invalid
)

//Transaction db model for transction object.
//Incidentally it is same as domain object Deposit, but that is not necessary
type Transaction struct {
	ID         int
	CustomerID int
	LoadAmount float64
	Time       time.Time
	Status     Status
}
