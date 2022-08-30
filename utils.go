package btobet

import (
	"encoding/base64"
	"time"
)

// TimeZone ...
var TimeZone = "Africa/Nairobi"

// Encode ...
func Encode(rawString string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawString))
}

//GetLocation returns timezone in Nairobi
func GetLocation() *time.Location {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil
	}

	return loc
}
