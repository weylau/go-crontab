package app

import (
	"flag"
	"fmt"
	"mcrontab/worker/config"
	"mcrontab/worker/executor"
	"mcrontab/worker/job"
	"mcrontab/worker/scheduler"
	"runtime"
)

type App struct {
}

var (
	configFile string
)

func (this *App) Run() {
	var (
		err error
	)

	if err = this.initApp(); err != nil {
		goto Err
	}
Err:
	fmt.Println(err)
}

func (this *App) initApp() (err error) {
	initArgs()
	initEnv()
	//初始化配置
	if err = initConfig(); err != nil {
		return
	}

	//启动任务执行器
	if err = initExecutor(); err != nil {
		return
	}
	//初始化调度器
	if err = initScheduler(); err != nil {
		return
	}
	//初始化任务管理器
	if err = initJobManager(); err != nil {
		return
	}

	return
}

//初始化命令行参数
func initArgs() {
	flag.StringVar(&configFile, "config", "src/mcrontab/worker/config/config.json", "master配置文件")
	flag.Parse()
}

//初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

//加载配置文件
func initConfig() (err error) {
	if err = config.LoadConfig(configFile); err != nil {
		return
	}
	return
}

func initJobManager() (err error) {
	err = job.InitJobManager()
	return
}

func initScheduler() (err error) {
	err = scheduler.InitScheduler()
	return
}

func initExecutor() (err error) {
	err = executor.InitExecutor()
	return
}

func NewApp() *App {
	return &App{}
}
