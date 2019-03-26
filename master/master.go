package master

import (
	"flag"
	"fmt"
	"mcrontab/master/config"
	"runtime"
)

type Master struct {
	configFile string
}

var (
	G_config *config.Config
)

func (this *Master) Run() {
	var (
		err error
	)

	if err = this.Init(); err != nil {
		goto Err
	}
Err:
	fmt.Println(err)
}

func (this *Master) Init() (err error) {
	this.initArgs()
	this.initEnv()
	if G_config, err = config.LoadConfig(this.configFile); err != nil {
		return
	}
	return
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

func (this *Master) InitHttpServer() (err error) {

	return
}
