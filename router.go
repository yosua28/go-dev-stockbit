package main

import (
	"api/controllers"
	"api/db"
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func router() *echo.Echo {

	db.Db = db.DBInit()
	e := echo.New()

	e.Use(printUrlMiddleware)
	e.Use(middleware.Logger())

	// Movie
	e.GET("/search-movies/:searchword/:pagination", controllers.SearchMovies).Name = "SearchMovies"
	e.GET("/detail-movies/:imdbid", controllers.DetailMovies).Name = "DetailMovies"
	return e
}

func printUrlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println(c.Request().URL)
		return next(c)
	}
}
