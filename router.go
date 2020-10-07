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

	// Session
	e.POST("/register", controllers.Register).Name = "Register"
	e.GET("/verifyemail", controllers.VerifyEmail).Name = "VerifyEmail"
	e.POST("/verifyotp", controllers.VerifyOtp).Name = "VerifyOtp"
	e.POST("/login", controllers.Login).Name = "Login"
	e.POST("/resendverification", controllers.ResendVerification).Name = "ResendVerification"

	//Admin OA Request
	admin.GET("/oarequestlist", controllers.GetOaRequestList).Name = "GetOaRequestList"
	admin.GET("/oarequestdata/:key", controllers.GetOaRequestData).Name = "GetOaRequestData"
	admin.POST("/updatestatusapproval/cs", controllers.UpdateStatusApprovalCS).Name = "UpdateStatusApprovalCS"

	//Admin Post
	admin.GET("/posts", controllers.GetAdminCmsPostList).Name = "GetAdminCmsPostList"

	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}
