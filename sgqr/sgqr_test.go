package sgqr

import (
	"testing"
)

func TestGeneratePaynowQrString(t *testing.T) {
	options := PayNowQROptions{
		MobileNumber:    "+6581010321",
		Editable:        false,
		Expiry:          "20251228",
		Amount:          "10.50",
		MerchantName:    "sgqr_test",
		ReferenceNumber: "REF123",
	}

	result, err := GeneratePayNowQrString(options)
	if err != nil {
		t.Errorf("Error generating PayNow QR string: %v\n", err)
	}
	t.Logf("PayNow QR String: %s\n", result)
}

func TestGenerateSGQRString(t *testing.T) {
	options := SGQROptions{
		ReceiverType:             "uen",
		MobileOrUENAccountNumber: "T11LL1111C",
		Editable:                 false,
		Expiry:                   "20251228",
		Amount:                   "10.50",
		SGQRID:                   "SGQR1234567890",
		MerchantName:             "sgqr_test",
		ReferenceNumber:          "REF123",
	}

	result, err := GenerateSGQRString(options)
	if err != nil {
		t.Errorf("Error generating SGQR string: %v\n", err)
	}
	t.Logf("SGQR String: %s\n", result)
}
