package api

import (
	"fmt"
	"mcrontab/master/config"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	Serv *http.Server
}

var (
	G_apiServer *ApiServer
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var (
		content []byte
	)

	content = []byte("<h1>Job Manage Of Golang</h1>")

	w.Write(content)
}

//任务列表
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		content []byte
	)

	content = []byte("<h1>handleJobList</h1>")

	w.Write(content)
}

//任务更新保存
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		content []byte
	)

	content = []byte("<h1>handleJobList</h1>")

	w.Write(content)
}

//任务删除
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		content []byte
	)

	content = []byte("<h1>handleJobDelete</h1>")

	w.Write(content)
}

func InitApiServer() (err error) {
	var (
		mux      *http.ServeMux
		listener net.Listener
		httpServ *http.Server
		addr     string
	)
	//路由配置
	mux = http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/jobs/list", handleJobList)
	mux.HandleFunc("/jobs/save", handleJobSave)
	mux.HandleFunc("/jobs/delete", handleJobDelete)
	addr = ":" + strconv.Itoa(config.G_config.HttpServerPort)
	//监听并启动服务
	if listener, err = net.Listen("tcp", addr); err != nil {
		return
	}
	fmt.Println("httpserver listen ", addr)
	httpServ = &http.Server{
		ReadTimeout:  time.Duration(config.G_config.HttpServerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.G_config.HttpServerWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{
		Serv: httpServ,
	}
	err = httpServ.Serve(listener)
	return
}
