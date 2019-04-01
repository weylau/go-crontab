package job

import (
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"golang.org/x/net/context"
	"mcrontab/common"
	"mcrontab/master/config"
	"mcrontab/worker/scheduler"
	"time"
)

type JobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	G_jobMamager *JobManager
)

//创建任务的锁
func (this *JobManager) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, this.kv, this.lease)
	return
}

//取出所有并监视所有job的更新删除操作
func watchJobs() {
	var (
		err      error
		dirKey   string
		getResp  *clientv3.GetResponse
		kvPair   *mvccpb.KeyValue
		job      *common.Job
		ctx      context.Context
		jobEvent *common.JobEvent
	)
	dirKey = common.JOB_SAVE_DIR
	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(config.G_config.EtcdServerOptionTimeout)*time.Millisecond)
	if getResp, err = G_jobMamager.kv.Get(ctx, dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	//把取出的任务放到调度器的事件处理队列
	for _, kvPair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvPair.Value); err == nil {
			//更新事件
			jobEvent = &common.JobEvent{
				EventType: common.JOB_EVENT_SAVE,
				Job:       job,
			}
			scheduler.G_scheduler.PushJobEvent(jobEvent)
		}
	}
}

func InitJobManager() (err error) {
	var (
		conf       clientv3.Config
		client     *clientv3.Client
		lease      clientv3.Lease
		watcher    clientv3.Watcher
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
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	jobManager.client = client
	jobManager.kv = kv
	jobManager.lease = lease
	jobManager.watcher = watcher

	G_jobMamager = jobManager
	return
}
