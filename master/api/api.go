package api

import (
	"encoding/json"
	"fmt"
	"mcrontab/common"
	"mcrontab/master/config"
	"mcrontab/master/job"
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
		err     error
		bytes   []byte
		joblist []*common.Job
	)
	if joblist, err = job.G_jobMamager.List(); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", joblist); err != nil {
		goto ERR
	}
	responseJson(w, bytes)
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		fmt.Println("err :", string(bytes))
		responseJson(w, bytes)
		return
	}

}

//任务更新保存
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		content []byte
		jobdata string
		newjob  common.Job
		oldjob  *common.Job
		bytes   []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	jobdata = r.PostForm.Get("job")
	newjob = common.Job{}
	content = []byte(jobdata)
	if err = json.Unmarshal(content, &newjob); err != nil {
		goto ERR
	}

	if oldjob, err = job.G_jobMamager.Save(&newjob); err != nil {
		goto ERR
	}

	if bytes, err = common.BuildResponse(0, "success", oldjob); err != nil {
		goto ERR
	}
	fmt.Println(string(bytes))
	responseJson(w, bytes)
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		responseJson(w, bytes)
	}
}

//任务删除
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error // interface{}
		name   string
		oldjob *common.Job
		bytes  []byte
	)

	// POST:   a=1&b=2&c=3
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	// 删除的任务名
	name = r.Form.Get("name")
	fmt.Println(name)
	// 去删除任务
	if oldjob, err = job.G_jobMamager.Delete(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", oldjob); err == nil {
		responseJson(w, bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		responseJson(w, bytes)
	}
}

func responseJson(w http.ResponseWriter, content []byte) (int, error) {
	w.Header().Add("Content-Type", "application/json;charset=UTF-8")
	return w.Write(content)
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
		IdleTimeout:  time.Duration(config.G_config.HttpServerIdleTimeout) * time.Millisecond,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{
		Serv: httpServ,
	}
	err = httpServ.Serve(listener)
	return
}
