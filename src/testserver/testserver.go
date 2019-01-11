package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"protocol"
	"runtime"
	"syscall"
	"tcp"

	"encoding/json"

	klog "logger"
)

type serverConfig struct {
	Port 			uint32 			`json:"Port"`
	TcpConfig		tcp.Config 		`json:"TcpConfig"`
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

	srv, err := tcp.NewServer(serverConfig.Port, &serverConfig.TcpConfig, &tcp.CallbackEcho{}, &protocol.EchoProtocol{})
	if nil != err {
		klog.LogWarn("Failed create server : %s", err.Error())
		return
	}

	go srv.Start()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.StopGoRoutineWait()
	klog.LogInfo("Main end")
}