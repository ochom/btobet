package btobet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gohttp "github.com/ochom/go-http"
)

//BtoBet contains methods for bto bet
type BtoBet interface {
	RegisterUser(ctx context.Context, mobile, password string) (*RegistrationResponse, error)
	CustomerLogin(ctx context.Context, loginRequest LoginRequest) (*LoginResponse, error)
	GetCustomerDetails(ctx context.Context, mobile string) (*CustomerDetails, error)
	AddPaymentAccount(ctx context.Context, mobile string) error
	WithdrawFromWallet(ctx context.Context, mobile, callbackURL string, amount int) error
	PlaceBet(ctx context.Context, betSlip BetSlipRequest) (*BetSlipResponse, error)
	CheckBetSlip(ctx context.Context, mobile, slipID string) (*BetStatusResponse, error)
}

// New ...
func New(timeout time.Duration, pu, pp, pa, bi string, pmi int) BtoBet {
	return &impl{
		http:            gohttp.New(timeout),
		paymentUsername: pu,
		paymentPassword: pp,
		paymentAPIKey:   pa,
		btobetID:        bi,
		paymentMethodID: pmi,
	}
}

type impl struct {
	http            gohttp.Service
	paymentUsername string
	paymentPassword string
	paymentAPIKey   string
	btobetID        string
	paymentMethodID int
}

// RegisterUser ...
func (s *impl) RegisterUser(ctx context.Context, mobile, password string) (*RegistrationResponse, error) {

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
		"apiKey":          s.paymentAPIKey,
		"activateAccount": "true",
		"loginAccount":    "false",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", s.paymentAPIKey),
	}

	status, res, err := s.http.Post(ctx, registerCustomerURL, headers, body)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status: %d", status)
	}

	var resp RegistrationResponse
	if err = json.Unmarshal(res, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

//CustomerLogin ...
func (s *impl) CustomerLogin(ctx context.Context, loginRequest LoginRequest) (*LoginResponse, error) {

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", s.paymentAPIKey),
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
		"apiKey":                  s.paymentAPIKey,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %v", err)
	}

	status, res, err := s.http.Post(ctx, loginURL, headers, body)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", status)
	}

	var loginResponse LoginResponse

	err = json.Unmarshal(res, &loginResponse)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &loginResponse, nil
}

//GetCustomerDetails ...
func (s *impl) GetCustomerDetails(ctx context.Context, mobile string) (*CustomerDetails, error) {

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", s.paymentAPIKey),
		"Content-Type":  "application/json",
	}

	data := map[string]string{
		"apiKey":      s.paymentAPIKey,
		"phoneNumber": mobile,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %v", err)
	}

	status, res, err := s.http.Post(ctx, getCustomerDetailsURL, headers, body)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v", status)
	}

	var customerDetails CustomerDetails

	err = json.Unmarshal(res, &customerDetails)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}

	return &customerDetails, nil
}

// AddPaymentAccount ...
func (s *impl) AddPaymentAccount(ctx context.Context, mobile string) error {
	customer, err := s.GetCustomerDetails(ctx, mobile)
	if err != nil {
		return err
	}

	if !customer.IsSuccessful {
		return fmt.Errorf("customer not registered: %s", customer.Errors[0].Description)
	}

	data := map[string]any{
		"apiKey":     s.paymentAPIKey,
		"internalID": customer.Customer.Account.InternalID,
		"paymentAccounts": []map[string]any{
			{
				"AccountReference": mobile,
				"HolderName":       mobile,
				"PaymentMethodID":  s.paymentMethodID,
			},
		},
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", s.paymentAPIKey),
		"Content-Type":  "application/json",
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	status, res, err := s.http.Post(ctx, addPaymentAccountURL, headers, payload)
	if err != nil {
		return err
	}

	if status != http.StatusOK {
		log.Println(string(res))
		return fmt.Errorf("adding payment account failed status: %d", status)
	}

	return nil
}

//WithdrawFromWallet ...
func (s *impl) WithdrawFromWallet(ctx context.Context, mobile, callbackURL string, amount int) error {
	err := s.AddPaymentAccount(ctx, mobile)
	if err != nil {
		return err
	}

	apiKey := Encode(fmt.Sprintf("%s:%s", s.paymentUsername, s.paymentPassword))

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

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal err: %v", err)
	}

	status, res, err := s.http.Post(ctx, withdrawURL, headers, body)
	if err != nil {
		return fmt.Errorf("http err : %v", err.Error())
	}

	if status != http.StatusOK {
		return fmt.Errorf("withdrawal failed status: %d error: %s", status, string(res))
	}

	return nil
}

//PlaceBet ...
func (s *impl) PlaceBet(ctx context.Context, betSlip BetSlipRequest) (*BetSlipResponse, error) {
	headers := map[string]string{
		"X-API-Key":    s.btobetID,
		"Content-Type": "application/json",
	}

	body, err := json.Marshal(betSlip)
	if err != nil {
		return nil, fmt.Errorf("json marshal err: %v", err)
	}

	status, res, err := s.http.Post(ctx, placeBetURL, headers, body)
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("http status err: %v, %s", status, string(res))
	}

	var response BetSlipResponse

	if err = json.Unmarshal(res, &response); err != nil {
		return nil, fmt.Errorf("json unmarshal err: %v", err)
	}
	return &response, nil
}

//CheckBetSlip ...
func (s *impl) CheckBetSlip(ctx context.Context, mobile, slipID string) (*BetStatusResponse, error) {
	headers := map[string]string{
		"X-API-Key":    s.btobetID,
		"Content-Type": "application/json",
	}

	url := fmt.Sprintf(checkSlipURL, mobile, slipID)
	status, res, err := s.http.Get(ctx, url, headers)
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("request failed status %v", status)
	}

	var data BetStatusResponse
	if err = json.Unmarshal(res, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
