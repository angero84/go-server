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

	log "logger"
)

type serverConfig struct {
	Port 			uint32 			`json:"Port"`
	TcpConfig		tcp.Config 		`json:"TcpConfig"`
}

func main() {

	log.LogInfo("Testserver started")

	runtime.GOMAXPROCS(runtime.NumCPU())

	serverConfigBytes, err := ioutil.ReadFile("configServer.json")
	if nil != err {
		log.LogFatal("Cannot read config file : %s", err.Error())
		return
	}

	serverConfig := &serverConfig{}
	err = json.Unmarshal(serverConfigBytes, serverConfig)
	if nil != err {
		log.LogFatal("Failed unmarshal config file : %s", err.Error())
		return
	}

	srv, err := tcp.NewServer(serverConfig.Port, &serverConfig.TcpConfig, &tcp.CallbackEcho{}, &protocol.EchoProtocol{})
	if nil != err {
		log.LogFatal("Failed create server : %s", err.Error())
		return
	}

	go srv.Start()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.StopGoRoutineWait()
	log.LogInfo("Main end")
}