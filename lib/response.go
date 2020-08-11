package lib

import (
	"net/http"

	"github.com/labstack/echo"
)

type Status struct {
	Code          uint   `json:"code"`
	MessageServer string `json:"message_server"`
	MessageClient string `json:"message_client"`
}

type Response struct {
	Status Status      `json:"status"`
	Data   interface{} `json:"data"`
}

func CustomError(code int, messages ...string) *echo.HTTPError {
	var response Response
	response.Status.Code = uint(code)
	response.Status.MessageServer = http.StatusText(code)
	response.Status.MessageClient = http.StatusText(code)
	for index, value := range messages {
		if index == 0 {
			response.Status.MessageServer = value
		}
		if index == 1 {
			response.Status.MessageClient = value
		}
	}


	return echo.NewHTTPError(code, response)
}


