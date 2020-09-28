package dal

var a map[int][]*Transaction

func init() {
	a = make(map[int][]*Transaction)
}

type memoryDataStore struct{}

func NewMemoryDataStore() *memoryDataStore {
	return &memoryDataStore{}
}

//RetrieveCustomerTxn ...
func (m *memoryDataStore) RetrieveCustomerTxn(custID, loadID int) (*Transaction, error) {
	if data, ok := a[custID]; ok {
		for _, t := range data {
			if t.ID == loadID {
				return t, nil
			}
		}
	}
	return nil, nil
}

//RetrieveCustomerTxns get all customer trasactions ...
func (m *memoryDataStore) RetrieveCustomerTxns(custID int) ([]*Transaction, error) {
	if data, ok := a[custID]; ok {
		return data, nil
	}
	return nil, nil
}

//RetrieveRecentCustomerTxns get all customer trasactions ...
func (m *memoryDataStore) RetrieveRecentCustomerTxns(custID int, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentTxns []*Transaction
	if data, ok := a[custID]; ok {

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
	a[txn.CustomerID] = append(a[txn.CustomerID], txn)
	return nil
}
