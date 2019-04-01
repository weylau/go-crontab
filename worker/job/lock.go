package job

import (
	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
	"mcrontab/common"
)

type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancelFunc context.CancelFunc // 用于终止自动续租
	leaseId    clientv3.LeaseID   // 租约ID
	isLocked   bool               // 是否上锁成功
}

//上锁
func (this *JobLock) Lock() (err error) {

	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	//创建租约
	if leaseGrantResp, err = this.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	//定义取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	leaseId = leaseGrantResp.ID
	//自动续租
	if keepRespChan, err = this.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		return
	}

	//处理续租应答
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepResp == nil {
					return
				}
			}
		}

	}()

	//创建事务
	txn = this.kv.Txn(context.TODO())

	lockKey = common.JOB_LOCK_DIR + this.jobName
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	if txnResp, err = txn.Commit(); err != nil {
		goto ERR
	}

	//成功返回，失败释放锁
	if !txnResp.Succeeded {
		err = common.ERR_LOCK_ALREADY_REQUIRED
		goto ERR
	}
	//抢锁成功
	this.leaseId = leaseId
	this.cancelFunc = cancelFunc
	this.isLocked = true
	return
ERR:
	//取消自动续租
	cancelFunc()
	//释放租约
	this.lease.Revoke(context.TODO(), leaseId)
	return
}

//解锁
func (this *JobLock) Unlock() (err error) {
	if this.isLocked {
		this.cancelFunc()
		this.lease.Revoke(context.TODO(), this.leaseId)
	}
	return
}

// 初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}
