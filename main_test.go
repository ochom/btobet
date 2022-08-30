package btobet

import (
	"context"
	"os"
	"testing"
	"time"

	gohttp "github.com/ochom/go-http"
)

func Test_impl_AddPaymentAccount(t *testing.T) {
	s := &impl{
		http:            gohttp.New(time.Second * 30),
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
		http:            gohttp.New(time.Second * 30),
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
