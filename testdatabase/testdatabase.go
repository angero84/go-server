package main

import (
	"github.com/angero84/go-server/kdatabase"
	klog "github.com/angero84/go-server/klogger"
	"os"
	"syscall"
	"os/signal"
	"fmt"
)



func main() {

	dbinfo := &kdatabase.KDBInfo{
		Driver:		"mysql",
		Account:	"trinitygames",
		Password:	"vmflwld",
		Host:		"freezing.tr",
		Port:		3306,
		Database:	"freezing",
	}

	db, err := kdatabase.NewKDB(dbinfo, nil)
	if nil != err {
		klog.LogWarn("Failed connect db : %s", err.Error())
		return
	}

	result := db.Query("call usp_user_dispatch_info_select(?)", 1564 )
	if nil == result {
		return
	}
	defer result.Close()

	for result.Next() {

		mileage := int(0)
		err := result.Scan(&mileage)
		if nil == err {
			println(mileage)
		}
	}

	if result.NextResultSet() {

		for result.Next() {
			country, state, remaintime := 0,0,0
			err := result.Scan(&country, &state, &remaintime)
			if nil == err {
				println(country, " ", state, " ", remaintime)
			}
		}

	}



	chSig := make(chan os.Signal)


	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	db.Destroy()
	klog.LogInfo("Main end")
}