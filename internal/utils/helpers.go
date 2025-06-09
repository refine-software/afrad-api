package utils

import (
	"crypto/rand"
	"encoding/json"
	"regexp"
	"time"
)

var digits = []byte("1234567890")

func GenerateRandomOTP() string {
	const otpLength = 6
	otp := make([]byte, otpLength)

	for i := range otpLength {
		randomByte := make([]byte, 1)

		for {
			_, err := rand.Read(randomByte)
			if err == nil {
				break
			}
		}

		otp[i] = digits[randomByte[0]%byte(len(digits))]
	}

	return string(otp)
}

// This function checks for an iraqi phone number,
// Starts with 07 and the 9 numbers after it are digits between 0 - 9
func ValidPhoneNumber(phone string) bool {
	// this regex checks the following
	// 1- the phone number starts with 07.
	// 2- the next 9 letters are number digits 0 - 9.
	// the ^ indicates the start of the string and $ is the end of it
	re := regexp.MustCompile(`^07[0-9]{9}$`)
	return re.MatchString(phone)
}

func CheckValidJSON(jsonStr string) bool {
	var jsonData map[string]any
	err := json.Unmarshal([]byte(jsonStr), &jsonData)
	if err != nil {
		return false
	}
	// Check if the JSON object is empty
	return len(jsonData) > 0
}

// Get the expiration time after certain number of days added to now().
func GetExpTimeAfterDays(numOfDays int) time.Time {
	return time.Now().Add((time.Hour * 24) * time.Duration(numOfDays))
}

// Get the expiration time after certain number of Minutes added to now().
func GetExpTimeAfterMins(numOfMins int) time.Time {
	return time.Now().Add((time.Minute * time.Duration(numOfMins)))
}
