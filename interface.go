package btobet

import "context"

//BtoBet contains methods for bto bet
type BtoBet interface {
	RegisterUser(ctx context.Context, mobile, password string) (*RegistrationResponse, error)
	CustomerLogin(ctx context.Context, loginRequest LoginRequest) (*LoginResponse, error)
	GetCustomerDetails(ctx context.Context, mobile string) (*CustomerDetails, error)
	AddPaymentAccount(ctx context.Context, mobile string) error
	WithdrawFromWallet(ctx context.Context, mobile, slipID, callbackURL string, amount int) error
	PlaceBet(ctx context.Context, betSlip BetSlipRequest) (*BetSlipResponse, error)
	CheckBetSlip(ctx context.Context, mobile, slipID string) (*BetStatusResponse, error)
}
