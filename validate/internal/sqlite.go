package dal

type sqliteDataStore struct{}

//NewSqliteDataStore returns sqlite datastore
func NewSqliteDataStore() *sqliteDataStore {
	return &sqliteDataStore{}
}

func (m *sqliteDataStore) GetAllTxns(custID int) ([]*Transaction, error) {
	var recentTxns []*Transaction
	return recentTxns, nil

}

//GetLastNValidTxns gets last "N" valid transactions for a customer ID. numberOfRecentTxn represents N
func (m *sqliteDataStore) GetLastNValidTxns(custID int, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	return recentTxns, nil
}

//SaveCustomerTxn saves customer transcation
func (m *sqliteDataStore) SaveCustomerTxn(txn *Transaction) error {
	return nil
}
