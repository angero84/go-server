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

)


func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())



	serverConfigBytes, err := ioutil.ReadFile("configServer.json")
	if nil != err {
		fmt.Println("config read : ", err)
		return
	}

	srv, err := tcp.NewServer(serverConfigBytes, &tcp.CallbackEcho{}, &protocol.EchoProtocol{})
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