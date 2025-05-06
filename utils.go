package btobet

import (
	"encoding/base64"
	"strings"
	"time"

	"github.com/ochom/gutils/helpers"
	"github.com/ochom/gutils/logs"
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

func BtoMobile(s string) string {
	mobile := helpers.ParseMobile(s)
	if mobile == "" {
		logs.Error("invalid mobile number")
		return ""
	}

	mobile = strings.TrimPrefix(mobile, "254")
	mobile = "0" + mobile
	return mobile
}
