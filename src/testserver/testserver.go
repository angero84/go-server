package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"encoding/json"

	"ktcp"
	"kprotocol"
	klog "klogger"
	"khandler"

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
		klog.LogWarn("Failed acceptor server : %s", err.Error())
		return
	}

	go acceptor.Start()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	acceptor.StopGoRoutineWait()
	klog.LogInfo("Main end")
}