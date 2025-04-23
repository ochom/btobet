package btobet

import (
	"fmt"
	"net/http"

	"github.com/ochom/gutils/env"
	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
)

// PlaceBet ...
func PlaceBet(betSlip BetSlipRequest) (*BetSlipResponse, error) {
	accessToken := env.Get("BTOBET_ACCESS_TOKEN")
	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	betSlip.Mobile = BtoMobile(betSlip.Mobile)
	res, err := gttp.Post(placeBetURL, headers, helpers.ToBytes(betSlip))
	if err != nil {
		logs.Error("error placing bet [%s]=> %s", betSlip.Mobile, err.Error())
		return nil, fmt.Errorf("http err : %v", err)
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error placing bet [%s]=> %s", betSlip.Mobile, string(res.Body))
		return nil, fmt.Errorf("http status err: %v, %s", res.StatusCode, string(res.Body))
	}

	data := helpers.FromBytes[BetSlipResponse](res.Body)
	return &data, nil
}

// CheckBetSlip ...
func CheckBetSlip(mobile, slipID string) (*BetStatusResponse, error) {
	accessToken := env.Get("BTOBET_ACCESS_TOKEN")
	headers := map[string]string{
		"X-API-Key":    accessToken,
		"Content-Type": "application/json",
	}

	mobile = BtoMobile(mobile)

	url := fmt.Sprintf(checkSlipURL, mobile, slipID)
	res, err := gttp.Get(url, headers)
	if err != nil {
		logs.Error("error checking bet slip [%s]=> %s", mobile, err.Error())
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		logs.Error("error checking bet slip [%s]=> %s", mobile, string(res.Body))
		return nil, fmt.Errorf("request failed status %v", res.StatusCode)
	}

	data := helpers.FromBytes[BetStatusResponse](res.Body)
	return &data, nil
}

// GetMarkets ...
func GetMarkets(eventCode string) (*MarketResponse, error) {
	accessToken := env.Get("BTOBET_ACCESS_TOKEN")
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

	if res.StatusCode != http.StatusOK {
		logs.Error("error getting markets [%s]=> %s", eventCode, string(res.Body))
		return nil, fmt.Errorf("request failed status %v", res.StatusCode)
	}

	data := helpers.FromBytes[MarketResponse](res.Body)
	return &data, nil
}
