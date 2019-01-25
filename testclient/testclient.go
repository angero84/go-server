package main

import (

	klog "github.com/angero84/go-server/klogger"

	"github.com/angero84/go-server/kprotocol"
	"github.com/angero84/go-server/ktcp"
	"github.com/angero84/go-server/khandler"
	"fmt"
	"strings"
)

func main() {

	cliOpt := &ktcp.KClientOpt{
		TargetRemoteIP:	"0.0.0.0",
		TargetPort:		8989,
		Reconnect:		true,
		ReconnectIntervalTime: 5000,
	}

	connhOpt := &ktcp.KConnHandleOpt{
		Handler:	khandler.NewKConnEventExample(khandler.NewMessageHandlerExample()),
		Protocol:	&kprotocol.KProtocol{},
	}

	client, err := ktcp.NewKClient(0, cliOpt, nil, connhOpt )
	if nil != err {
		klog.LogWarn("Create client failed : %s", err.Error())
		return
	}

	cmd := ""

	LOOP:
	for {

		fmt.Scanln(&cmd)
		splits := strings.Split(cmd,":")
		if 0 < len(splits) {
			switch splits[0] {
			case "exit":
				break LOOP
			case "disconnect":
				client.Disconnect()
			case "connect":
				client.ConnectAsync(nil)
			case "login":
				req := &kprotocol.ProtocolLoginRequest{}
				req.UserID = "angero"
				req.Password = "1234qwer"
				err := client.Send(req)
				if nil != err {
					klog.LogWarn("Send err : %s", err.Error())
				}
			case "chat":
				if 1 < len(splits) {
					req := &kprotocol.ProtocolChattingRequest{}
					req.Chat = splits[1]
					req.ChatType = "[normal]"
					err := client.Send(req)
					if nil != err {
						klog.LogWarn("Send err : %s", err.Error())
					}
				}
			default:

			}
		}
	}


}

