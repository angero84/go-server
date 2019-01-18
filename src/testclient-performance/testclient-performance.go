package main

import (

	klog "klogger"

	"kprotocol"
	"ktcp"
	"khandler"
	"fmt"

	"kcontainer"
	"robot"
	"runtime"
)



func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	kdlogger, err := klog.NewKDefaultLogger(&klog.KDefaultLoggerOpt{
		LogTypeDepth:		klog.KLogType_Fatal,
		LoggerName:			"perform",
		RootDirectoryName:	"log",
		UseQueue:			false,
	})

	if nil != err {
		println("Failed NewKDefaultLogger : ", err.Error())
		return
	}

	klog.SetDefaultLoggerInstance(kdlogger)

	cliOpt := &ktcp.KClientOpt{
		TargetRemoteIP:	"0.0.0.0",
		TargetPort:		8989,
		Reconnect:		true,
		ReconnectIntervalTime: 5000,
	}

	connOpt := &ktcp.KConnOpt{}
	connOpt.SetDefault()

	connhOpt := &ktcp.KConnHandleOpt{
		Handler:	khandler.NewKConnHandlerEchoClient(),
		Protocol:	&kprotocol.KProtocolEcho{},
	}

	robotOpt := &robot.ClientRobotOpt{
		RobotingInterval:	20,
	}

	container, err := kcontainer.NewKContainer(2000)
	if nil != err {
		klog.LogWarn("Failed to create container : %s", err.Error())
		return
	}

	for i := uint64(0) ; i < 500 ; i++ {

		client, err := ktcp.NewKClient(i, cliOpt, connOpt, connhOpt )
		if nil != err {
			klog.LogWarn("Create KClient failed %v", err.Error())
			return
		}

		robot, err := robot.NewClientRobot(client, robotOpt)
		if nil != err {
			klog.LogWarn("Create Robot failed %v", err.Error())
			return
		}

		container.Add(robot)

	}

	cmd := ""
	fmt.Scanln(&cmd)

}

