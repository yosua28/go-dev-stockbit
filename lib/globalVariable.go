package lib

//mail Template Name - mm_mail_master
var EMAIL_ACTIVATION string = "EMAIL-ACTIVATION"

//MOTION-PAY LINKING API NAME
var MAX_AMOUNT_MOTION_PAY int64 = 2000000
var REGISTRATION_PREMIUM_NO_OTP string = "REGISTRATION_PREMIUM_NO_OTP"
var PATH_REGISTRATION_PREMIUM_NO_OTP string = "v1/merchants/users/registration/noauth/premium"

var CASH_BALANCE string = "CASH_BALANCE"
var PATH_CASH_BALANCE string = "v1/merchants/users/balance/cash"

var POINT_BALANCE string = "POINT_BALANCE"
var PATH_POINT_BALANCE string = "v1/merchants/users/balance/points"

var STATUS_NON_LINKED string = "NON LINKED"
var STATUS_LINKED string = "LINKED"
var STATUS_UNLINKED string = "UNLINKED"

var PREMIUM string = "PREMIUM"

//MOTION PAY - PAYMENT
var CREATE_ORDER string = "CREATE_ORDER"
var PATH_CREATE_ORDER string = "v1/merchants/orders"

var CREATE_OTP string = "CREATE_OTP"
var PATH_CREATE_OTP string = "v1/merchants/pay/otp"

var PAY_ORDER string = "PAY_ORDER"
var PATH_PAY_ORDER string = "v1/merchants/pay"

var ORDER_DETAIL string = "ORDER_DETAIL"
var PATH_ORDER_DETAIL string = "v1/merchants/orders"

//ROLE GROUP
var ROLE_CS = "11"
var ROLE_KYC = "12"
var ROLE_FUND_ADMIN = "13"

//PAYMENT METHOD LOOKUP
var PAYMENT_MOTION_PAY = "285"
var PAYMENT_VIRTUAL_ACCOUNT = "287"
var PAYMENT_TRANSFER_MANUAL = "284"

var UNSETTLED = "243"
var SETTLED = "244"
