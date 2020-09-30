package io

import (
	"context"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/milind/velocitylimit/config"
	validate "github.com/milind/velocitylimit/validate"
)

func TestProcess(t *testing.T) {
	config.Configuration.Currency = "$"
	config.Configuration.DayLimit = 5000.00
	config.Configuration.WeekLimit = 20000.00
	config.Configuration.MaxAttemptsPerDay = 3
	config.AppendMode = false
	validatr, _ := validate.NewValidator("memory")

	type args struct {
		ctx       context.Context
		input     Ingester
		output    Sink
		validator *validate.Validator
	}
	tests := []struct {
		name           string
		args           args
		input          string
		expectedOutput string
	}{
		{
			name:           "1",
			args:           args{ctx: context.Background(), validator: validatr},
			input:          "{\"id\":\"75887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}",
			expectedOutput: "{\"id\":\"75887\",\"customer_id\":\"528\",\"accepted\":true}",
		},
		{
			name:           "2",
			args:           args{ctx: context.Background(), validator: validatr},
			input:          "{\"id\":\"75886\",\"customer_id\":\"3528\",\"load_amount\":\"#3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}\n{\"id\":\"75888\",\"customer_id\":\"1528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}",
			expectedOutput: "{\"id\":\"75888\",\"customer_id\":\"1528\",\"accepted\":true}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = ioutil.WriteFile("testIn.txt", []byte(tt.input), 0644)
			tt.args.input, _ = NewInputFile("testIn.txt")
			tt.args.output, _ = NewOutputFile("testOut.txt")
			Process(tt.args.ctx, tt.args.input, tt.args.output, tt.args.validator)
			data, _ := ioutil.ReadFile("testOut.txt")
			if strings.TrimSpace(string(data)) != tt.expectedOutput {
				t.Errorf("Process() = %v, want %v", string(data), tt.expectedOutput)
			}
		})
	}
}

func Test_preProcess(t *testing.T) {
	tm, _ := time.Parse(time.RFC3339, "2000-02-20T00:00:00Z")
	config.Configuration.Currency = "$"
	type args struct {
		fund string
	}
	tests := []struct {
		name    string
		args    args
		want    *validate.Deposit
		wantErr bool
	}{
		{
			name: "valid",
			args: args{fund: "{\"id\":\"75887\",\"customer_id\":\"528\",\"load_amount\":\"$3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want: &validate.Deposit{
				ID:         75887,
				CustomerID: 528,
				LoadAmount: 3318.47,
				Time:       tm,
			},
			wantErr: false,
		},
		{
			name:    "invalid_currency",
			args:    args{fund: "{\"id\":\"75887\",\"customer_id\":\"528\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "string_custid",
			args:    args{fund: "{\"id\":\"75887\",\"customer_id\":\"dff\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative__custid",
			args:    args{fund: "{\"id\":\"75887\",\"customer_id\":\"-528\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "string_loadid",
			args:    args{fund: "{\"id\":\"jkj\",\"customer_id\":\"528\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative__loadid",
			args:    args{fund: "{\"id\":\"-75887\",\"customer_id\":\"528\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid_json",
			args:    args{fund: "{\"id\":\"-75887\",customer_id\":\"528\",\"load_amount\":\"%3318.47\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invald_loadamnt",
			args:    args{fund: "{\"id\":\"-75887\",\"customer_id\":\"528\",\"load_amount\":\"%dfg\",\"time\":\"2000-02-20T00:00:00Z\"}"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := preProcess(tt.args.fund)
			if (err != nil) != tt.wantErr {
				t.Errorf("preProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postProcess(t *testing.T) {
	type args struct {
		fundRslt *validate.DepositStatus
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{fundRslt: &validate.DepositStatus{ID: 123, CustomerID: 456, Accepted: true}},
			want:    "{\"id\":\"123\",\"customer_id\":\"456\",\"accepted\":true}",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := postProcess(tt.args.fundRslt)
			if (err != nil) != tt.wantErr {
				t.Errorf("postProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("postProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
