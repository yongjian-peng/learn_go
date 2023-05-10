package haodapayreq

type CallBackReq struct {
	Data       CallBackOrderData `json:"data,omitempty"`
	Event      string            `json:"event,omitempty"`
	Gateway    string            `json:"gateway,omitempty"`
	GatewayUpi string            `json:"gateway_upi,omitempty"`
	Status     string            `json:"status,omitempty"`
}

type CallBackOrderData struct {
	Utr         string `json:"UTR,omitempty"`
	AccountNo   string `json:"account_no,omitempty"`
	Amount      string `json:"amount,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	Ifsc        string `json:"ifsc,omitempty"`
	Name        string `json:"name,omitempty"`
	OrderID     string `json:"order_id,omitempty"`
	PayerUPIID  string `json:"payer_UPIID,omitempty"`
	PaymentMode string `json:"payment_mode,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Reference   string `json:"reference,omitempty"`
	Remarks     string `json:"remarks,omitempty"`
	Status      string `json:"status,omitempty"`
}

type CallBackPayoutReq struct {
	Data   CallBackPayoutData `json:"data,omitempty"`
	Event  string             `json:"event,omitempty"`
	Status string             `json:"status,omitempty"`
}

type CallBackPayoutData struct {
	Utr                      string `json:"UTR,omitempty"`
	Amount                   string `json:"amount,omitempty"`
	BeneficiaryAccountIfsc   string `json:"beneficiary_account_ifsc,omitempty"`
	BeneficiaryAccountName   string `json:"beneficiary_account_name,omitempty"`
	BeneficiaryAccountNumber string `json:"beneficiary_account_number,omitempty"`
	BeneficiaryBankName      string `json:"beneficiary_bank_name,omitempty"`
	BeneficiaryUpiHandle     string `json:"beneficiary_upi_handle,omitempty"`
	CreatedAt                string `json:"created_at,omitempty"`
	PaymentMode              string `json:"payment_mode,omitempty"`
	PayoutID                 string `json:"payout_id,omitempty"`
	Remarks                  string `json:"remarks,omitempty"`
	TransferDate             string `json:"transfer_date,omitempty"`
	Reference                string `json:"reference,omitempty"`
	Checksum                 string `json:"checksum,omitempty"`
}
