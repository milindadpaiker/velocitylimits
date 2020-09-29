package dal

var (
	//Memory db model saves valid and invalid transactions in different maps
	validTxns, invalidTxns map[uint][]*Transaction
)

func init() {
	validTxns = make(map[uint][]*Transaction)
	invalidTxns = make(map[uint][]*Transaction)
}

type memoryDataStore struct{}

//NewMemoryDataStore returns memory datastore
func NewMemoryDataStore() *memoryDataStore {
	return &memoryDataStore{}
}

func (m *memoryDataStore) GetAllTxns(custID uint) ([]*Transaction, error) {
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
func (m *memoryDataStore) GetLastNValidTxns(custID uint, numberOfRecentTxn uint) ([]*Transaction, error) {
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

//GetLastNValidTxns gets last "N" valid transactions for a customer ID. numberOfRecentTxn represents N
func (m *memoryDataStore) GetLastNTxns(custID uint, numberOfRecentTxn uint) ([]*Transaction, error) {
	var recentValidTxns, recentInvalidTxns []*Transaction
	var nv = numberOfRecentTxn
	var niv = numberOfRecentTxn
	if data, ok := validTxns[custID]; ok {
		for i := len(data) - 1; i >= 0 && nv > 0; i-- {
			recentValidTxns = append(recentValidTxns, data[i])
			nv = nv - 1
		}
	}
	if data, ok := invalidTxns[custID]; ok {
		for i := len(data) - 1; i >= 0 && niv > 0; i-- {
			recentInvalidTxns = append(recentInvalidTxns, data[i])
			niv = niv - 1
		}
	}

	var latestN []*Transaction
	var j, k int
	for i := 0; i < int(numberOfRecentTxn); i++ {
		if j < len(recentValidTxns) && k < len(recentInvalidTxns) {
			if recentValidTxns[j].Time.After(recentInvalidTxns[k].Time) {
				latestN = append(latestN, recentValidTxns[j])
				//latestN[i] = recentValidTxns[j]
				j = j + 1
			} else {
				latestN = append(latestN, recentInvalidTxns[k])
				//latestN[i] = recentInvalidTxns[k]
				k = k + 1
			}
		} else if j < len(recentValidTxns) {
			latestN = append(latestN, recentValidTxns[j])
			//latestN[i] = recentValidTxns[j]
			j = j + 1

		} else if k < len(recentInvalidTxns) {
			latestN = append(latestN, recentInvalidTxns[k])
			//latestN[i] = recentInvalidTxns[k]
			k = k + 1
		}

	}
	return latestN, nil
}

//SaveCustomerTxn saves customer transcation
func (m *memoryDataStore) SaveCustomerTxn(txn *Transaction) error {
	if txn.Status == Valid {
		validTxns[txn.CustomerID] = append(validTxns[txn.CustomerID], txn)
	} else {
		invalidTxns[txn.CustomerID] = append(invalidTxns[txn.CustomerID], txn)
	}
	return nil
}
