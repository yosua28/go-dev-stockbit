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

	// Session
	e.POST("/register", controllers.Register).Name = "Register"
	e.GET("/verifyemail", controllers.VerifyEmail).Name = "VerifyEmail"
	e.POST("/verifyotp", controllers.VerifyOtp).Name = "VerifyOtp"
	e.POST("/login", controllers.Login).Name = "Login"
	e.POST("/resendverification", controllers.ResendVerification).Name = "ResendVerification"
	auth.GET("/user", controllers.GetUserLogin).Name = "GetUserLogin"
	auth.POST("/uploadprofilepic", controllers.UploadProfilePic).Name = "UploadProfilePic"
	auth.PUT("/changepassword", controllers.ChangePassword).Name = "ChangePassword"
	auth.GET("/servertime", controllers.CurrentTime).Name = "CurrentTime"

	//Admin OA Request
	admin.GET("/oarequestlist", controllers.GetOaRequestList).Name = "GetOaRequestList"
	admin.GET("/oarequestdata/:key", controllers.GetOaRequestData).Name = "GetOaRequestData"
	admin.POST("/updatestatusapproval/cs", controllers.UpdateStatusApprovalCS).Name = "UpdateStatusApprovalCS"
	admin.POST("/updatestatusapproval/compliance", controllers.UpdateStatusApprovalCompliance).Name = "UpdateStatusApprovalCompliance"
	admin.GET("/oarequestlist/dotransaction", controllers.GetOaRequestListDoTransaction).Name = "GetOaRequestListDoTransaction"
	admin.GET("/downloadformatsinvest", controllers.DownloadOaRequestFormatSinvest).Name = "DownloadOaRequestFormatSinvest"

	//Admin Post
	admin.GET("/posts", controllers.GetAdminCmsPostList).Name = "GetAdminCmsPostList"
	admin.GET("/post/:key", controllers.GetAdminCmsPostData).Name = "GetAdminCmsPostData"
	admin.POST("/post/create", controllers.CreateAdminCmsPost).Name = "CreateAdminCmsPost"
	admin.POST("/post/update/:key", controllers.UpdateAdminCmsPost).Name = "UpdateAdminCmsPost"
	admin.POST("/post/delete", controllers.DeleteAdminCmsPost).Name = "DeleteAdminCmsPost"

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

	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}
