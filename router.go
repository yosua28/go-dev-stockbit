package main

import (
	"api/db"
	"api/controllers"

	"github.com/labstack/echo"
)

func router() *echo.Echo {
	
	db.Db = db.DBInit()
	e := echo.New()
	e.GET("/posts", controllers.GetCmsPostTypeData).Name = "test"
	return e
}