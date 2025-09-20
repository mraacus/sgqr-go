package sgqr

import (
	"fmt"
	"time"
)

func getMerchantAccountInformationPaynowQR(payNowQROptions PayNowQROptions) SGQRDataObject {
	// 02-26 reserved for card networks - Table 4.1
	// 26-51 reserved for additional payment networks - Table 4.2
	payNowIndicator := SGQRDataObject{
		ID:        "00",
		Name:      "PayNow Indicator",
		MaxLength: 32,
		// This is how paynow is indicated in the SGQR spec
		Value: "SG.PAYNOW",
	}
	mobileOrUenAccount := SGQRDataObject{
		ID:        "01",
		Name:      "Mobile Or UEN Account",
		MaxLength: 1,
		Value:     "0",
	}
	mobileOrUENAccountNumber := SGQRDataObject{
		ID:        "02",
		MaxLength: 13,
		Name:      "Mobile or UEN Account Number",
		Value:     payNowQROptions.MobileNumber,
	}
	editable := SGQRDataObject{
		ID:        "03",
		MaxLength: 1,
		Name:      "Payment amount editable",
	}
	if payNowQROptions.Editable {
		editable.Value = "1"
	} else {
		editable.Value = "0"
	}

	merchantAccountInformationValue := []SGQRDataObject{
		payNowIndicator,
		mobileOrUenAccount,
		mobileOrUENAccountNumber,
		editable,
	}

	if len(payNowQROptions.Expiry) == 8 {
		expiry := SGQRDataObject{
			ID:        "04",
			Name:      "Expiry Date",
			MaxLength: 8,
			Value:     payNowQROptions.Expiry,
		}
		merchantAccountInformationValue = append(merchantAccountInformationValue, expiry)
	} else {
		// Default to 1 hour expiry if not provided
		currentTime := time.Now()
		expiry := SGQRDataObject{
			ID:        "04",
			Name:      "Expiry Date",
			MaxLength: 8,
			Value:     currentTime.Add(time.Hour * 1).Format("20060102"),
		}
		merchantAccountInformationValue = append(merchantAccountInformationValue, expiry)
	}

	return SGQRDataObject{
		ID:        "26",
		Name:      "Merchant Account Information - PayNow",
		MaxLength: 99,
		Value:     merchantAccountInformationValue,
	}
}

/**
 * Refer to https://github.com/hisenyuan/EMVCo-SGQR-encode-decode-crc/blob/master/src/main/java/com/hisen/emvco/docs/EMVCo-Merchant-Presented-QR-Specification-v1_0_uqpay.pdf
 * for the PayNow QR data objects specification
 */
func getSGQRRootObject(payNowQROptions PayNowQROptions) SGQRRootObject {

	// ID: "00", Payload Format Indicator - Defines the version of QR Code and conventions
	// Mandatory
	payloadFormatIndicator := SGQRDataObject{
		ID:        "00",
		Name:      "Payload Format Indicator",
		MaxLength: 2,
		// 01 is the default qr code format version convention
		Value: "01",
	}

	// ID: "01", Point of Initiation Method - Identifies the communication technology and static vs dynamic
	// Optional
	pointOfInitiationMethod := SGQRDataObject{
		ID:        "01",
		Name:      "Point of Initiation Method",
		MaxLength: 2,
		// 11 for same QR for more than 1 transaction
		// 12 for single use QR
		Value: "12",
	}

	/** ID: "02"-"50", Merchant Account Information - Identifies the merchants
	 *  "26"-"50" represent the merchant's registered payment schemes (paynow, paylah, grabpay... etc)
	 *  To scale, we can add additional payment schems based on the payload.
	 *  Mandatory
	 */
	merchantAccountInformation := getMerchantAccountInformationPaynowQR(payNowQROptions)

	// ID: "51", SGQR ID - Unique SGQR identifier - Omitted for PayNow QR

	// ID: "52", Merchant Category Code - MCC of the merchant
	// Mandatory
	merchantCategoryCode := SGQRDataObject{
		ID:        "52",
		Name:      "Merchant Category Code",
		MaxLength: 4,
		// If this is not utilised by a payment scheme, default to “0000”
		// Merchants that require MCC - (Ref [B] - ISO 18245. Retail financial services – Merchant category codes)
		Value: "0000",
	}

	// ID: "53", Transaction Currency - Currency of the transaction
	// Mandatory
	transactionCurrency := SGQRDataObject{
		ID:        "53",
		Name:      "Transaction Currency",
		MaxLength: 3,
		// Default to SGD
		Value: "702", // 702 is the ISO 4217 code for SGD
	}

	// ID: "54", Transaction Amount - Amount of the transaction
	// Conditional - Absent if the amount is input by the customer on the application
	transactionAmount := SGQRDataObject{
		ID:        "54",
		Name:      "Transaction Amount",
		MaxLength: 13,
		Value:     payNowQROptions.Amount,
	}

	// ID: "55"-"57" - Tip/Convenience Fee - Optional

	// ID: "58", Country Code - Country code of the merchant
	// Mandatory
	countryCode := SGQRDataObject{
		ID:        "58",
		Name:      "Country Code",
		MaxLength: 2,
		Value:     "SG", // Default to Singapore - SG
	}

	// ID: "59", Merchant Name - Name of the merchant
	// Mandatory
	merchantName := SGQRDataObject{
		ID:        "59",
		Name:      "Merchant Name",
		MaxLength: 25,
		Value:     payNowQROptions.CompanyName,
	}

	// ID: "60", Merchant City - City of the merchant
	// Mandatory
	merchantCity := SGQRDataObject{
		ID:        "60",
		Name:      "Merchant City",
		MaxLength: 15,
		Value:     "Singapore", // Default to Singapore
	}

	// ID: "61", Postal Code - Optional (For relevant payment systems)

	// ID: "62", Additional Data Fields
	// Refer to Section 4.8 in the SGQR spec for more details
	// Optional
	referenceNumber := SGQRDataObject{
		ID:        "01",
		Name:      "Reference Number",
		MaxLength: 25,
		Value:     payNowQROptions.ReferenceNumber,
	}
	additionalDataFields := SGQRDataObject{
		ID:        "62",
		Name:      "Additional Data Fields",
		MaxLength: 99,
		Value:     []SGQRDataObject{referenceNumber},
	}

	// ID: "63", CRC16 - Cyclic Redundancy Check checksum for all the SGQR data objects in the SGQR
	// Used for integrity checking of the data
	// Mandatory
	crc := SGQRDataObject{
		ID:        "63",
		Name:      "CRC",
		MaxLength: 4,
		Value:     "",
	}

	SGQRDataObjects := []SGQRDataObject{
		payloadFormatIndicator,
		pointOfInitiationMethod,
		merchantAccountInformation,
		merchantCategoryCode,
		transactionCurrency,
		transactionAmount,
		countryCode,
		merchantName,
		merchantCity,
		additionalDataFields,
		crc,
	}

	return SGQRRootObject{
		DataObjects: SGQRDataObjects,
	}
}

func GeneratePayNowQrString(payNowQROptions PayNowQROptions) (string, error) {
	if err := validatePayNowQROptions(payNowQROptions); err != nil {
		return "", fmt.Errorf("invalid PayNow QR options: %v", err)
	}

	// Get the PayNow QR data object
	SGQRRootObject := getSGQRRootObject(payNowQROptions)

	// Generate the SGQR string
	return SGQRRootObject.getString()
}

func validatePayNowQROptions(payNowQROptions PayNowQROptions) error {
	if payNowQROptions.MobileNumber == "" {
		return fmt.Errorf("mobile number is required")
	}
	if payNowQROptions.Expiry != "" {
		if _, err := time.Parse("20060102", payNowQROptions.Expiry); err != nil {
			return fmt.Errorf("expiry must be a valid date in the format YYYYMMDD: %v", err)
		}
	}
	if payNowQROptions.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	return nil
}
