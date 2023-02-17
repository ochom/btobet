package btobet

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func Test_impl_AddPaymentAccount(t *testing.T) {
	s := &impl{
		paymentUsername: os.Getenv("PAYMENTS_USERNAME"),
		paymentPassword: os.Getenv("PAYMENTS_PASSWORD"),
		paymentAPIKey:   os.Getenv("PAYMENTS_API_KEY"),
		btobetID:        os.Getenv("BTOBET_ACCESS_TOKEN"),
	}

	type args struct {
		ctx    context.Context
		mobile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "0708113456",
			args: args{
				ctx:    context.Background(),
				mobile: "0708113456",
			},
			wantErr: false,
		},
		{
			name: "07XXXXXXXX",
			args: args{
				ctx:    context.Background(),
				mobile: "07XXXXXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.AddPaymentAccount(tt.args.ctx, tt.args.mobile); (err != nil) != tt.wantErr {
				t.Errorf("impl.AddPaymentAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_impl_GetCustomerDetails(t *testing.T) {
	s := &impl{
		paymentUsername: os.Getenv("PAYMENTS_USERNAME"),
		paymentPassword: os.Getenv("PAYMENTS_PASSWORD"),
		paymentAPIKey:   os.Getenv("PAYMENTS_API_KEY"),
		btobetID:        os.Getenv("BTOBET_ACCESS_TOKEN"),
	}
	type args struct {
		ctx    context.Context
		mobile string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "0708113456",
			args: args{
				ctx:    context.Background(),
				mobile: "0708113456",
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCustomerDetails(tt.args.ctx, tt.args.mobile)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.GetCustomerDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("impl.GetCustomerDetails() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func Test_impl_GetMarkets(t *testing.T) {
	s := &impl{
		paymentUsername: os.Getenv("PAYMENTS_USERNAME"),
		paymentPassword: os.Getenv("PAYMENTS_PASSWORD"),
		paymentAPIKey:   os.Getenv("PAYMENTS_API_KEY"),
		btobetID:        os.Getenv("BTOBET_ACCESS_TOKEN"),
	}
	type args struct {
		ctx        context.Context
		mobile     string
		eventCode  string
		marketCode string
	}
	tests := []struct {
		name    string
		args    args
		want    *MarketResponse
		wantErr bool
	}{
		{
			name: "success got markets",
			args: args{
				ctx:        context.Background(),
				mobile:     "0743119767",
				eventCode:  "3617",
				marketCode: "1",
			},
			want: &MarketResponse{
				EventCode:       3617,
				ResponseCode:    0,
				ResponseMessage: "Success",
				Outcomes: []struct {
					Outcome   string  `json:"Outcome,omitempty"`
					Odd       float32 `json:"Odd,omitempty"`
					ShortCode string  `json:"ShortCode,omitempty"`
				}{
					{
						Outcome:   "1",
						Odd:       1.38,
						ShortCode: "1",
					},
					{
						Outcome:   "2",
						Odd:       6.25,
						ShortCode: "2",
					},
					{
						Outcome:   "X",
						Odd:       4.4,
						ShortCode: "X",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetMarkets(tt.args.ctx, tt.args.mobile, tt.args.eventCode, tt.args.marketCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.GetMarkets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("impl.GetMarkets() = %v, want %v", got, tt.want)
			}
		})
	}
}
