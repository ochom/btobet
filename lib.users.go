package btobet

import (
	"fmt"
	"net/http"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// GetCustomerDetails ...
func GetCustomerDetails(mobile string) (*CustomerDetails, error) {
	paymentAPIKey := env.Get("PAYMENTS_API_KEY")

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

	logs.Error("getting customer details [%s]=> %s", mobile, string(helpers.ToBytes(payload)))
	res, err := gttp.Post(getCustomerDetailsURL, headers, payload)
	if err != nil {
		logs.Error("error getting customer details [%s]=> %s", mobile, err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		logs.Error("error getting customer details [%s]=> %s", mobile, string(res.Body))
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	data := helpers.FromBytes[CustomerDetails](res.Body)
	return &data, nil
}

// GetCustomerMetadata ...
func GetCustomerMetadata(mobile string) (map[string]any, error) {
	paymentAPIKey := env.Get("PAYMENTS_API_KEY")

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

	logs.Info("getting customer metadata [%s]=> %s", mobile, string(helpers.ToBytes(payload)))

	res, err := gttp.Post(getCustomerDetailsURL, headers, payload)
	if err != nil {
		logs.Error("error getting customer metadata [%s]=> %s", mobile, err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		logs.Error("error getting customer metadata [%s]=> %s", mobile, string(res.Body))
		return nil, fmt.Errorf("http status err: %v", res.Status)
	}

	data := helpers.FromBytes[map[string]any](res.Body)
	return data, nil
}
