package master

import (
	"encoding/json"
	"github.com/jazzmr/dcrontab/common"
	"io/ioutil"
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
		err        error
		job        common.Job
		inputBytes []byte
		oldJob     *common.Job
		bytes      []byte
	)
	// 解析表单
	//if err = r.ParseForm(); err != nil {
	//	goto ERR
	//}

	inputBytes, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err = json.Unmarshal(inputBytes, &job); err != nil {
		goto ERR
	}

	if oldJob, err = GetJobMgr().SaveJob(&job); err != nil {
		goto ERR
	}
	// 返回正常应答

	if bytes, err = common.BuildResponse(1, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	// 返回异常应答
	if bytes, err = common.BuildResponse(0, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}
