package dal

import "time"

//Transaction db model for transction object.
//Incidentally it is same as domain object Deposit, but that is not necessary
type Transaction struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	LoadAmount float64   `json:"load_amount"`
	Time       time.Time `json:"time"`
}
