package lib

import (
	"api/models"
	"fmt"
	"log"
	"strconv"
	"time"
)

func PromoOnceAday() {
	blast("312")
}

func PromoOnceBeginOnceBefore() {
	blast("311")
}

func PromoOnce() {
	blast("310")
}

func blast(notifType string) {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("START CRON PROMO - TYPE " + notifType + " : " + time.Now().Format(dateLayout))
	var err error

	//cek promo
	var promos []models.TrPromoCron
	if notifType == "310" {
		_, err = models.AdminGetAllPromoOnce(&promos)
	} else if notifType == "311" {
		_, err = models.AdminGetAllPromoOnceBeginOnceBefore(&promos)
	} else if notifType == "312" {
		_, err = models.AdminGetAllPromoOnceAday(&promos)
	}

	if err == nil {
		if len(promos) > 0 {
			//cek user
			var customer []models.UserBlastPromo
			_, err = models.AdminGetAllUserBlastPromo(&customer)
			if err == nil {
				if len(customer) > 0 {
					//insert message
					var bindVarMessage []interface{}
					for _, pro := range promos { //promo
						for _, cus := range customer { //customer
							var row []string
							row = append(row, "247")                                    //umessage_type
							row = append(row, strconv.FormatUint(cus.UserLoginKey, 10)) //umessage_recipient_key
							row = append(row, time.Now().Format(dateLayout))            //umessage_receipt_date
							row = append(row, "0")                                      //flag_read
							row = append(row, "1")                                      //flag_sent
							row = append(row, *pro.PromoTitle)                          //umessage_subject
							row = append(row, pro.PromoDescription)                     //umessage_body
							row = append(row, "248")                                    //umessage_category
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

					_, err = models.CreateMultipleUserMessage(bindVarMessage)
					if err != nil {
						fmt.Println("err create multiple user message : " + err.Error())
					} else {
						fmt.Println("SUKSES CREATE MULTIPLE USER MESSAGE")
					}

					//push notif
					var heading string
					var content string
					if len(promos) > 1 {
						heading = "Banyak PROMO baru lho!!"
						content = "Kamu bisa Subscribe atau Topup dengan promo-promo menarik saat ini juga. Segera Subscribe atau Topup sebelum voucher promo Kamu tidak berlaku lagi. Cek voucher di sini."
					} else {
						heading = *promos[0].PromoTitle
						content = promos[0].PromoDescription
					}

					DataNotif := make(map[string]interface{})
					DataNotif["category"] = "PROMO"
					log.Println("=====================================")
					log.Println(playerIds)
					log.Println(heading)
					log.Println(content)
					log.Println(DataNotif)
					log.Println("=====================================")
					BlastAllNotificationHelper(playerIds, heading, content, DataNotif)
				}
			} else {
				fmt.Println("err get customer : " + err.Error())
			}

		}
	} else {
		fmt.Println("err get promo : " + err.Error())
	}

	fmt.Println("======END CRON PROMO BLAST============")
}
