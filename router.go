package main

import (
	"api/db"
	"api/controllers"
	"api/config"
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

	// Post
	e.GET("/posts/:field/:key", controllers.GetCmsPostList).Name = "GetCmsPostList"
	e.GET("/posts/:key", controllers.GetCmsPostData).Name = "GetCmsPost"

	// Fund Type
	e.GET("/fundtype", controllers.GetMsFundTypeList).Name = "GetMsFundTypeList"

	// Session
	e.POST("/register", controllers.Register).Name = "Register"
	e.POST("/verifyemail", controllers.VerifyEmail).Name = "VerifyEmail"
	e.POST("/verifyotp", controllers.VerifyOtp).Name = "VerifyOtp"
	e.POST("/login", controllers.VerifyOtp).Name = "Login"
	return e
}

func printUrlMiddleware (next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}