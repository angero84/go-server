package main

import (
	"runtime"

	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tcp"
	"protocol"
	"io/ioutil"

	"encoding/json"

	log "logger"

	"time"
)

type serverConfig struct {
	Port 			uint32 			`json:"Port"`
	TcpConfig		tcp.Config 		`json:"TcpConfig"`
}

func main() {


	tempChan := make(chan int)
	tempChan = nil

	go func() {

		select {
		case <-time.After(time.Second*5):
			tempChan <- 1
		}

	}()

	log.LogInfo("testserver started")

	runtime.GOMAXPROCS(runtime.NumCPU())

	serverConfigBytes, err := ioutil.ReadFile("configServer.json")
	if nil != err {
		fmt.Println("config read : ", err)
		return
	}

	serverConfig := &serverConfig{}
	err = json.Unmarshal(serverConfigBytes, serverConfig)
	if nil != err {
		return
	}

	srv, err := tcp.NewServer(serverConfig.Port, &serverConfig.TcpConfig, &tcp.CallbackEcho{}, &protocol.EchoProtocol{})
	if nil != err {
		return
	}

	go srv.Start()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}