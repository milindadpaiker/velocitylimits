package dal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const dialect = "sqlite3"

type sqliteDataStore struct {
	conn *gorm.DB
}

//NewSqliteDataStore returns sqlite datastore
func NewSqliteDataStore() (*sqliteDataStore, error) {
	db, err := gorm.Open(sqlite.Open("test2.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&Transaction{})
	return &sqliteDataStore{conn: db}, nil
}

func (m *sqliteDataStore) GetAllTxns(custID uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	m.conn.First(recentTxns, 10)
	return recentTxns, nil

}

//GetLastNValidTxns gets last "N" valid transactions for a customer ID. numberOfRecentTxn represents N
func (m *sqliteDataStore) GetLastNValidTxns(custID, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	return recentTxns, nil
}

func (m *sqliteDataStore) GetLastNTxns(custID uint, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	return recentTxns, nil
}

//SaveCustomerTxn saves customer transcation
func (m *sqliteDataStore) SaveCustomerTxn(txn *Transaction) error {
	m.conn.Create(txn)
	return nil
}
