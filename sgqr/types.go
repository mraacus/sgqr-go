package sgqr

type ReceiverType string

const (
	ReceiverTypeMobile ReceiverType = "mobile"
	ReceiverTypeUEN    ReceiverType = "uen"
)

// SGQROptions represents the options for generating a SGQR string
//   - Fields are mandatory unless stated optional. Passing empty or invalid values for mandatory fields will return an error
//   - You may pass empty strings "" for optional fields
type SGQROptions struct {
	ReceiverType             ReceiverType // mobile or UEN receiver
	MobileOrUENAccountNumber string	 	  // +65[mobile number] eg. "+6599998888" for mobile numbers of UEN Number for company entities
	Editable                 bool   	  // Indicates if the payment amount is editable
	Expiry                   string  	  // "YYYYMMDD" eg. "20240101"
	Amount                   string  	  // Amount as a numeric string with 2 decimal places eg. "10.40"
	SGQRID                   string  	  // Unique identifier for the SGQR, provided by the SGQR Centralised Repository
	MerchantName             string  	  // Optional - Name of the company or individual recipient
	ReferenceNumber          string  	  // Optional - Reference number for the transaction
}

// PayNowQROptions represents the options for generating a QR string for mobile number PayNow
//   - Fields are mandatory unless stated optional. Passing empty or invalid values for mandatory fields will return an error
//   - You may pass empty strings "" for optional fields
type PayNowQROptions struct {
	MobileNumber    string	// "+65[mobile number]"" eg. "+6599998888" for mobile numbers of UEN Number for company entities eg. "T11LL1111C"
	Editable        bool	// Indicates if the payment amount is editable
	Expiry          string  // "YYYYMMDD" eg. "20240101"
	Amount          string  // Amount as a numeric string with 2 decimal places eg. "10.40"
	MerchantName    string  // Optional - Name of the company or individual recipient
	ReferenceNumber string  // Optional - Reference number for the transaction
}

type SGQRDataObject struct {
	ID        string
	Name      string
	MaxLength int
	Value     interface{}
}

type SGQRRootObject struct {
	DataObjects []SGQRDataObject
}
