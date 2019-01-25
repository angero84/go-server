package main

import (
	"kdatabase"
	klog "klogger"
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

	result := db.QueryResults("call usp_user_dispatch_info_select(?)", 1517 )
	if nil == result {
		return
	}

	chSig := make(chan os.Signal)


	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	db.Destroy()
	klog.LogInfo("Main end")
}