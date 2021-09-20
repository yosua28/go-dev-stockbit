package main

import (
	"api/controllers"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func Goroutines() {

	searchMoviesTest()
	detailMoviesTest()

	str := "test(oke stoctbit)"
	res := controllers.FindFirstStringInBracket(str)
	fmt.Println(res)

	arrAnagram := []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}
	data := controllers.LogisTest(arrAnagram)
	fmt.Println(data)
}

func searchMoviesTest() {
	gin.SetMode(gin.TestMode)

	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.GET("/search-movies/Batmen/2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success get search",
		})
	})

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	req, err := http.NewRequest(http.MethodGet, "/search-movies/Batmen/2", nil)
	if err != nil {
		fmt.Println("Couldn't create request: %v\n", err)
	}

	log.Println(req)

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		fmt.Println("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		fmt.Println("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}

func detailMoviesTest() {
	// Switch to test mode so you don't get such noisy output
	gin.SetMode(gin.TestMode)

	// Setup your router, just like you did in your main function, and
	// register your routes
	r := gin.Default()
	r.GET("/detail-movies/tt4853102", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "success detail search",
		})
	})

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	req, err := http.NewRequest(http.MethodGet, "/detail-movies/tt4853102", nil)
	if err != nil {
		fmt.Println("Couldn't create request: %v\n", err)
	}

	log.Println(req)

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)
	fmt.Println(w.Body)

	// Check to see if the response was what you expected
	if w.Code == http.StatusOK {
		fmt.Println("Expected to get status %d is same ast %d\n", http.StatusOK, w.Code)
	} else {
		fmt.Println("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}
}
