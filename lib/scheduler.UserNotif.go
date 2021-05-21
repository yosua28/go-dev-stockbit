package lib

import (
	"api/models"
	"fmt"
	"log"
	"strconv"
	"time"
)

func UserNotifOnceAday() {
	blastUserNotif("312")
}

func UserNotifOnceBeginOnceBefore() {
	blastUserNotif("311")
}

func UserNotifOnce() {
	blastUserNotif("310")
}

func blastUserNotif(notifType string) {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON USER NOTIF - TYPE " + notifType + " : " + time.Now().Format(dateLayout))
	var err error

	//cek promo
	var uNotif []models.UserNotifCron
	if notifType == "310" {
		_, err = models.AdminGetAllUserNotifOnceCron(&uNotif)
	} else if notifType == "311" {
		_, err = models.AdminGetAllUserNotifOnceBeginOnceBefore(&uNotif)
	} else if notifType == "312" {
		_, err = models.AdminGetAllUserNotifOnceAday(&uNotif)
	}

	if err == nil {
		if len(uNotif) > 0 {
			//cek user
			var customer []models.UserBlastPromo
			_, err = models.AdminGetAllUserBlastPromo(&customer)
			if err == nil {
				if len(customer) > 0 {
					//insert message
					var bindVarMessage []interface{}
					for _, un := range uNotif { //user notif
						for _, cus := range customer { //customer
							var row []string
							row = append(row, "247")                                    //umessage_type
							row = append(row, strconv.FormatUint(un.NotifHdrKey, 10))   //notif_hdr_key
							row = append(row, strconv.FormatUint(cus.UserLoginKey, 10)) //umessage_recipient_key
							row = append(row, time.Now().Format(dateLayout))            //umessage_receipt_date
							row = append(row, "0")                                      //flag_read
							row = append(row, "1")                                      //flag_sent
							row = append(row, *un.UmessageSubject)                      //umessage_subject
							row = append(row, *un.UmessageBody)                         //umessage_body
							row = append(row, "249")                                    //umessage_category
							row = append(row, "0")                                      //flag_archieved
							row = append(row, time.Now().Format(dateLayout))            //archieved_date
							row = append(row, "1")                                      //rec_status
							row = append(row, time.Now().Format(dateLayout))            //rec_created_date
							row = append(row, "CRON")                                   //rec_created_by
							bindVarMessage = append(bindVarMessage, row)
						}
					}

					var playerIds []string

					for _, cus := range customer { //customer
						playerIds = append(playerIds, *cus.TokenNotif)
					}

					_, err = models.CreateMultipleUserMessageFromUserNotif(bindVarMessage)
					if err != nil {
						fmt.Println("err create multiple user message : " + err.Error())
					} else {
						fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
					}

					//push notif

					var heading string
					var content string
					if len(uNotif) > 0 {
						for _, un := range uNotif { //user notif
							heading = *un.UmessageSubject
							content = *un.UmessageBody

							DataNotif := make(map[string]interface{})
							DataNotif["category"] = "MESSAGE"
							log.Println("=====================================")
							log.Println(playerIds)
							log.Println(heading)
							log.Println(content)
							log.Println(DataNotif)
							log.Println("=====================================")
							BlastAllNotificationHelper(playerIds, heading, content, DataNotif)
						}
					}
				}
			} else {
				fmt.Println("err get customer : " + err.Error())
			}

		}
	} else {
		fmt.Println("err get user notif : " + err.Error())
	}

	fmt.Println("======END CRON USER NOTIF BLAST============")
}
