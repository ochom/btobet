package btobet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
)

// Controller ...
type Controller struct {
	paymentUsername string
	paymentPassword string
	paymentAPIKey   string
	btobetID        string
	paymentMethodID int
}

// New ...
func New() (*Controller, error) {
	return &Controller{
		paymentUsername: helpers.GetEnv("PAYMENTS_USERNAME", ""),
		paymentPassword: helpers.GetEnv("PAYMENTS_PASSWORD", ""),
		paymentAPIKey:   helpers.GetEnv("PAYMENTS_API_KEY", ""),
		btobetID:        helpers.GetEnv("BTOBET_ACCESS_TOKEN", ""),
		paymentMethodID: helpers.GetEnvInt("PAYMENT_METHOD_ID", 0),
	}, nil
}

func parseMobile(s string) (string, error) {
	mobile := helpers.ParseMobile(s)
	if mobile == "" {
		return "", fmt.Errorf("invalid mobile number")
	}

	// replace 254 with 0
	mobile = fmt.Sprintf("0%s", mobile[3:])
	return mobile, nil
}

// wrapRequest ...
func wrapRequest(url string, headers map[string]string, payload any) (*gttp.Response, error) {
	printable := map[string]interface{}{
		"headers": headers,
		"payload": payload,
		"url":     url,
	}

	fmt.Printf("register user: %+v\n", printable)
	res, err := gttp.NewRequest(url, headers, payload).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	return res, nil
}

// RegisterUser ...
func (c *Controller) RegisterUser(mobile, password string) (*RegistrationResponse, error) {
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
		"apiKey":          c.paymentAPIKey,
		"activateAccount": "true",
		"loginAccount":    "false",
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
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
func (c *Controller) CustomerLogin(loginRequest LoginRequest) (*LoginResponse, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
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
		"returnBalance":           "false",
		"returnApplicableBonuses": "false",
		"returnCustomerDetails":   "false",
		"deviceType":              "Default",
		"apiKey":                  c.paymentAPIKey,
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
func (c *Controller) GetCustomerDetails(mobile string) (*CustomerDetails, error) {
	mobile, err := parseMobile(mobile)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
		"Content-Type":  "application/json",
	}

	payload := map[string]string{
		"apiKey":      c.paymentAPIKey,
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
func (c *Controller) AddPaymentAccount(mobile string) error {
	mobile, err := parseMobile(mobile)
	if err != nil {
		return err
	}

	customer, err := c.GetCustomerDetails(mobile)
	if err != nil {
		return err
	}

	if !customer.IsSuccessful {
		return fmt.Errorf("customer not registered: %s", customer.Errors[0].Description)
	}

	payload := map[string]any{
		"apiKey":     c.paymentAPIKey,
		"internalID": customer.Customer.Account.InternalID,
		"paymentAccounts": []map[string]any{
			{
				"AccountReference": mobile,
				"HolderName":       mobile,
				"PaymentMethodID":  c.paymentMethodID,
			},
		},
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
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
func (c *Controller) WithdrawFromWallet(mobile, callbackURL string, amount int) error {
	mobile, err := parseMobile(mobile)
	if err != nil {
		return err
	}

	if err := c.AddPaymentAccount(mobile); err != nil {
		return err
	}

	apiKey := Encode(fmt.Sprintf("%s:%s", c.paymentUsername, c.paymentPassword))

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
func (c *Controller) PlaceBet(betSlip BetSlipRequest) (*BetSlipResponse, error) {
	headers := map[string]string{
		"X-API-Key":    c.btobetID,
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
func (c *Controller) CheckBetSlip(mobile, slipID string) (*BetStatusResponse, error) {
	mobile, err := parseMobile(mobile)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-API-Key":    c.btobetID,
		"Content-Type": "application/json",
	}

	url := fmt.Sprintf(checkSlipURL, mobile, slipID)
	res, err := gttp.NewRequest(url, headers, nil).Get()
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
func (c *Controller) GetMarkets(eventCode string) (*MarketResponse, error) {
	headers := map[string]string{
		"X-API-Key": c.btobetID,
		"Accept":    "application/json",
	}

	url := fmt.Sprintf(getMarketsURL, eventCode)
	res, err := gttp.NewRequest(url, headers, nil).Get()
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
