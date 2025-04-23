package btobet

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// AddPaymentAccount ...
func AddPaymentAccount(mobile string) error {
	mobile = BtoMobile(mobile)
	customer, err := GetCustomerDetails(mobile)
	if err != nil {
		logs.Error("AddPaymentAccount: error getting customer details: %s", err.Error())
		return err
	}

	if !customer.IsSuccessful {
		logs.Error("AddPaymentAccount: customer not registered: %s", customer.Errors[0].Description)
		return fmt.Errorf("customer not registered: %s", customer.Errors[0].Description)
	}

	payload := map[string]any{
		"apiKey":     env.Get("PAYMENTS_API_KEY"),
		"internalID": customer.Customer.Account.InternalID,
		"paymentAccounts": []map[string]any{
			{
				"AccountReference": mobile,
				"HolderName":       mobile,
				"PaymentMethodID":  env.Int("PAYMENT_METHOD_ID", 0),
			},
		},
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", payload["apiKey"]),
		"Content-Type":  "application/json",
	}

	logs.Info("adding payment account [%s]=> %s", mobile, string(helpers.ToBytes(payload)))
	res, err := gttp.Post(addPaymentAccountURL, headers, helpers.ToBytes(payload))
	if err != nil {
		logs.Error("AddPaymentAccount: error adding payment account: %s", err.Error())
		return err
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("AddPaymentAccount: error adding payment account: %s", string(res.Body))
		return fmt.Errorf("adding payment account failed status: %d", res.StatusCode)
	}

	return nil
}

// WithdrawFromWallet ...
func WithdrawFromWallet(mobile, callbackURL string, amount int) error {
	mobile = BtoMobile(mobile)
	if err := AddPaymentAccount(mobile); err != nil {
		logs.Error("WithdrawFromWallet: error adding payment account: %s", err.Error())
		return err
	}

	paymentUsername := env.Get("PAYMENTS_USERNAME")
	paymentPassword := env.Get("PAYMENTS_PASSWORD")

	apiKey := Encode(fmt.Sprintf("%s:%s", paymentUsername, paymentPassword))

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", apiKey),
		"Content-Type":  "application/json",
	}

	now := time.Now().In(GetLocation()).Format("20060102150405")
	payload := map[string]any{
		"PspId":        now,
		"OrderId":      now,
		"Currency":     "KES",
		"WithdrawalId": nil,
		"Amount":       amount,
		"Username":     mobile,
		"PosId":        2331007,
		"CashierId":    "1",
		"CallbackURL":  callbackURL,
	}

	logs.Info("withdrawing from wallet [%s]=> %s", mobile, string(helpers.ToBytes(payload)))
	res, err := gttp.Post(withdrawURL, headers, helpers.ToBytes(payload))
	if err != nil {
		logs.Error("WithdrawFromWallet: error withdrawing from wallet: %s", err.Error())
		return fmt.Errorf("http err : %v", err.Error())
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("WithdrawFromWallet: error withdrawing from wallet: %s", string(res.Body))
		return fmt.Errorf("withdrawal failed status: %d error: %s", res.StatusCode, string(res.Body))
	}

	return nil
}
