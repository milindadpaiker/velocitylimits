package validate

import (
	"context"
	"testing"

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
	}{
		{
			name:    "within-daily-limit-1",
			da:      &dailyLimitsRule{dayLimit: 100.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 200.23}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
