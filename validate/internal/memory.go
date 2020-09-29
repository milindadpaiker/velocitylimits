package dal

var (
	validTxns, invalidTxns map[int][]*Transaction
)

func init() {
	validTxns = make(map[int][]*Transaction)
	invalidTxns = make(map[int][]*Transaction)
}

type memoryDataStore struct{}

func NewMemoryDataStore() *memoryDataStore {
	return &memoryDataStore{}
}

func (m *memoryDataStore) GetAllTxns(custID int) ([]*Transaction, error) {
	var recentTxns []*Transaction
	if data, ok := validTxns[custID]; ok {
		recentTxns = append(recentTxns, data...)
	}
	if data, ok := invalidTxns[custID]; ok {
		recentTxns = append(recentTxns, data...)
	}
	return recentTxns, nil

}

//GetLastNValidTxns gets last "N" valid transactions for a customer ID. numberOfRecentTxn represents N
func (m *memoryDataStore) GetLastNValidTxns(custID int, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	if data, ok := validTxns[custID]; ok {

		for i := len(data) - 1; i >= 0 && numberOfRecentTxn > 0; i-- {

			recentTxns = append(recentTxns, data[i])
			numberOfRecentTxn = numberOfRecentTxn - 1
		}
		return recentTxns, nil
	}
	return nil, nil
}

//SaveCustomerTxn ...
func (m *memoryDataStore) SaveCustomerTxn(txn *Transaction) error {
	if txn.Status == Valid {
		validTxns[txn.CustomerID] = append(validTxns[txn.CustomerID], txn)
	} else {
		invalidTxns[txn.CustomerID] = append(invalidTxns[txn.CustomerID], txn)
	}
	return nil
}
