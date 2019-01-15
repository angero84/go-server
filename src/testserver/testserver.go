package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"khandler"
	klog "klogger"
	"kprotocol"
	"ktcp"
)

type serverConfig struct {
	Port 			uint32 				`json:"Port"`
	AcceptorOpt		ktcp.KAcceptorOpt	`json:"AcceptorOpt"`
}

func main() {

	klog.LogInfo("Testserver started")

	runtime.GOMAXPROCS(runtime.NumCPU())

	serverConfigBytes, err := ioutil.ReadFile("configServer.json")
	if nil != err {
		klog.LogWarn("Cannot read config file : %s", err.Error())
		return
	}

	serverConfig := &serverConfig{}
	err = json.Unmarshal(serverConfigBytes, serverConfig)
	if nil != err {
		klog.LogWarn("Failed unmarshal config file : %s", err.Error())
		return
	}

	handler := khandler.NewKConnHandlerEcho()
	protocol := &kprotocol.KProtocolEcho{}

	acceptor, err := ktcp.NewAcceptor(serverConfig.Port, &serverConfig.AcceptorOpt, handler, protocol)
	if nil != err {
		klog.LogWarn("Failed to create acceptor : %s", err.Error())
		return
	}

	chSig := make(chan os.Signal)

	go func () {
		err = acceptor.Start()
		if nil != err {
			klog.LogFatal("Failed start acceptor : %s", err.Error())
			chSig <- syscall.SIGTERM
		}
	}()


	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	acceptor.StopGoRoutineWait()
	klog.LogInfo("Main end")
}