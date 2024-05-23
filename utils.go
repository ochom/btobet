package btobet

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

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

	mobile = strings.TrimPrefix(mobile, "254")
	mobile = "0" + mobile
	return mobile, nil
}
