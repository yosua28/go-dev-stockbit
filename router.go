package main

import (
	"api/config"
	"api/controllers"
	"api/db"
	"api/lib"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func router() *echo.Echo {

	db.Db = db.DBInit()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.CORSAllowOrigin,
	}))

	e.Use(printUrlMiddleware)
	e.Use(middleware.Logger())

	auth := e.Group("/auth")
	admin := e.Group("/admin")

	auth.Use(lib.AuthenticationMiddleware)
	admin.Use(lib.AuthenticationMiddleware)

	// Post
	e.GET("/posts/:field/:key", controllers.GetCmsPostList).Name = "GetCmsPostList"
	e.GET("/posts/:key", controllers.GetCmsPostData).Name = "GetCmsPost"
	e.GET("/posttype", controllers.GetCmsPostTypeList).Name = "GetCmsPostTypeList"

	// Fund Type
	auth.GET("/fundtype", controllers.GetMsFundTypeList).Name = "GetMsFundTypeList"

	// Product
	auth.GET("/product", controllers.GetMsProductList).Name = "GetMsProductList"
	auth.GET("/product/:key", controllers.GetMsProductData).Name = "GetMsProductData"
	auth.GET("/productlist", controllers.ProductListMutasi).Name = "ProductListMutasi"

	// Message
	auth.GET("/message", controllers.GetMessageList).Name = "GetMessageList"
	auth.GET("/message/:key", controllers.GetMessageData).Name = "GetMessageData"
	auth.PATCH("/patchmessage", controllers.PatchMessage).Name = "PatchMessage"

	// Nav
	auth.GET("/nav/:duration/:product_key", controllers.GetTrNavProduct).Name = "GetTrNavProduct"

	// Lookup
	auth.GET("/lookup", controllers.GetGenLookup).Name = "GetGenLookup"

	// City
	auth.GET("/city/:field/:key", controllers.GetMsCityList).Name = "GetMsCityList"

	// Bank
	auth.GET("/bank", controllers.GetMsBankList).Name = "GetMsBankList"

	// Country
	auth.GET("/country", controllers.GetMsCountryList).Name = "GetMsCountryList"

	// Quiz
	auth.GET("/quiz", controllers.GetCmsQuiz).Name = "GetCmsQuiz"
	auth.POST("/quizresult", controllers.PostQuizAnswer).Name = "PostQuizAnswer"

	// Request
	auth.POST("/oarequest", controllers.CreateOaPersonalData).Name = "CreateOaPersonalData"

	// Invest
	auth.GET("/investpurpose", controllers.GetCmsInvestPurpose).Name = "GetCmsInvestPurpose"
	auth.GET("/investpartner", controllers.GetCmsInvestParter).Name = "GetCmsInvestParter"

	// Transaction
	auth.POST("/createtransaction", controllers.CreateTransaction).Name = "CreateTransaction"
	auth.POST("/uploadtransferpic", controllers.UploadTransferPic).Name = "UploadTransferPic"
	auth.GET("/transaction", controllers.GetTransactionList).Name = "GetTransactionList"
	auth.GET("/portofolio", controllers.Portofolio).Name = "Portofolio"
	auth.GET("/mailportofolio", controllers.SendEmailPortofolio).Name = "SendEmailPortofolio"
	auth.GET("/mailtransaction", controllers.SendEmailTransaction).Name = "SendEmailTransaction"

	// Session
	e.POST("/register", controllers.Register).Name = "Register"
	e.GET("/verifyemail", controllers.VerifyEmail).Name = "VerifyEmail"
	e.POST("/verifyotp", controllers.VerifyOtp).Name = "VerifyOtp"
	e.POST("/login", controllers.Login).Name = "Login"
	e.POST("/loginbo", controllers.LoginBo).Name = "LoginBo"
	e.POST("/resendverification", controllers.ResendVerification).Name = "ResendVerification"
	auth.GET("/user", controllers.GetUserLogin).Name = "GetUserLogin"
	auth.POST("/uploadprofilepic", controllers.UploadProfilePic).Name = "UploadProfilePic"
	auth.PUT("/changepassword", controllers.ChangePassword).Name = "ChangePassword"
	auth.GET("/servertime", controllers.CurrentTime).Name = "CurrentTime"
	e.POST("/forgotpassword", controllers.ForgotPassword).Name = "ForgotPassword"

	//Admin OA Request
	admin.GET("/oarequestlist", controllers.GetOaRequestList).Name = "GetOaRequestList"
	admin.GET("/oarequestdata/:key", controllers.GetOaRequestData).Name = "GetOaRequestData"
	admin.POST("/updatestatusapproval/cs", controllers.UpdateStatusApprovalCS).Name = "UpdateStatusApprovalCS"
	admin.POST("/updatestatusapproval/compliance", controllers.UpdateStatusApprovalCompliance).Name = "UpdateStatusApprovalCompliance"
	admin.GET("/oarequestlist/dotransaction", controllers.GetOaRequestListDoTransaction).Name = "GetOaRequestListDoTransaction"
	admin.GET("/downloadformatsinvest", controllers.DownloadOaRequestFormatSinvest).Name = "DownloadOaRequestFormatSinvest"
	admin.POST("/uploadformatsinvest", controllers.UploadOaRequestFormatSinvest).Name = "UploadOaRequestFormatSinvest"

	//Admin Post
	admin.GET("/posts", controllers.GetAdminCmsPostList).Name = "GetAdminCmsPostList"
	admin.GET("/post/:key", controllers.GetAdminCmsPostData).Name = "GetAdminCmsPostData"
	admin.POST("/post/create", controllers.CreateAdminCmsPost).Name = "CreateAdminCmsPost"
	admin.POST("/post/update/:key", controllers.UpdateAdminCmsPost).Name = "UpdateAdminCmsPost"
	admin.POST("/post/delete", controllers.DeleteAdminCmsPost).Name = "DeleteAdminCmsPost"
	admin.GET("/typelist", controllers.AdminGetListPostType).Name = "AdminGetListPostType"
	admin.GET("/subtypelist", controllers.AdminGetListPostSubtype).Name = "AdminGetListPostSubtype"
	admin.GET("/subtypelist/:post_type", controllers.AdminGetListPostSubtypeByType).Name = "AdminGetListPostSubtypeByType"

	//Admin CMS Financial Calc
	admin.GET("/financialcalc", controllers.GetAdminCmsFinancialCalcList).Name = "GetAdminCmsFinancialCalcList"

	//Admin Transaction
	admin.GET("/transactionlist", controllers.GetTransactionApprovalList).Name = "GetTransactionApprovalList"
	admin.GET("/transactionlist/cutoff", controllers.GetTransactionCutOffList).Name = "GetTransactionCutOffList"
	admin.GET("/transactionlist/batch", controllers.GetTransactionBatchList).Name = "GetTransactionBatchList"
	admin.GET("/transactionlist/confirmation", controllers.GetTransactionConfirmationList).Name = "GetTransactionConfirmationList"
	admin.GET("/transactionlist/correctionadmin", controllers.GetTransactionCorrectionAdminList).Name = "GetTransactionCorrectionAdminList"
	admin.GET("/transactionlist/posting", controllers.GetTransactionPostingList).Name = "GetTransactionPostingList"
	admin.GET("/transaction/:key", controllers.GetTransactionDetail).Name = "GetTransactionDetail"
	admin.POST("/transactionapproval/cs", controllers.TransactionApprovalCs).Name = "TransactionApprovalCs"
	admin.POST("/transactionapproval/compliance", controllers.TransactionApprovalCompliance).Name = "TransactionApprovalCompliance"
	admin.POST("/transaction/updatenavdate", controllers.UpdateNavDate).Name = "UpdateNavDate"
	admin.POST("/transactionapproval/cutoff", controllers.TransactionApprovalCutOff).Name = "TransactionApprovalCutOff"
	admin.GET("/transaction/downloadformatsinvest", controllers.DownloadTransactionFormatSinvest).Name = "DownloadTransactionFormatSinvest"
	admin.POST("/transaction/downloadformatexcel", controllers.GetFormatExcelDownloadList).Name = "GetFormatExcelDownloadList"
	admin.POST("/transaction/uploadexcelconfirmation", controllers.UploadExcelConfirmation).Name = "UploadExcelConfirmation"
	admin.POST("/transactionapproval/posting", controllers.ProsesPosting).Name = "ProsesPosting"

	//Admin Transaction type
	admin.GET("/transactiontypelist", controllers.GetTransactionType).Name = "GetTransactionType"

	//Admin Product
	admin.GET("/productlist", controllers.GetListProductAdmin).Name = "GetListProductAdmin"
	admin.GET("/product/:key", controllers.GetProductDetailAdmin).Name = "GetProductDetailAdmin"
	admin.POST("/product/delete", controllers.DeleteProductAdmin).Name = "DeleteProductAdmin"
	admin.POST("/product/create", controllers.CreateAdminMsProduct).Name = "CreateAdminMsProduct"
	admin.POST("/product/update", controllers.UpdateAdminMsProduct).Name = "UpdateAdminMsProduct"

	//Admin Product
	admin.GET("/currencylist", controllers.GetListMsCurrency).Name = "GetListMsCurrency"

	//Admin Custodian Bank
	admin.GET("/custodianbanklist", controllers.GetListMsCustodianBank).Name = "GetListMsCustodianBank"

	//Admin Fund Structure
	admin.GET("/fundstructurelist", controllers.GetListMsFundStructure).Name = "GetListMsFundStructure"

	//Admin Fund Type
	admin.GET("/fundtypelist", controllers.AdminGetListMsFundType).Name = "AdminGetListMsFundType"

	//Admin Product Category
	admin.GET("/productcategorylist", controllers.AdminGetListMsProductCategory).Name = "AdminGetListMsProductCategory"

	//Admin Product Type
	admin.GET("/producttypelist", controllers.AdminGetListMsProductType).Name = "AdminGetListMsProductType"

	//Admin Product Type
	admin.GET("/riskprofilelist", controllers.AdminGetListMsRiskProfile).Name = "AdminGetListMsRiskProfile"

	//Admin Product
	admin.GET("/productfeelist", controllers.GetListProductFeeAdmin).Name = "GetListProductFeeAdmin"
	admin.GET("/productfee/:key", controllers.GetProductFeeDetailAdmin).Name = "GetProductFeeDetailAdmin"
	admin.GET("/productlist/dropdown", controllers.GetListProductAdminDropdown).Name = "GetListProductAdminDropdown"
	admin.POST("/productfee/delete", controllers.DeleteProductFeeAdmin).Name = "DeleteProductFeeAdmin"
	admin.POST("/productfee/create", controllers.CreateAdminMsProductFee).Name = "CreateAdminMsProductFee"
	admin.POST("/productfee/update", controllers.UpdateAdminMsProductFee).Name = "UpdateAdminMsProductFee"
	admin.POST("/productfee/feeitem/create", controllers.CreateAdminMsProductFeeItem).Name = "CreateAdminMsProductFeeItem"
	admin.POST("/productfee/feeitem/update", controllers.UpdateAdminMsProductFeeItem).Name = "UpdateAdminMsProductFeeItem"
	admin.POST("/productfee/feeitem/delete", controllers.DeleteAdminMsProductFeeItem).Name = "DeleteAdminMsProductFeeItem"
	admin.GET("/productfee/feeitem/:key", controllers.DetailAdminMsProductFeeItem).Name = "DetailAdminMsProductFeeItem"

	//Admin Product Bank Account
	admin.GET("/productbankaccountlist", controllers.GetListProductBankAccountAdmin).Name = "GetListProductBankAccountAdmin"
	admin.GET("/productbankaccount/:key", controllers.GetProductBankAccountDetailAdmin).Name = "GetProductBankAccountDetailAdmin"
	admin.POST("/productbankaccount/delete", controllers.DeleteProductBankAccountAdmin).Name = "DeleteProductBankAccountAdmin"
	admin.POST("/productbankaccount/create", controllers.CreateAdminMsProductBankAccount).Name = "CreateAdminMsProductBankAccount"
	admin.POST("/productbankaccount/update", controllers.UpdateAdminMsProductBankAccount).Name = "UpdateAdminMsProductBankAccount"

	//Admin User Management
	admin.GET("/logout", controllers.LogoutAdmin).Name = "LogoutAdmin"
	admin.GET("/usermanagementlist", controllers.GetListScUserLoginAdmin).Name = "GetListScUserLoginAdmin"
	admin.GET("/usermanagement/:key", controllers.GetDetailScUserLoginAdmin).Name = "GetDetailScUserLoginAdmin"
	admin.POST("/usermanagement/disableenable", controllers.DisableEnableUser).Name = "DisableEnableUser"
	admin.POST("/usermanagement/lockunlock", controllers.LockUnlockUser).Name = "LockUnlockUser"
	admin.GET("/rolelist", controllers.AdminGetListScRole).Name = "AdminGetListScRole"
	admin.GET("/usercategorylist", controllers.AdminGetListScUserCategory).Name = "AdminGetListScUserCategory"
	admin.GET("/userdeptlist", controllers.AdminGetListScUserDept).Name = "AdminGetListScUserDept"
	admin.POST("/usermanagement/create", controllers.CreateAdminScUserLogin).Name = "CreateAdminScUserLogin"
	admin.POST("/usermanagement/update", controllers.UpdateAdminScUserLogin).Name = "UpdateAdminScUserLogin"
	admin.POST("/usermanagement/changepassword", controllers.ChangePasswordUser).Name = "ChangePasswordUser"
	admin.POST("/usermanagement/changerole", controllers.ChangeRoleUser).Name = "ChangeRoleUser"
	admin.POST("/usermanagement/delete", controllers.DeleteUser).Name = "DeleteUser"

	//Admin Role Management
	admin.GET("/rolemanagementlist", controllers.GetListRoleManagementAdmin).Name = "GetListRoleManagementAdmin"

	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}
