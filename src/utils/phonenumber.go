package utils

import "github.com/nyaruka/phonenumbers"

func FormatPhoneNumber(phoneNumber string) (string, error) {
	num, err := phonenumbers.Parse(phoneNumber, "DE")
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL), nil
}
