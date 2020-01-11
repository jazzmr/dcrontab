package master

import (
	"github.com/coreos/etcd/clientv3"
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

func (jobMgr *JobMgr) aa() {

}
