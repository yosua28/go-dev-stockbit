package lib

import (
	"api/models"
	"fmt"
	"log"

	"github.com/tbalthazar/onesignal-go"
)

func CreateNotifCustomerFromApp(heading string, content string) {
	if Profile.TokenNotif == nil {
		log.Println("token kosong")
	} else {
		log.Println("token : " + *Profile.TokenNotif)
		playerID := &Profile.TokenNotif
		CreateNotificationHelper(**playerID, heading, content)
	}
}

func CreateNotifCustomerFromAdminByCustomerId(customerId string, heading string, content string) {
	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByCustomerKey(&userData, customerId)
	if err == nil {
		if userData.TokenNotif == nil {
			log.Println("token kosong")
		} else {
			log.Println("token : " + *userData.TokenNotif)
			playerID := &userData.TokenNotif
			CreateNotificationHelper(**playerID, heading, content)
		}
	}
}

func CreateNotifCustomerFromAdminByUserLoginKey(userLoginKey string, heading string, content string) {
	var userData models.ScUserLogin
	_, err := models.GetScUserLoginByKey(&userData, userLoginKey)
	if err == nil {
		if userData.TokenNotif == nil {
			log.Println("token kosong")
		} else {
			log.Println("token : " + *userData.TokenNotif)
			playerID := &userData.TokenNotif
			CreateNotificationHelper(**playerID, heading, content)
		}
	}
}

func CreateNotificationHelper(playerID string, heading string, content string) *onesignal.NotificationCreateResponse {
	log.Println("playerID : " + playerID)
	log.Println("Heading : " + heading)
	log.Println("Content : " + content)
	// playerID = "a00cfd56-b91a-464f-8da9-f36b376190b4" //yosua
	client := onesignal.NewClient(nil)
	client.AppKey = "YzE1YzdmN2UtNjgwNi00ZDc5LWI4ZDQtZjQyMzU3NzMzMGI5"
	notificationReq := &onesignal.NotificationRequest{
		AppID:            "8d260c99-2cb7-4159-94bf-d575d9c772dc",
		Headings:         map[string]string{"en": heading},
		Contents:         map[string]string{"en": content},
		SmallIcon:        "https://devapi.mncasset.com/images/mail/icon_md.png",
		LargeIcon:        "https://devapi.mncasset.com/images/mail/icon_mncduit.jpg",
		IncludePlayerIDs: []string{playerID},
	}
	createRes, _, err := client.Notifications.Create(notificationReq)
	if err != nil {
		log.Println("OneSignal Message")
		fmt.Println(err)
	} else {
		return createRes
	}
	return createRes
}
