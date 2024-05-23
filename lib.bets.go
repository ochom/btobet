package btobet

import (
	"fmt"
	"net/http"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// PlaceBet ...
func PlaceBet(betSlip BetSlipRequest) (*BetSlipResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN")
	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	mobile, err := parseMobile(betSlip.Mobile)
	if err != nil {
		logs.Error("PlaceBet: error parsing mobile: %s", err.Error())
		return nil, err
	}

	betSlip.Mobile = mobile

	logs.Info("placing bet [%s]=> %s", betSlip.Mobile, string(helpers.ToJSON(betSlip)))

	res, err := gttp.Post(placeBetURL, headers, betSlip)
	if err != nil {
		logs.Error("error placing bet [%s]=> %s", betSlip.Mobile, err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.Status != http.StatusOK {
		logs.Error("error placing bet [%s]=> %s", betSlip.Mobile, string(res.Body))
		return nil, fmt.Errorf("http status err: %v, %s", res.Status, string(res.Body))
	}

	data := helpers.FromJSON[BetSlipResponse](res.Body)
	return &data, nil
}

// CheckBetSlip ...
func CheckBetSlip(mobile, slipID string) (*BetStatusResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN")
	mobile, err := parseMobile(mobile)
	if err != nil {
		logs.Error("CheckBetSlip: error parsing mobile: %s", err.Error())
		return nil, err
	}

	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	logs.Info("checking bet slip [%s]=> %s", mobile, slipID)

	url := fmt.Sprintf(checkSlipURL, mobile, slipID)
	res, err := gttp.Get(url, headers)
	if err != nil {
		logs.Error("error checking bet slip [%s]=> %s", mobile, err.Error())
		return nil, err
	}

	if res.Status != http.StatusOK {
		logs.Error("error checking bet slip [%s]=> %s", mobile, string(res.Body))
		return nil, fmt.Errorf("request failed status %v", res.Status)
	}

	data := helpers.FromJSON[BetStatusResponse](res.Body)
	return &data, nil
}

// GetMarkets ...
func GetMarkets(eventCode string) (*MarketResponse, error) {
	accessToken := helpers.GetEnv("BTOBET_ACCESS_TOKEN")
	headers := map[string]string{
		"X-API-Key": accessToken,
		"Accept":    "application/json",
	}

	logs.Info("getting markets [%s]=> %s", eventCode, accessToken)

	url := fmt.Sprintf(getMarketsURL, eventCode)
	res, err := gttp.Get(url, headers)
	if err != nil {
		logs.Error("error getting markets [%s]=> %s", eventCode, err.Error())
		return nil, err
	}

	if res.Status != http.StatusOK {
		logs.Error("error getting markets [%s]=> %s", eventCode, string(res.Body))
		return nil, fmt.Errorf("request failed status %v", res.Status)
	}

	data := helpers.FromJSON[MarketResponse](res.Body)
	return &data, nil
}
