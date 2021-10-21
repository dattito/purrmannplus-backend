package utils

import "github.com/nyaruka/phonenumbers"

// Brings the phone number in a international format
func FormatPhoneNumber(phoneNumber string) (string, error) {
	num, err := phonenumbers.Parse(phoneNumber, "DE")
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL), nil
}
