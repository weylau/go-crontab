package httpserver

import (
	"mcrontab/master/config"
	"net"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	Serv *http.Server
}

var (
	G_httpServer *HttpServer
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var (
		content []byte
	)

	content = []byte("hello world")

	w.Write(content)
}

func InitHttpServer() (err error) {
	var (
		mux      *http.ServeMux
		listener net.Listener
		httpServ *http.Server
	)
	//路由配置
	mux = http.NewServeMux()
	mux.HandleFunc("/", handleIndex)

	//监听并启动服务
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.G_config.HttpServerPort)); err != nil {
		return
	}
	httpServ = &http.Server{
		ReadTimeout:  time.Duration(config.G_config.HttpServerReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.G_config.HttpServerWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_httpServer = &HttpServer{
		Serv: httpServ,
	}
	err = httpServ.Serve(listener)
	return
}
