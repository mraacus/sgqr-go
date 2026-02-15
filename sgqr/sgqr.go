package sgqr

import (
	"fmt"
	"time"
)

func getMerchantAccountInformationPaynow(sgqrOptions SGQROptions) SGQRDataObject {
	// 02-26 reserved for card networks - Table 4.1
	// 26-51 reserved for additional payment networks - Table 4.2
	payNowIndicator := SGQRDataObject{
		ID:        "00",
		Name:      "PayNow Indicator",
		MaxLength: 32,
		// This is how paynow is indicated in the SGQR spec
		// TODO: Discover the indicators for other payment schemes
		Value: "SG.PAYNOW",
	}
	mobileOrUenAccount := SGQRDataObject{
		ID:        "01",
		Name:      "Mobile Or UEN Account",
		MaxLength: 1,
	}
	if sgqrOptions.ReceiverType == "mobile" {
		mobileOrUenAccount.Value = "0"
	} else {
		mobileOrUenAccount.Value = "2"
	}
	mobileOrUENAccountNumber := SGQRDataObject{
		ID:        "02",
		MaxLength: 13,
		Name:      "Mobile or UEN Account Number",
		Value:     sgqrOptions.MobileOrUENAccountNumber,
	}
	editable := SGQRDataObject{
		ID:        "03",
		MaxLength: 1,
		Name:      "Payment amount editable",
	}
	if sgqrOptions.Editable {
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

	if len(sgqrOptions.Expiry) == 8 {
		expiry := SGQRDataObject{
			ID:        "04",
			Name:      "Expiry Date",
			MaxLength: 8,
			Value:     sgqrOptions.Expiry,
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
 * for the SGQR data objects specification
 */
func getSGQRObject(sgqrOptions SGQROptions) SGQRRootObject {

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
	merchantAccountInformation := getMerchantAccountInformationPaynow(sgqrOptions)

	/** ID: "51", SGQR ID - Unique SGQR idenfiier determined by the Repository
	 *  SGQR Centralised Repository generates the SGQR and keeps the records
	 *  identifier for the SGQR, used to identify the QR Code
	 *  PAINPOINT: We can help with the application flow for SGQR codes
	 *  Mandatory
	 */
	sgqrID := SGQRDataObject{
		ID:        "51",
		Name:      "SGQR ID",
		MaxLength: 99,
		// Placeholder - This is an object, refer to SGQR spec v1.7
		Value: sgqrOptions.SGQRID,
	}

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
		Value:     sgqrOptions.Amount,
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
		Value:     sgqrOptions.MerchantName,
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
		Value:     sgqrOptions.ReferenceNumber,
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

	sgqrDataObjects := []SGQRDataObject{
		payloadFormatIndicator,
		pointOfInitiationMethod,
		merchantAccountInformation,
		sgqrID,
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
		DataObjects: sgqrDataObjects,
	}
}

/**
 *	Converts a SGQRDataObject into its string representation
 *	- ID - 2 characters
 *	- Length - 2 characters
 *	- Value - Length of value in characters
 *  Eg: ID: "53", Length: "03", Value: "702" -> "5303702"
 *  If the value is a list of SGQRDataObjects, calls getString() on each sub object and concatenates the results
 */
func (obj SGQRDataObject) getString() (string, error) {
	valueStr := ""
	switch obj.Value.(type) {
	case string:
		valueStr = obj.Value.(string)
		if len(valueStr) > obj.MaxLength {
			return "", fmt.Errorf("id %s value %s is out of max length %d", obj.ID, valueStr, obj.MaxLength)
		}
	case []SGQRDataObject:
		subObjects := obj.Value.([]SGQRDataObject)
		for _, subObj := range subObjects {
			subObjStr, err := subObj.getString()
			if err != nil {
				return "", err
			}
			valueStr += subObjStr
		}
	}
	lengthStr := fmt.Sprintf("%02d", len(valueStr))
	return obj.ID + lengthStr + valueStr, nil
}

/**
 *	Generates a CRC string for the SGQRDataObject
 *	- ID - 2 characters
 *	- Length - 2 characters (always "04" for CRC)
 *	- Value - 4 characters (CRC16 checksum of the value)
 */
func (obj SGQRDataObject) getCRCString(value string) string {
	value += obj.ID + "04" // CRC require length "04"
	data := []byte(value)
	checkSum := crc16(data)
	valueStr := fmt.Sprintf("%04X", checkSum)
	return obj.ID + "04" + valueStr
}

/**
 *	Generates the full SGQR string representation
 *	- Iterates through each SGQRDataObject in the root object
 *	- Calls getString() on each object
 *	- If the object ID is "63", calls getCRCString() with the current value string
 *	- Concatenates the results into a single string
 */
func (root *SGQRRootObject) getString() (string, error) {
	valueStr := ""
	for _, dataObject := range root.DataObjects {
		if dataObject.ID == "63" {
			valueStr += dataObject.getCRCString(valueStr)
		} else {
			dataObjectStr, err := dataObject.getString()
			if err != nil {
				return "", err
			}
			valueStr += dataObjectStr
		}
	}
	return valueStr, nil
}

func GenerateSGQRString(sgqrOptions SGQROptions) (string, error) {
	if err := validateSGQROptions(sgqrOptions); err != nil {
		return "", err
	}

	// Get the SGQR data object
	sgqrRootObject := getSGQRObject(sgqrOptions)

	// Generate the SGQR string
	return sgqrRootObject.getString()
}
