package models

type TrTransactionFifo struct{
	TransFifoKey              uint64    `db:"trans_fifo_key"         json:"trans_fifo_key"`
	TransRedKey              *uint64    `db:"trans_red_key"          json:"trans_red_key"`
	TransSubKey              *uint64    `db:"trans_sub_key"          json:"trans_sub_key"`
	SubAcaKey                *uint64    `db:"sub_aca_key"            json:"sub_aca_key"`
	HoldingDays              *uint64    `db:"holding_days"           json:"holding_days"`
	TransUnit                 float32   `db:"trans_unit"             json:"trans_unit"`
	FeeNavMode                uint64    `db:"fee_nav_mode"           json:"fee_nav_mode"`
	TransAmount               float32   `db:"trans_amount"           json:"trans_amount"`
	TransFeeAmount            float32   `db:"trans_fee_amount"       json:"trans_fee_amount"`
	TransFeeTax               float32   `db:"trans_fee_tax"          json:"trans_fee_tax"`
	TransNettAmount           float32   `db:"trans_nett_amount"      json:"trans_nett_amount"`
}