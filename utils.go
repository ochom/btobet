package btobet

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ochom/gutils/gttp"
	"github.com/ochom/gutils/helpers"
)

// TimeZone ...
var TimeZone = "Africa/Nairobi"

// Encode ...
func Encode(rawString string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawString))
}

// GetLocation returns time zone in Nairobi
func GetLocation() *time.Location {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil
	}

	return loc
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

	// print as json
	b, err := json.MarshalIndent(printable, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("json err : %v", err)
	}

	fmt.Printf("%s\n", string(b))
	res, err := gttp.NewRequest(url, headers, payload).Post()
	if err != nil {
		return nil, fmt.Errorf("http err : %v", err)
	}

	return res, nil
}
