package validate

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/milind/velocitylimit/config"
	dal "github.com/milind/velocitylimit/validate/internal"
)

func Test_dailyLimitsRule_Do(t *testing.T) {
	type args struct {
		ctx     context.Context
		deposit *Deposit
	}
	tests := []struct {
		name    string
		da      *dailyLimitsRule
		args    args
		want    bool
		wantErr bool
		given   []*dal.Transaction
	}{
		{
			name:    "txn-amount-bigger-than-dailylimit",
			da:      &dailyLimitsRule{dayLimit: 199.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 200.23, CustomerID: 123, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    false,
			wantErr: false,
			given:   nil,
		},
		{
			name:    "total-txn-per-day-one-cent-bigger-than-dailylimit",
			da:      &dailyLimitsRule{dayLimit: 1000, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 100.01, CustomerID: 123, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    false,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 2890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "total-txn-per-day-one-cent-less-than-dailylimit",
			da:      &dailyLimitsRule{dayLimit: 1000.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 199.99, CustomerID: 1231, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 1231, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 2890, LoadAmount: 100.00, Time: time.Date(2000, time.February, 14, 21, 10, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 3890, LoadAmount: 100.00, Time: time.Date(2000, time.February, 14, 21, 20, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "total-txn-per-day-equal-to-dailylimit-fourth-attempt",
			da:      &dailyLimitsRule{dayLimit: 1000.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 200, CustomerID: 1231, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 1231, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 2890, LoadAmount: 100.00, Time: time.Date(2000, time.February, 14, 21, 10, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 3890, LoadAmount: 100.00, Time: time.Date(2000, time.February, 14, 21, 20, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "total-txn-per-day-equal-to-dailylimit-third-attempt",
			da:      &dailyLimitsRule{dayLimit: 1000.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 99.99, CustomerID: 1231, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 1231, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 2890, LoadAmount: 300.01, Time: time.Date(2000, time.February, 14, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "negative-is-a-valid-txn",
			da:      &dailyLimitsRule{dayLimit: 1000.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: -99.99, CustomerID: 1231, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 1231, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 1231, ID: 2890, LoadAmount: 400, Time: time.Date(2000, time.February, 14, 21, 10, 0, 0, time.UTC), Status: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupDataStore(tt.da.dal, tt.given)
			got, err := tt.da.Do(tt.args.ctx, tt.args.deposit)
			if (err != nil) != tt.wantErr {
				t.Errorf("dailyLimitsRule.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dailyLimitsRule.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupDataStore(mds DataStore, txns []*dal.Transaction) {
	dal.NewMemoryDataStore().ClearData()
	for _, t := range txns {
		_ = mds.SaveCustomerTxn(t)
	}
}

func TestNewDailyLimitsRule(t *testing.T) {
	type args struct {
		cfg config.Config
		ds  DataStore
	}
	mds := dal.NewMemoryDataStore()
	tests := []struct {
		name string
		args args
		want *dailyLimitsRule
	}{
		{
			name: "test1",
			args: args{ds: mds, cfg: config.Config{
				VelocityLimitConfig: config.VelocityLimitConfig{
					DayLimit:          100,
					MaxAttemptsPerDay: 3,
				},
			}},
			want: &dailyLimitsRule{dal: mds, dayLimit: 100, attempts: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDailyLimitsRule(tt.args.cfg, tt.args.ds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDailyLimitsRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dailyLimitsRule_String(t *testing.T) {
	tests := []struct {
		name string
		da   *dailyLimitsRule
		want string
	}{
		{
			name: "1",
			da:   &dailyLimitsRule{},
			want: "DailyLimitsRule",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.da.String(); got != tt.want {
				t.Errorf("dailyLimitsRule.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
