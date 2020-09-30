package validate

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/milind/velocitylimit/config"
	dal "github.com/milind/velocitylimit/validate/internal"
)

func TestNewValidator(t *testing.T) {
	type args struct {
		backend string
	}
	tests := []struct {
		name    string
		args    args
		want    *Validator
		wantErr bool
	}{
		{
			name:    "default",
			args:    args{backend: "test"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "sqlite",
			args:    args{backend: "sqlite"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewValidator(tt.args.backend)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.backend == "sqlite" {
				if _, err := os.Stat("velocitylimits.db"); os.IsNotExist(err) {
					t.Error("Failed to create sqlite db: velocitylimits.db")
				} else {
					fmt.Println(os.Remove("velocitylimits.db"))
				}
			}
		})
	}
}

func TestValidator_Process(t *testing.T) {
	type args struct {
		inFund *Deposit
	}
	config.Configuration.DayLimit = 5000.00
	config.Configuration.WeekLimit = 20000.00
	config.Configuration.MaxAttemptsPerDay = 3
	config.Configuration.Currency = "$"
	vm, _ := NewValidator("memory")

	tests := []struct {
		name    string
		v       *Validator
		args    args
		want    *DepositStatus
		wantErr bool
		given   []*dal.Transaction
	}{
		{
			name:    "1",
			v:       vm,
			args:    args{inFund: &Deposit{ID: 345, CustomerID: 123, LoadAmount: 6000, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC)}},
			want:    &DepositStatus{ID: 345, CustomerID: 123, Accepted: false},
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 11, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 12, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 2890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 13, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "duplicate_txn",
			v:       vm,
			args:    args{inFund: &Deposit{ID: 890, CustomerID: 123, LoadAmount: 100, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC)}},
			want:    nil,
			wantErr: true,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 11, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 12, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 2890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 13, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "valid_txn",
			v:       vm,
			args:    args{inFund: &Deposit{ID: 345, CustomerID: 1223, LoadAmount: 1000, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC)}},
			want:    &DepositStatus{ID: 345, CustomerID: 1223, Accepted: true},
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 11, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 12, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 2890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 13, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupDataStore(tt.v.dal, tt.given)
			got, err := tt.v.Process(tt.args.inFund)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validator.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Validator.Process() = %v, want %v", got, tt.want)
			}
		})
	}
}
