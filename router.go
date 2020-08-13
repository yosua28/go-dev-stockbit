package main

import (
	"api/db"
	"api/controllers"

	"github.com/labstack/echo"
)

func router() *echo.Echo {
	
	db.Db = db.DBInit()
	e := echo.New()

	// Post
	e.GET("/posts/:field/:key", controllers.GetCmsPostList).Name = "GetCmsPostList"
	e.GET("/posts/:key", controllers.GetCmsPostList).Name = "GetCmsPost"

	// Fund Type
	e.GET("/fundtype", controllers.GetMsFundTypeList).Name = "GetMsFundTypeList"

	// User
	e.POST("/register", controllers.Register).Name = "Register"
	return e
}