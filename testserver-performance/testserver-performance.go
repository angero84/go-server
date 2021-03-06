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
	"github.com/angero84/go-server/kcontainer"
)

type serverConfig struct {
	Port 			uint16 				`json:"Port"`
	AcceptorOpt		ktcp.KAcceptorOpt	`json:"AcceptorOpt"`
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	kdlogger, err := klog.NewKDefaultLogger(&klog.KDefaultLoggerOpt{
		LogTypeDepth:		klog.KLogType_Fatal,
		LoggerName:			"perform",
		RootDirectoryName:	"log",
		UseQueue:			false,
		StoringPeriodDay:	30,
	})

	if nil != err {
		println("Failed NewKDefaultLogger : ", err.Error())
		return
	}

	klog.SetDefaultLoggerInstance(kdlogger)

	klog.LogInfo("Performance test server started")

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

	container, err := kcontainer.NewKMapConn(2000)
	if nil != err {
		klog.LogWarn("Failed to create container : %s", err.Error())
		return
	}

	connhOpt := &ktcp.KConnHandleOpt{
		Handler:	khandler.NewKConnEventEchoServer(),
		Protocol:	&kprotocol.KProtocol{},
	}

	connhOpt.Handler.(*khandler.KConnEventEchoServer).SetContainer(container)

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