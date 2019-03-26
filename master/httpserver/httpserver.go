package httpserver

import "net/http"

type HttpServer struct {
	Serv *http.Server
}

func (this *HttpServer) InitHttpServer(err error) {

}
