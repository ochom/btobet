package btobet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ochom/gutils/gttp"
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
	paymentUsername, err := GetEnv("PAYMENTS_USERNAME")
	if err != nil {
		return nil, err
	}

	paymentPassword, err := GetEnv("PAYMENTS_PASSWORD")
	if err != nil {
		return nil, err
	}

	paymentAPIKey, err := GetEnv("PAYMENTS_API_KEY")
	if err != nil {
		return nil, err
	}
	btobetID, err := GetEnv("BTOBET_ACCESS_TOKEN")
	if err != nil {
		return nil, err
	}
	paymentMethodID, err := GetIntEnv("PAYMENT_METHOD_ID")
	if err != nil {
		return nil, err
	}

	return &Controller{
		paymentUsername: paymentUsername,
		paymentPassword: paymentPassword,
		paymentAPIKey:   paymentAPIKey,
		btobetID:        btobetID,
		paymentMethodID: paymentMethodID,
	}, nil
}

// RegisterUser ...
func (c *Controller) RegisterUser(mobile, password string) (*RegistrationResponse, error) {

	data := map[string]interface{}{
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

	res, err := gttp.NewRequest(registerCustomerURL, headers, data).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", res.Status)
	}

	var resp RegistrationResponse
	if err = json.Unmarshal(res.Body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// CustomerLogin ...
func (c *Controller) CustomerLogin(loginRequest LoginRequest) (*LoginResponse, error) {

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
		"Content-Type":  "application/json",
	}

	data := map[string]string{
		"login":                   loginRequest.Username,
		"password":                loginRequest.Password,
		"ipAddress":               loginRequest.IPaddress,
		"returnBalance":           "false",
		"returnApplicableBonuses": "false",
		"returnCustomerDetails":   "false",
		"deviceType":              "Default",
		"apiKey":                  c.paymentAPIKey,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %v", err)
	}

	res, err := gttp.NewRequest(loginURL, headers, body).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(res.Body, &loginResponse)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &loginResponse, nil
}

// GetCustomerDetails ...
func (c *Controller) GetCustomerDetails(mobile string) (*CustomerDetails, error) {

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", c.paymentAPIKey),
		"Content-Type":  "application/json",
	}

	data := map[string]string{
		"apiKey":      c.paymentAPIKey,
		"phoneNumber": mobile,
	}

	res, err := gttp.NewRequest(getCustomerDetailsURL, headers, data).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	var customerDetails CustomerDetails

	err = json.Unmarshal(res.Body, &customerDetails)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &customerDetails, nil
}

// AddPaymentAccount ...
func (c *Controller) AddPaymentAccount(mobile string) error {
	customer, err := c.GetCustomerDetails(mobile)
	if err != nil {
		return err
	}

	if !customer.IsSuccessful {
		return fmt.Errorf("customer not registered: %s", customer.Errors[0].Description)
	}

	data := map[string]any{
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

	res, err := gttp.NewRequest(addPaymentAccountURL, headers, data).Post()
	if err != nil {
		return err
	}

	if res.Status != http.StatusOK {
		log.Println(string(res.Body))
		return fmt.Errorf("adding payment account failed status: %d", res.Status)
	}

	return nil
}

// WithdrawFromWallet ...
func (c *Controller) WithdrawFromWallet(mobile, callbackURL string, amount int) error {
	err := c.AddPaymentAccount(mobile)
	if err != nil {
		return err
	}

	apiKey := Encode(fmt.Sprintf("%s:%s", c.paymentUsername, c.paymentPassword))

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", apiKey),
		"Content-Type":  "application/json",
	}

	now := time.Now().In(GetLocation()).Format("20060102150405")

	data := map[string]any{
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

	res, err := gttp.NewRequest(withdrawURL, headers, data).Post()
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

	body, err := json.Marshal(betSlip)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %v", err)
	}

	res, err := gttp.NewRequest(placeBetURL, headers, body).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v, %s", res.Status, string(res.Body))
	}

	var response BetSlipResponse

	if err = json.Unmarshal(res.Body, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return &response, nil
}

// CheckBetSlip ...
func (c *Controller) CheckBetSlip(mobile, slipID string) (*BetStatusResponse, error) {
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
