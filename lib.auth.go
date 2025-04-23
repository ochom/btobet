package btobet

import (
	"fmt"
	"net/http"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// RegisterUser ...
func RegisterUser(phone, password string) (*RegistrationResponse, error) {
	mobile := BtoMobile(phone)

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
		"apiKey":          env.Get("PAYMENTS_API_KEY"),
		"activateAccount": "true",
		"loginAccount":    "false",
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Basic %s", payload["apiKey"]),
	}

	logs.Info("registering user [%s]=> %s", mobile, string(helpers.ToBytes(payload)))

	res, err := gttp.Post(registerCustomerURL, headers, helpers.ToBytes(payload))
	if err != nil {
		logs.Error("error registering user [%s]=> %s", mobile, err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error registering user [%s]=> %s", mobile, string(res.Body))
		return nil, fmt.Errorf("http status: %d", res.StatusCode)
	}

	data := helpers.FromBytes[RegistrationResponse](res.Body)
	return &data, nil
}

// CustomerLogin ...
func CustomerLogin(loginRequest LoginRequest) (*LoginResponse, error) {
	payload := map[string]string{
		"login":                   BtoMobile(loginRequest.Username),
		"password":                loginRequest.Password,
		"ipAddress":               loginRequest.IPaddress,
		"returnBalance":           "true",
		"returnApplicableBonuses": "true",
		"returnCustomerDetails":   "true",
		"deviceType":              "Default",
		"apiKey":                  env.Get("PAYMENTS_API_KEY"),
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", payload["apiKey"]),
		"Content-Type":  "application/json",
	}

	res, err := gttp.Post(loginURL, headers, helpers.ToBytes(payload))
	if err != nil {
		logs.Error("error logging in user [%s]=> %s", payload["login"], err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error logging in user [%s]=> %s", payload["login"], string(res.Body))
		return nil, fmt.Errorf("http status err: %v", res.StatusCode)
	}

	data := helpers.FromBytes[LoginResponse](res.Body)
	return &data, nil
}
