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
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)

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

	common.WriteResponse(w, 1, "success", oldJob)
	//httputil.WriteJSONResponse(w, 1, oldJob)
	return
ERR:
	// 返回异常应答
	common.WriteResponse(w, 0, err.Error(), nil)
}

// 删除任务接口
// GET /job/delete name=jobName
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		jobName string
		delJob  *common.Job
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	jobName = r.FormValue("name")

	if delJob, err = GetJobMgr().DeleteJob(jobName); err != nil {
		goto ERR
	}

	common.WriteResponse(w, 1, "success", delJob)

	return
ERR:
	common.WriteResponse(w, 0, err.Error(), nil)
}

func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		jobList []*common.Job
	)
	jobList = GetJobMgr().ListJobs()
	common.WriteResponse(w, 1, "success", jobList)
}
