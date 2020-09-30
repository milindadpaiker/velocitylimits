package dal

import (
	"reflect"
	"testing"
)

func TestNewMemoryDataStore(t *testing.T) {
	tests := []struct {
		name string
		want *memoryDataStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemoryDataStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemoryDataStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryDataStore_GetAllTxns(t *testing.T) {
	type args struct {
		custID uint
	}
	tests := []struct {
		name    string
		m       *memoryDataStore
		args    args
		want    []*Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetAllTxns(tt.args.custID)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryDataStore.GetAllTxns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryDataStore.GetAllTxns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryDataStore_GetLastNValidTxns(t *testing.T) {
	type args struct {
		custID            uint
		numberOfRecentTxn uint
	}
	tests := []struct {
		name    string
		m       *memoryDataStore
		args    args
		want    []*Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetLastNValidTxns(tt.args.custID, tt.args.numberOfRecentTxn)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryDataStore.GetLastNValidTxns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryDataStore.GetLastNValidTxns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryDataStore_GetLastNTxns(t *testing.T) {
	type args struct {
		custID            uint
		numberOfRecentTxn uint
	}
	tests := []struct {
		name    string
		m       *memoryDataStore
		args    args
		want    []*Transaction
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetLastNTxns(tt.args.custID, tt.args.numberOfRecentTxn)
			if (err != nil) != tt.wantErr {
				t.Errorf("memoryDataStore.GetLastNTxns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("memoryDataStore.GetLastNTxns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_memoryDataStore_SaveCustomerTxn(t *testing.T) {
	type args struct {
		txn *Transaction
	}
	tests := []struct {
		name    string
		m       *memoryDataStore
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SaveCustomerTxn(tt.args.txn); (err != nil) != tt.wantErr {
				t.Errorf("memoryDataStore.SaveCustomerTxn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
