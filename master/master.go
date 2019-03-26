package master

import (
	"flag"
	"runtime"
)

type Master struct {
	configFile string
}

func (this *Master) Run() {
	this.initArgs()
	this.initEnv()
}

//初始化命令行参数
func (this *Master) initArgs() {
	flag.StringVar(&this.configFile, "config", "config/config.json", "master配置文件")
	flag.Parse()
}

//初始化线程数量
func (this *Master) initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
