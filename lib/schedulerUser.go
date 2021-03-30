package lib

import (
	"api/models"
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
)

func Unlockuser() {
	dateLayout := "2006-01-02 15:04:05"
	fmt.Println("CRON UNLOCK USER : " + time.Now().Format(dateLayout))
	var users []models.UserLoginKeyLocked

	_, err := models.GetUserLocked(&users)
	if err == nil {
		if len(users) < 1 {
			fmt.Println("NO DATA CRON UNLOCK USER")
		} else {
			var userIds []string
			for _, us := range users {
				if _, ok := Find(userIds, strconv.FormatUint(us.UserLoginKey, 10)); !ok {
					userIds = append(userIds, strconv.FormatUint(us.UserLoginKey, 10))
				}
			}
			if len(userIds) > 0 {
				params := make(map[string]string)
				params["ulogin_locked"] = "0"
				params["ulogin_failed_count"] = "0"
				params["rec_modified_by"] = "CRON UNLOCK USER"
				params["rec_modified_date"] = time.Now().Format(dateLayout)
				_, err = models.UpdateScUserLoginByKeyIn(params, userIds, "user_login_key")
				if err != nil {
					log.Error("ERROR CRON UNLOCK USER : " + err.Error())
				} else {
					fmt.Println("UNLOCK USER DONE. Jumlah Data : ")
					fmt.Println(len(userIds))
				}
			}
		}
	} else {
		fmt.Println("NO DATA CRON UNLOCK USER")
	}
	fmt.Println("======END CRON UNLOCK USER============")
}
