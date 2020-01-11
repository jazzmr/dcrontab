package master

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/jazzmr/dcrontab/common"
	"time"
)

type JobMgr struct {
	client *clientv3.Client
	lease  clientv3.Lease
}

var (
	jobMgr *JobMgr
)

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
	)

	config = clientv3.Config{
		Endpoints:   GetConfig().EtcdEndpoints,
		DialTimeout: time.Duration(GetConfig().EtcdDialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	jobMgr = &JobMgr{
		client: client,
		lease:  client.Lease,
	}
	return
}

func GetJobMgr() *JobMgr {
	return jobMgr
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// /cron/jobs/jobName -> json
	var (
		jobKey   string
		jobValue []byte
		putRes   *clientv3.PutResponse
	)
	jobKey = common.JOB_STORE_PREFIX + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return
	}

	// 保存到etcd
	if putRes, err = jobMgr.client.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putRes.PrevKv != nil {
		// 忽略旧值获取错误
		if e := json.Unmarshal(putRes.PrevKv.Value, &oldJob); e != nil {
			fmt.Println(e)
		}
	}
	return
}

// 删除任务
func (jobMgr *JobMgr) DeleteJob(name string) (delJob *common.Job, err error) {

	var (
		jobKey string
		delRes *clientv3.DeleteResponse
	)
	jobKey = common.JOB_STORE_PREFIX + name

	if delRes, err = jobMgr.client.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if delRes.Deleted > 0 {
		if e := json.Unmarshal(delRes.PrevKvs[0].Value, &delJob); e != nil {
			fmt.Println(e)
		}
	}
	return
}

func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job) {
	var (
		getRes *clientv3.GetResponse
		err    error
		job    *common.Job
	)

	if getRes, err = jobMgr.client.Get(context.TODO(), common.JOB_STORE_PREFIX, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kv := range getRes.Kvs {
		job = new(common.Job)
		if err = json.Unmarshal(kv.Value, job); err == nil {
			jobList = append(jobList, job)
		}
	}
	return
}
