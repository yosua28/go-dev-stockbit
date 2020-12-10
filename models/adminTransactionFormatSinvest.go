package models

import (
	"api/db"
	"log"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type TransactionFormatSinvest struct {
	// SubscriptionRedeemption *[]SubscriptionRedeemption `json:"subscription_redeemption"`
	// SwitchTransaction       *[]SwitchTransaction       `json:"switch"`
	SubscriptionRedeemption *[]OaRequestCsvFormatFiksTxt `json:"subscription_redeemption"`
	SwitchTransaction       *[]OaRequestCsvFormatFiksTxt `json:"switch"`
}

type SubscriptionRedeemption struct {
	TransactionDate             string `json:"transaction_date"`
	TransactionType             string `json:"transaction_type"`
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
	SaReferenceNo               string `json:"sa_reference_no"`
}

type SwitchTransaction struct {
	TransactionDate        string `json:"transaction_date"`
	TransactionType        string `json:"transaction_type"`
	SACode                 string `json:"sa_code"`
	InvestorFundUnitACNo   string `json:"investor_fund_unit"`
	SwitchOutFundCode      string `json:"switch_out_fund_code"`
	SwitchOutAmountNominal string `json:"switch_out_amount_nominal"`
	SwitchOutAmountUnit    string `json:"switch_out_amount_unit"`
	SwitchOutAmountAll     string `json:"switch_out_amount_all"`
	SwitchingFeeChargeFund string `json:"switching_fee_charge_fund"`
	FeeNominal             string `json:"fee_nominal"`
	FeeUnit                string `json:"fee_unit"`
	FeePercent             string `json:"fee_percent"`
	SwitchInFundCode       string `json:"switch_in_fund_code"`
	PaymentDate            string `json:"payment_date"`
	TransferType           string `json:"transfer_type"`
	SaReferenceNo          string `json:"sa_reference_no"`
}

type DataTransactionParent struct {
	TransactionKey  uint64  `db:"transaction_key"           json:"transaction_key"`
	SinvestFundCode *string `db:"sinvest_fund_code"       json:"sinvest_fund_code"`
	TransAmount     decimal.Decimal `db:"trans_amount"              json:"trans_amount"`
	TransUnit       decimal.Decimal `db:"trans_unit"                json:"trans_unit"`
	FlagRedemtAll   *uint8  `db:"flag_redempt_all"          json:"flag_redempt_all"`
}

func GetDataParentTransactionSwitch(c *[]DataTransactionParent, value []string) (int, error) {
	inQuery := strings.Join(value, ",")
	query2 := `SELECT 
			tr.transaction_key AS transaction_key,
			mp.sinvest_fund_code AS sinvest_fund_code,
			tr.trans_amount AS trans_amount,
			tr.trans_unit AS trans_unit,
			tr.flag_redempt_all AS flag_redempt_all 
			FROM tr_transaction AS tr 
			INNER JOIN ms_product AS mp ON tr.product_key = mp.product_key`
	query := query2 + " WHERE tr.transaction_key IN(" + inQuery + ")"

	// Main query
	log.Println(query)
	err := db.Db.Select(c, query)
	if err != nil {
		log.Println(err)
		return http.StatusBadGateway, err
	}

	return http.StatusOK, nil
}
