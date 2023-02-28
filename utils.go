package btobet

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"
)

// TimeZone ...
var TimeZone = "Africa/Nairobi"

// Encode ...
func Encode(rawString string) string {
	return base64.StdEncoding.EncodeToString([]byte(rawString))
}

// GetLocation returns timezone in Nairobi
func GetLocation() *time.Location {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil
	}

	return loc
}

// GetEnv ...
func GetEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf("Environment variable %s not set", key)
}

// GetIntEnv ...
func GetIntEnv(key string) (int, error) {
	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return num, nil
	}
	return 0, fmt.Errorf("Environment variable %s not set", key)
}
