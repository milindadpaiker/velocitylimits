package dal

import (
	"time"
)

//Status represents customer transaction status:- valid/invalid
type Status int

const (
	//Valid represents valid customer transaction after passing all rules
	Valid Status = 1 + iota
	//Invalid represents invalid customer transaction
	Invalid
)

//Transaction db model for transction object.
//Incidentally it is same as domain object Deposit, but that is not necessary
type Transaction struct {
	CustomerID uint `gorm:"primaryKey"`
	ID         int  `gorm:"primaryKey"`
	LoadAmount float64
	Time       time.Time
	Status     Status
}
