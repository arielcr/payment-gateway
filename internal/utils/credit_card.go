package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"regexp"
	"strconv"
)

// TokenizeCreditCard generates a token for a credit card
func TokenizeCreditCard(creditCardNumber string) (string, error) {
	token, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateRandomToken generates a random token for the credit card
func generateRandomToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return token, nil
}

func ValidateCreditCard(cardNumber string) error {
	if !luhnCheck(cardNumber) {
		return errors.New("invalid credit card number")
	}
	return nil
}

func GetCreditCardBrand(cardNumber string) string {
	visaPattern := regexp.MustCompile(`^4[0-9]{12}(?:[0-9]{3})?$`)
	mastercardPattern := regexp.MustCompile(`^5[1-5][0-9]{14}$`)
	amexPattern := regexp.MustCompile(`^3[47][0-9]{13}$`)
	dinersPattern := regexp.MustCompile(`^3(?:0[0-5]|[68][0-9])[0-9]{11}$`)

	if visaPattern.MatchString(cardNumber) {
		return "Visa"
	} else if mastercardPattern.MatchString(cardNumber) {
		return "Mastercard"
	} else if amexPattern.MatchString(cardNumber) {
		return "American Express"
	} else if dinersPattern.MatchString(cardNumber) {
		return "Diners Club"
	}

	return "Unknown"
}

func GetLastFourDigits(cardNumber string) (string, error) {
	if len(cardNumber) < 4 {
		return "", errors.New("credit card number must be at least 4 digits long")
	}
	return cardNumber[len(cardNumber)-4:], nil
}

func luhnCheck(cardNumber string) bool {
	var sum int
	double := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}
		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		double = !double
	}
	return sum%10 == 0
}
