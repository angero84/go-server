package main

import (

	klog "klogger"

	"kprotocol"
	"ktcp"
	"khandler"
	"fmt"
	"sync"
)



func main() {

	cliOpt := &ktcp.KClientOpt{
		TargetRemoteIP:	"0.0.0.0",
		TargetPort:		8989,
		Reconnect:		true,
		ReconnectIntervalTime: 5000,
	}

	connhOpt := &ktcp.KConnHandleOpt{
		Handler:	khandler.NewKConnHandlerJson(khandler.NewProcessorExampleJson()),
		Protocol:	&kprotocol.KProtocolJson{},
	}

	client, err := ktcp.NewKClient(0, cliOpt, nil, connhOpt )
	if nil != err {
		klog.LogWarn("Create client failed : %s", err.Error())
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.Connect()
		if nil != err {
			klog.LogWarn("Connect err : %s", err.Error())
			return
		}
	}()

	wg.Wait()

	if client.Connected() {
		cmd := ""

		LOOP:
		for {

			fmt.Scanln(&cmd)

			switch cmd {
			case "exit":
				break LOOP
			case "disconnect":
				client.Disconnect()
			case "connect":
				client.ConnectAsync(nil)
			default:

				chat := kprotocol.ProtocolJsonRequestChatting{}
				chat.Chat = cmd
				chat.ChatType = "[normal]"

				err := client.Send(chat.MakePacket())
				if nil != err {
					klog.LogWarn("Send err : %s", err.Error())
				}
			}

		}
	}

}

