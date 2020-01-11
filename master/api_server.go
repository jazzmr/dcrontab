package master

import (
	"encoding/json"
	"github.com/jazzmr/dcrontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	httpServer *apiServer
	//httpServerOnce sync.Once
)

// 任务的HTTP接口
type apiServer struct {
	server *http.Server
}

func InitApiServer() (err error) {
	var (
		mux      *http.ServeMux
		listener net.Listener
		server   *http.Server
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(GetConfig().ApiPort)); err != nil {
		return
	}

	server = &http.Server{
		Handler:      mux,
		ReadTimeout:  time.Duration(GetConfig().ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(GetConfig().ApiWriteTimeout) * time.Millisecond,
	}

	go func() {
		if e := server.Serve(listener); e != nil {
			panic(e)
		}
	}()
	//httpServerOnce.Do(func() {
	httpServer = &apiServer{server: server}
	//})
	return
}

func GetApiServer() *apiServer {
	return httpServer
}

// 保存任务接口
// POST job={"name"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err       error
		job       common.Job
		formValue string
	)
	// 解析表单
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	formValue = r.PostFormValue("job")
	if err = json.Unmarshal([]byte(formValue), &job); err != nil {
		goto ERR
	}

	return
ERR:
}
