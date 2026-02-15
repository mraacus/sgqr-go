package sgqr

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

func validateSGQROptions(sgqrOptions SGQROptions) error {
	if sgqrOptions.ReceiverType == ReceiverTypeMobile {
		if err := validateMobileNumberString(sgqrOptions.MobileOrUENAccountNumber); err != nil {
			return err
		}
	}
	if sgqrOptions.ReceiverType == ReceiverTypeUEN {
		if err := validateUENString(sgqrOptions.MobileOrUENAccountNumber); err != nil {
			return err
		}
	}
	if err := validateExpiryString(sgqrOptions.Expiry); err != nil {
		return err
	}
	if err := validateAmountString(sgqrOptions.Amount); err != nil {
		return err
	}
	if err := validateSGQRString(sgqrOptions.SGQRID); err != nil {
		return err
	}
	return nil
}

func validatePayNowQROptions(payNowQROptions PayNowQROptions) error {
	if err := validateMobileNumberString(payNowQROptions.MobileNumber); err != nil {
		return err
	}
	if err := validateExpiryString(payNowQROptions.Expiry); err != nil {
		return err
	}
	if err := validateAmountString(payNowQROptions.Amount); err != nil {
		return err
	}
	return nil
}

func validateMobileNumberString(s string) error {
	if s == "" {
		return fmt.Errorf("non empty mobile number is required")
	}
	if s[0] != '+' {
		return fmt.Errorf("mobile numbers must start with a country code. eg. '+6599998888'")
	}
	if !isNumericString(s[1:]) {
		return fmt.Errorf("mobile numbers must be numeric exclusding country code definition. eg. '+6599998888'")
	}
	return nil
}

func validateUENString(s string) error {
	if s == "" {
		return fmt.Errorf("non empty UEN is required")
	}
	if len(s) != 10 && len(s) != 9 {
		return fmt.Errorf("UEN must be 9 or 10 characters long")
	}
	if !isAlphanumericString(s) {
		return fmt.Errorf("UEN must be alphanumeric: %v", s)
	}
	return nil
}

func validateExpiryString(s string) error {
	if s == "" {
		return fmt.Errorf("non empty expiry date is required")
	}
	if _, err := time.Parse("20060102", s); err != nil {
		return fmt.Errorf("expiry must be a valid date in the format YYYYMMDD: %v", err)
	}
	return nil
}

func validateAmountString(s string) error {
	if s == "" {
		return fmt.Errorf("non empty amount is required")
	}
	if !isFloatString(s) {
		return fmt.Errorf("amount must be a valid float: %v", s)
	}
	return nil
}

func validateSGQRString(s string) error {
	if s == "" {
		return fmt.Errorf("non empty SGQR ID is required")
	}
	if s[:4] != "SGQR" {
		return fmt.Errorf("valid SGQR ID starts with 'SGQR': %v", s)
	}
	if !isNumericString(s[4:]) {
		return fmt.Errorf("valid SGQR ID must be numeric excluding 'SGQR' prefix: %v", s)
	}
	return nil
}

func isNumericString(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func isFloatString(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isAlphanumericString(s string) bool {
	for _, c := range s {
        if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
            return false
        }
	}
	return true
}