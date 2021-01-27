package controllers

import (
	"fmt"
	"log"

	"github.com/maddevsio/fcm"
)

func CreateNotification(token string, title string, body string) {
	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	serverKey := "AAAACgKWnuY:APA91bEdhBJmGR5mYx9Lyd8jimPSh8bAz65ao6cCmmA3-O1vJBIML7a6-IyQ0b9giER2-EYpWBWriJdODPdSMmaAZsWCcxdgsnx_lpSACx5HCKMug8wDs0XrrDbzsbiaVo6rl3_ui84q"

	c := fcm.NewFCM(serverKey)
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{"/topics/dev-all"},
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: title,
			Body:  body,
			Sound: "default",
			Badge: "3",
		},
	})
	if err != nil {
		log.Println("------------------")
		log.Println(err)
		log.Println("------------------")
	}
	log.Println("------------------")
	log.Println(err)
	log.Println(response)
	log.Println("------------------")
	fmt.Println("Status Code   :", response.StatusCode)
	fmt.Println("Success       :", response.Success)
	fmt.Println("Fail          :", response.Fail)
	fmt.Println("Canonical_ids :", response.CanonicalIDs)
	fmt.Println("Topic MsgId   :", response.MsgID)
	fmt.Println("Topic Error   :", response.RetryAfter)
}
