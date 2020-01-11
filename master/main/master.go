package main

import (
	"flag"
	"fmt"
	"github.com/jazzmr/dcrontab/master"
	"runtime"
)

var (
	banner = `
██████╗  ██████╗██████╗  ██████╗ ███╗   ██╗████████╗ █████╗ ██████╗ 
██╔══██╗██╔════╝██╔══██╗██╔═══██╗████╗  ██║╚══██╔══╝██╔══██╗██╔══██╗
██║  ██║██║     ██████╔╝██║   ██║██╔██╗ ██║   ██║   ███████║██████╔╝
██║  ██║██║     ██╔══██╗██║   ██║██║╚██╗██║   ██║   ██╔══██║██╔══██╗
██████╔╝╚██████╗██║  ██║╚██████╔╝██║ ╚████║   ██║   ██║  ██║██████╔╝
╚═════╝  ╚═════╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═════╝
`
	confFile string
)

func main() {
	var (
		err error
	)

	initArgs()

	// 初始化线程
	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 初始化ApiServer
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}
	fmt.Print(banner)
	return
ERR:
	fmt.Println(err)
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	flag.StringVar(&confFile, "config", "master/main/master.json", "指定master json")
	flag.Parse()
}
