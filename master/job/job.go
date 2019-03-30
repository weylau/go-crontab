package job

import (
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"mcrontab/common"
	"mcrontab/master/config"
	"time"
)

type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
}

var (
	G_jobMamager *JobManager
)

//查询任务列表
func (this *JobManager) List() (joblist []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)
	dirKey = common.JOB_SAVE_DIR

	if getResp, err = this.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}
	joblist = make([]*common.Job, 0)
	// 遍历所有任务, 进行反序列化
	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			err = nil
			continue
		}
		joblist = append(joblist, job)
	}
	return
}

//保存任务
func (this *JobManager) Save(job *common.Job) (oldjob *common.Job, err error) {
	var (
		key      string
		value    []byte
		oldvalue common.Job
		putResp  *clientv3.PutResponse
	)
	key = common.JOB_SAVE_DIR + job.JobName
	if value, err = json.Marshal(job); err != nil {
		return
	}

	if putResp, err = this.kv.Put(context.TODO(), key, string(value), clientv3.WithPrevKV()); err != nil {
		return
	}
	if putResp.PrevKv != nil {
		if err = json.Unmarshal([]byte(putResp.PrevKv.Value), oldvalue); err != nil {
			return
		}
		oldjob = &oldvalue
	}

	return
}

//删除job
func (this *JobManager) Delete(jobName string) (oldjob *common.Job, err error) {
	var (
		key      string
		delResp  *clientv3.DeleteResponse
		oldvalue common.Job
	)

	key = common.JOB_SAVE_DIR + jobName

	// 从etcd中删除它
	if delResp, err = this.kv.Delete(context.TODO(), key, clientv3.WithPrevKV()); err != nil {
		return
	}

	// 返回被删除的任务信息
	if len(delResp.PrevKvs) != 0 {
		// 解析一下旧值, 返回它
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldvalue); err != nil {
			err = nil
			return
		}
		oldjob = &oldvalue
	}
	return
}

func InitJobManager() (err error) {
	var (
		conf       clientv3.Config
		client     *clientv3.Client
		kv         clientv3.KV
		jobManager *JobManager
	)
	jobManager = &JobManager{}
	conf = clientv3.Config{
		Endpoints:   []string{config.G_config.EtcdServerAddr},
		DialTimeout: time.Duration(config.G_config.EtcdServerConnectTimeout) * time.Millisecond,
	}
	if client, err = clientv3.New(conf); err != nil {
		return
	}
	kv = clientv3.NewKV(client)
	jobManager.client = client
	jobManager.kv = kv
	G_jobMamager = jobManager
	return
}
