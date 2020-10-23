package models

type TransactionFormatSinvest struct {
	SubscriptionRedeemption *[]SubscriptionRedeemption `json:"subscription_redeemption"`
	SwitchOutSwitchIn       *[]SwitchOutSwitchIn       `json:"switchout_switchin"`
}

type SubscriptionRedeemption struct {
	TransactionDate             string `json:"transaction_date"`
	TransactionType             uint64 `json:"transaction_type"`
	SACode                      string `json:"sa_code"`
	InvestorFundUnitACNo        string `json:"investor_fund_unit"`
	FundCode                    string `json:"fund_code"`
	AmountNominal               string `json:"amount_nominal"`
	AmountUnit                  string `json:"amount_unit"`
	AmountAllUnit               string `json:"amount_all_unit"`
	FeeNominal                  string `json:"fee_nominal"`
	FeeUnit                     string `json:"fee_unit"`
	FeePercent                  string `json:"fee_percent"`
	RedmPaymentACSequenceCode   string `json:"redm_payment_ac_sequence_code"`
	RedmPaymentBankBICCode      string `json:"redm_payment_bank_bic_code"`
	RedmPaymentBankBIMemberCode string `json:"redm_payment_bank_bi_member_code"`
	RedmPaymentACCode           string `json:"redm_payment_ac_code"`
	PaymentDate                 string `json:"payment_date"`
	TransferType                string `json:"transfer_type"`
	SaReferenceNo               uint64 `json:"sa_reference_no"`
}

type SwitchOutSwitchIn struct {
	TransactionDate        string  `json:"transaction_date"`
	TransactionType        uint64  `json:"transaction_type"`
	SACode                 string  `json:"sa_code"`
	InvestorFundUnitACNo   string  `json:"investor_fund_unit"`
	SwitchOutFundCode      string  `json:"switch_out_fund_code"`
	SwitchOutAmountNominal float32 `json:"switch_out_amount_nominal"`
	SwitchOutAmountUnit    float32 `json:"switch_out_amount_unit"`
	SwitchOutAmountAll     string  `json:"switch_out_amount_all"`
	SwitchingFeeChargeFund string  `json:"switching_fee_charge_fund"`
	FeeNominal             float32 `json:"fee_nominal"`
	FeeUnit                float32 `json:"fee_unit"`
	FeePercent             float32 `json:"fee_percent"`
	SwitchInFundCode       string  `json:"switch_in_fund_code"`
	PaymentDate            string  `json:"payment_date"`
	TransferType           string  `json:"transfer_type"`
	SaReferenceNo          uint64  `json:"sa_reference_no"`
}
