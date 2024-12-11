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
func GetCustomerDetails(phone string) (*CustomerDetails, error) {
	payload := map[string]string{
		"apiKey":      env.Get("PAYMENTS_API_KEY"),
		"phoneNumber": BtoMobile(phone),
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", payload["apiKey"]),
		"Content-Type":  "application/json",
	}

	res, err := gttp.NewClient(gttp.GoFiber).Post(getCustomerDetailsURL, headers, payload)
	if err != nil {
		logs.Error("error getting customer details [%s]=> %s", payload["phoneNumber"], err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error getting customer details [%s]=> %s", payload["phoneNumber"], string(res.Body))
		return nil, fmt.Errorf("http status err: %v", res.StatusCode)
	}

	data := helpers.FromBytes[CustomerDetails](res.Body)
	return &data, nil
}

// GetCustomerMetadata ...
func GetCustomerMetadata(phone string) (map[string]any, error) {
	payload := map[string]string{
		"apiKey":      env.Get("PAYMENTS_API_KEY"),
		"phoneNumber": BtoMobile(phone),
	}

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Basic %s", payload["apiKey"]),
		"Content-Type":  "application/json",
	}

	res, err := gttp.NewClient(gttp.GoFiber).Post(getCustomerDetailsURL, headers, payload)
	if err != nil {
		logs.Error("error getting customer metadata [%s]=> %s", payload["phoneNumber"], err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error getting customer metadata [%s]=> %s", payload["phoneNumber"], string(res.Body))
		return nil, fmt.Errorf("http status err: %v", res.StatusCode)
	}

	data := helpers.FromBytes[map[string]any](res.Body)
	return data, nil
}
