package validate

import (
	"context"
	"testing"

	dal "github.com/milind/velocitylimit/validate/internal"
)

func Test_weeklyLimitsRule_Do(t *testing.T) {
	type args struct {
		ctx     context.Context
		deposit *Deposit
	}
	tests := []struct {
		name    string
		da      *weeklyLimitsRule
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "within-weekly-limit-1",
			da:      &weeklyLimitsRule{weekLimit: 100.00, attempts: 3, dal: dal.NewMemoryDataStore()},
			args:    args{ctx: nil, deposit: &Deposit{LoadAmount: 200.23}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.da.Do(tt.args.ctx, tt.args.deposit)
			if (err != nil) != tt.wantErr {
				t.Errorf("weeklyLimitsRule.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("weeklyLimitsRule.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
