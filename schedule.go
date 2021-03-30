package main

import (
	"api/lib"
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func task() {
	fmt.Println("I am running task.")
}

func scheduler() {
	// gocron.Every(1).Second().Do(task)
	//unlock user
	gocron.Every(30).Minutes().Do(unlockuser)

	//promo blast
	gocron.Every(1).Day().At("09:02").Do(promoOnceAday)            //lookup:312
	gocron.Every(1).Day().At("10:02").Do(promoOnceBeginOnceBefore) //lookup:311
	gocron.Every(1).Day().At("11:02").Do(promoOnce)                //lookup:310
	<-gocron.Start()
}

func unlockuser() {
	lib.Unlockuser()
}

func promoOnceAday() {
	lib.PromoOnceAday()
}

func promoOnceBeginOnceBefore() {
	lib.PromoOnceBeginOnceBefore()
}

func promoOnce() {
	lib.PromoOnce()
}
