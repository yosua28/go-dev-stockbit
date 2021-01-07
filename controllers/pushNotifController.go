package controllers

import (
	"log"

	"github.com/douglasmakey/go-fcm"
)

func CreateNotification(token string) {
	// init client
	// client := fcm.NewClient("AAAASUHZQg4:APA91bG165ItNXI_35P27hHURe3PHXhFqxtOMo4XpsIdSgsK2_No--tb9CXGIplhP1DZfLlsyLAB2LKct16wbpyInBKIlrBnepzwoVaAjyoGMYDe4KArzQR20R3T7hwIkzIyKwXICLtp")
	client := fcm.NewClient("AAAACgKWnuY:APA91bEdhBJmGR5mYx9Lyd8jimPSh8bAz65ao6cCmmA3-O1vJBIML7a6-IyQ0b9giER2-EYpWBWriJdODPdSMmaAZsWCcxdgsnx_lpSACx5HCKMug8wDs0XrrDbzsbiaVo6rl3_ui84q")

	// You can use your HTTPClient
	//client.SetHTTPClient(client)

	data := map[string]interface{}{
		"message": "From Go-FCM",
		"details": map[string]string{
			"name":  "Name",
			"user":  "Admin",
			"thing": "none",
		},
	}

	// You can use PushMultiple or PushSingle
	// client.PushMultiple([]string{"token 1", "token 2"}, data)
	client.PushSingle(token, data)

	// registrationIds remove and return map of invalid tokens
	badRegistrations := client.CleanRegistrationIds()
	log.Println(badRegistrations)

	status, err := client.Send()
	if err != nil {
		log.Println("---------------------------")
		log.Println(err)
		log.Println("---------------------------")
	}

	log.Println(status.Results)
}
