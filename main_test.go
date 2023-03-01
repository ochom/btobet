package btobet_test

import (
	"testing"

	"github.com/ochom/btobet"
)

func Test_impl_AddPaymentAccount(t *testing.T) {
	s, err := btobet.New()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
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
				mobile: "0708113456",
			},
			wantErr: false,
		},
		{
			name: "07XXXXXXXX",
			args: args{
				mobile: "07XXXXXXXX",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.AddPaymentAccount(tt.args.mobile); (err != nil) != tt.wantErr {
				t.Errorf("impl.AddPaymentAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_impl_GetCustomerDetails(t *testing.T) {
	s, err := btobet.New()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
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
				mobile: "0708113456",
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetCustomerDetails(tt.args.mobile)
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
	s, err := btobet.New()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		eventCode string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "success got markets",
			args: args{
				eventCode: "3617",
			},
			wantNil: false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetMarkets(tt.args.eventCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("impl.GetMarkets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (got == nil) != tt.wantNil {
				t.Errorf("impl.GetMarkets() = %v, wantNil %v", got, tt.wantNil)
			}

		})
	}
}
