package models

type OaRequestCsvFormat struct {
	Type                  string `json:"type"`
	SACode                string `json:"sa_code"`
	SID                   string `json:"sid"`
	FirstName             string `json:"first_name"`
	MiddleName            string `json:"middle_name"`
	LastName              string `json:"last_name"`
	CountryOfNationality  string `json:"country_of_nationality"`
	IDno                  string `json:"id_no"`
	IDExpirationDate      string `json:"id_expiration_date"`
	NpwpNo                string `json:"npwp_no"`
	NpwpRegistrationDate  string `json:"npwp_registration_date"`
	CountryOfBirth        string `json:"country_of_birthday"`
	PlaceOfBirth          string `json:"place_of_birth"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	EducationalBackground string `json:"education_background"`
	MotherMaidenName      string `json:"mother_maiden_name"`
	Religion              string `json:"religion"`
	Occupation            string `json:"occupation"`
	IncomeLevel           string `json:"income_level"`
	MaritalStatus         string `json:"marital_status"`
	SpouseName            string `json:"spouse_name"`
	InvestorRiskProfile   string `json:"investor_risk_profile"`
	InvestmentObjective   string `json:"investment_objective"`
	SourceOfFund          string `json:"source_of_fund"`
	AssetOwner            string `json:"asset_owner"`
	KTPAddress            string `json:"ktp_address"`
	KTPCityCode           string `json:"ktp_city_code"`
	KTPPostalCode         string `json:"ktp_postal_code"`

	CorrespondenceAddress    string `json:"correspondence_address"`
	CorrespondenceCityCode   string `json:"correspondence_City_code"`
	CorrespondenceCityName   string `json:"correspondence_city_name"`
	CorrespondencePostalCode string `json:"correspondence_postal_code"`
	CountryOfCorrespondence  string `json:"country_of_correspondence"`

	DomicileAddress              string `json:"domicile_address"`
	DomicileCityCode             string `json:"domicile_city_code"`
	DomicileCityName             string `json:"domicile_city_name"`
	DomicilePostalCode           string `json:"domicile_postal_code"`
	CountryOfDomicile            string `json:"country_of_domicile"`
	HomePhone                    string `json:"home_phone"`
	MobilePhone                  string `json:"mobile_phone"`
	Facsimile                    string `json:"facsimile"`
	Email                        string `json:"email"`
	StatementType                string `json:"statement_type"`
	FATCA                        string `json:"fatca"`
	ForeignTIN                   string `json:"foreign_tin"`
	ForeignTINIssuanceCountry    string `json:"foreign_tin_issuance_country"`
	REDMPaymentBankBICCode1      string `json:"redm_payment_bank_bic_code1"`
	REDMPaymentBankBIMemberCode1 string `json:"redm_payment_bank_bi_member_code1"`
	REDMPaymentBankName1         string `json:"redm_payment_bank_name1"`
	REDMPaymentBankCountry1      string `json:"redm_payment_bank_country1"`
	REDMPaymentBankBranch1       string `json:"redm_payment_bank_branch1"`
	REDMPaymentACCcy1            string `json:"redm_payment_ac_ccy1"`
	REDMPaymentACNo1             string `json:"redm_payment_ac_no1"`
	REDMPaymentACName1           string `json:"redm_payment_ac_name1"`
	REDMPaymentBankBICCode2      string `json:"redm_payment_bank_bic_code2"`
	REDMPaymentBankBIMemberCode2 string `json:"redm_payment_bank_bi_member_code2"`
	REDMPaymentBankName2         string `json:"redm_payment_bank_name2"`
	REDMPaymentBankCountry2      string `json:"redm_payment_bank_country2"`
	REDMPaymentBankBranch2       string `json:"redm_payment_bank_branch2"`
	REDMPaymentACCcy2            string `json:"redm_payment_ac_ccy2"`
	REDMPaymentACNo2             string `json:"redm_payment_ac_no2"`
	REDMPaymentACName2           string `json:"redm_payment_ac_name2"`
	REDMPaymentBankBICCode3      string `json:"redm_payment_bank_bic_code3"`
	REDMPaymentBankBIMemberCode3 string `json:"redm_payment_bank_bi_member_code3"`
	REDMPaymentBankName3         string `json:"redm_payment_bank_name3"`
	REDMPaymentBankCountry3      string `json:"redm_payment_bank_country3"`
	REDMPaymentBankBranch3       string `json:"redm_payment_bank_branch3"`
	REDMPaymentACCcy3            string `json:"redm_payment_ac_ccy3"`
	REDMPaymentACNo3             string `json:"redm_payment_ac_no3"`
	REDMPaymentACName3           string `json:"redm_payment_ac_name3"`
	ClientCode                   string `json:"client_code"`
}

type OaRequestCsvFormatFiksTxt struct {
	DataRow string `json:"data_row"`
}
