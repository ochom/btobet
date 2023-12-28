package btobet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
)

// RegisterUser ...
func RegisterUser(mobile, password string) (*RegistrationResponse, error) {
	paymentAPIKey := helpers.GetEnv("PAYMENTS_API_KEY", "")

	mobile, err := parseMobile(mobile)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"customer": map[string]interface{}{
			"PreferredNotificationType": 1,
			"CustomerV3": map[string]interface{}{
				"CustomerDetails": map[string]interface{}{
					"FirstName":               "Kwikbet",
					"LastName":                "Kwikbet",
					"Email":                   fmt.Sprintf("%s@kwikbet.co.ke", mobile),
					"Username":                mobile,
					"PhoneNumber":             mobile,
					"MobileNumber":            mobile,
					"City":                    "Nairobi",
					"Postcode":                "00100",
					"Address":                 "Kwikbet",
					"Gender":                  "Male",
					"LanguageISO":             "EN",
					"CountryISO":              "KE",
					"CurrencyISO":             "KES",
					"Password":                password,
					"DateOfBirth":             "1990-01-01",
					"IPAddress":               "",
					"Browser":                 "Chrome",
					"CivilIdentificationCode": "123132132113112",
					"Note":                    "",
					"EmploymentStatus":        0,
					"Longitude":               nil,
					"Latitude":                nil,
					"TimeZoneName":            "SA Pacific Standard Time",
					"IsTestCustomer":          "false",
					"PassportNumber":          "",
					"Profession":              "",
				},
			},
		},
		"deviceType":      "Default",
		"apiKey":          paymentAPIKey,
		"activateAccount": "true",
		"loginAccount":    "false",
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", paymentAPIKey),
	}

	res, err := wrapRequest(registerCustomerURL, headers, payload)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", res.Status)
	}

	var data RegistrationResponse
	if err = json.Unmarshal(res.Body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// CustomerLogin ...
func CustomerLogin(loginRequest LoginRequest) (*LoginResponse, error) {
	paymentAPIKey := helpers.GetEnv("PAYMENTS_API_KEY", "")

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", paymentAPIKey),
		"Content-Type":  "application/json",
	}

	mobile, err := parseMobile(loginRequest.Username)
	if err != nil {
		return nil, err
	}

	payload := map[string]string{
		"login":                   mobile,
		"password":                loginRequest.Password,
		"ipAddress":               loginRequest.IPaddress,
		"returnBalance":           "true",
		"returnApplicableBonuses": "true",
		"returnCustomerDetails":   "true",
		"deviceType":              "Default",
		"apiKey":                  paymentAPIKey,
	}

	res, err := wrapRequest(loginURL, headers, payload)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	var data LoginResponse
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &data, nil
}

// GetCustomerDetails ...
func GetCustomerDetails(mobile string) (*CustomerDetails, error) {
	paymentAPIKey := helpers.GetEnv("PAYMENTS_API_KEY", "")

	mobile, err := parseMobile(mobile)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", paymentAPIKey),
		"Content-Type":  "application/json",
	}

	payload := map[string]string{
		"apiKey":      paymentAPIKey,
		"phoneNumber": mobile,
	}

	res, err := wrapRequest(getCustomerDetailsURL, headers, payload)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	var data CustomerDetails
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &data, nil
}

// AddPaymentAccount ...
func AddPaymentAccount(mobile string) error {
	paymentAPIKey := helpers.GetEnv("PAYMENTS_API_KEY", "")
	paymentMethodID := helpers.GetEnvInt("PAYMENT_METHOD_ID", 0)
	mobile, err := parseMobile(mobile)
	if err != nil {
		return err
	}

	customer, err := GetCustomerDetails(mobile)
	if err != nil {
		return err
	}

	if !customer.IsSuccessful {
		return fmt.Errorf("customer not registered: %s", customer.Errors[0].Description)
	}

	payload := map[string]any{
		"apiKey":     paymentAPIKey,
		"internalID": customer.Customer.Account.InternalID,
		"paymentAccounts": []map[string]any{
			{
				"AccountReference": mobile,
				"HolderName":       mobile,
				"PaymentMethodID":  paymentMethodID,
			},
		},
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", paymentAPIKey),
		"Content-Type":  "application/json",
	}

	res, err := wrapRequest(addPaymentAccountURL, headers, payload)
	if err != nil {
		return err
	}

	if res.Status != http.StatusOK {
		return fmt.Errorf("adding payment account failed status: %d", res.Status)
	}

	return nil
}

// WithdrawFromWallet ...
func WithdrawFromWallet(mobile, callbackURL string, amount int) error {
	mobile, err := parseMobile(mobile)
	if err != nil {
		return err
	}

	if err := AddPaymentAccount(mobile); err != nil {
		return err
	}

	paymentUsername := helpers.GetEnv("PAYMENTS_USERNAME", "")
	paymentPassword := helpers.GetEnv("PAYMENTS_PASSWORD", "")

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

	res, err := wrapRequest(withdrawURL, headers, payload)
	if err != nil {
		return fmt.Errorf("http err : %v", err.Error())
	}

	if res.Status != http.StatusOK {
		return fmt.Errorf("withdrawal failed status: %d error: %s", res.Status, string(res.Body))
	}

	return nil
}

// PlaceBet ...
func PlaceBet(betSlip BetSlipRequest) (*BetSlipResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN", "")
	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	mobile, err := parseMobile(betSlip.Mobile)
	if err != nil {
		return nil, err
	}

	betSlip.Mobile = mobile

	res, err := wrapRequest(withdrawURL, headers, betSlip)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v, %s", res.Status, string(res.Body))
	}

	var data BetSlipResponse
	if err := json.Unmarshal(res.Body, &data); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return &data, nil
}

// CheckBetSlip ...
func CheckBetSlip(mobile, slipID string) (*BetStatusResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN", "")
	mobile, err := parseMobile(mobile)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	url := fmt.Sprintf(checkSlipURL, mobile, slipID)
	res, err := gttp.Get(url, headers)
	if err != nil {
		return nil, err
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed status %v", res.Status)
	}

	var data BetStatusResponse
	if err = json.Unmarshal(res.Body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetMarkets ...
func GetMarkets(eventCode string) (*MarketResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN", "")
	headers := map[string]string{
		"X-API-Key": accessToken,
		"Accept":    "application/json",
	}

	url := fmt.Sprintf(getMarketsURL, eventCode)
	res, err := gttp.Get(url, headers)
	if err != nil {
		return nil, err
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed status %v", res.Status)
	}

	var data MarketResponse
	if err = json.Unmarshal(res.Body, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
