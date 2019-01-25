package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/angero84/go-server/khandler"
	klog "github.com/angero84/go-server/klogger"
	"github.com/angero84/go-server/kprotocol"
	"github.com/angero84/go-server/ktcp"
)

type serverConfig struct {
	Port 			uint16 				`json:"Port"`
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

	connhOpt := &ktcp.KConnHandleOpt{
		Handler:	khandler.NewKConnEventExample(khandler.NewMessageHandlerExample()),
		Protocol:	&kprotocol.KProtocol{},
	}

	acceptor, err := ktcp.NewKAcceptor(serverConfig.Port, &serverConfig.AcceptorOpt, connhOpt )
	if nil != err {
		klog.LogWarn("Failed to create acceptor : %s", err.Error())
		return
	}

	chSig := make(chan os.Signal)

	go func () {
		err = acceptor.Listen()
		if nil != err {
			klog.LogFatal("Failed start acceptor : %s", err.Error())
			chSig <- syscall.SIGTERM
		}
	}()


	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	acceptor.Destroy()
	klog.LogInfo("Main end")
}