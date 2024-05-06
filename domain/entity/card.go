package entity

import (
	"fmt"
	"regexp"
	"time"
)

const (
	verificationCodeTotalLength = 3
	numberTotalLength           = 16
	ownerNameMinLength          = 8
)

type Card struct {
	ownerName        string
	verificationCode string
	lastDigits       string
	validDate        time.Time
}

func NewCard(ownerName, verificationCode, number string, validDate time.Time) (*Card, error) {
	c := &Card{}
	if validDate.UTC().Before(time.Now().UTC()) {
		return nil, fmt.Errorf("card is expired")
	}
	c.validDate = validDate.UTC().Truncate(time.Millisecond)
	if len(ownerName) < ownerNameMinLength {
		return nil, fmt.Errorf("invalid ownerName length. Expected at least %d, got %d", ownerNameMinLength, len(ownerName))
	}
	c.ownerName = ownerName
	r, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return nil, err
	}
	cleanedVerificationCode := r.ReplaceAllString(verificationCode, "")
	if len(cleanedVerificationCode) != verificationCodeTotalLength {
		return nil, fmt.Errorf("invalid verificationCode length. Expected %d, got %d", verificationCodeTotalLength, len(cleanedVerificationCode))
	}
	c.verificationCode = cleanedVerificationCode
	cleanedCardNumber := r.ReplaceAllString(number, "")
	if len(cleanedCardNumber) != numberTotalLength {
		return nil, fmt.Errorf("invalid cardNumber length. Expected %d, got %d", numberTotalLength, len(cleanedCardNumber))
	}
	c.lastDigits = cleanedCardNumber[len(cleanedCardNumber)-4:]

	return c, nil
}

func RestoreCard(ownerName, verificationCode, number string, validDate time.Time) *Card {
	return &Card{
		ownerName:        ownerName,
		verificationCode: verificationCode,
		lastDigits:       number,
		validDate:        validDate,
	}
}

func (c *Card) GetOwnerName() string {
	return c.ownerName
}

func (c *Card) GetVerificationCode() string {
	return c.verificationCode
}

func (c *Card) GetLastDigits() string {
	return c.lastDigits
}

func (c *Card) GetValidDate() time.Time {
	return c.validDate
}
