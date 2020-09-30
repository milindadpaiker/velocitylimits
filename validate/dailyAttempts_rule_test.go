package validate

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/milind/velocitylimit/config"
	dal "github.com/milind/velocitylimit/validate/internal"
)

func TestNewDailyAttemptsRule(t *testing.T) {
	type args struct {
		cfg config.Config
		ds  DataStore
	}
	mds := dal.NewMemoryDataStore()
	tests := []struct {
		name string
		args args
		want *dailyAttemptsRule
	}{
		{
			name: "test1",
			args: args{ds: mds, cfg: config.Config{
				VelocityLimitConfig: config.VelocityLimitConfig{
					DayLimit:          100,
					MaxAttemptsPerDay: 3,
					WeekLimit:         100,
				},
			}},
			want: &dailyAttemptsRule{dal: mds, attempts: 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDailyAttemptsRule(tt.args.cfg, tt.args.ds); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDailyAttemptsRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dailyAttemptsRule_Do(t *testing.T) {
	type args struct {
		ctx     context.Context
		deposit *Deposit
	}
	tests := []struct {
		name    string
		da      *dailyAttemptsRule
		args    args
		want    bool
		wantErr bool
		given   []*dal.Transaction
	}{
		{
			name:    "total-txn-per-day-one-cent-bigger-than-dailylimit",
			da:      &dailyAttemptsRule{attempts: 3, dal: dal.NewMemoryDataStore()},
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
			name:    "total-txn-per-day-one-cent-bigger-than-dailylimit",
			da:      &dailyAttemptsRule{attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 100.01, CustomerID: 123, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
			},
		},
		{
			name:    "attemots3",
			da:      &dailyAttemptsRule{attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 100.01, CustomerID: 123, Time: time.Date(2000, time.February, 14, 23, 0, 0, 0, time.UTC)}},
			want:    true,
			wantErr: false,
			given: []*dal.Transaction{
				{CustomerID: 123, ID: 890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 13, 20, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 13, 21, 0, 0, 0, time.UTC), Status: 1},
				{CustomerID: 123, ID: 1890, LoadAmount: 300.00, Time: time.Date(2000, time.February, 14, 21, 0, 0, 0, time.UTC), Status: 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupDataStore(tt.da.dal, tt.given)
			got, err := tt.da.Do(tt.args.ctx, tt.args.deposit)
			if (err != nil) != tt.wantErr {
				t.Errorf("dailyAttemptsRule.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("dailyAttemptsRule.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dailyAttemptsRule_String(t *testing.T) {
	tests := []struct {
		name string
		da   *dailyAttemptsRule
		want string
	}{
		{
			name: "1",
			da:   &dailyAttemptsRule{},
			want: "DailyAttemptsRule",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.da.String(); got != tt.want {
				t.Errorf("dailyAttemptsRule.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
