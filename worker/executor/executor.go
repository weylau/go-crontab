package executor

import (
	"math/rand"
	"mcrontab/common"
	"mcrontab/worker/job"
	"mcrontab/worker/scheduler"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (this *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			err              error
			jobExecuteResult *common.JobExecuteResult
			jobLock          *job.JobLock
			cmd              *exec.Cmd
			output           []byte
		)
		jobExecuteResult = &common.JobExecuteResult{
			ExecuteInfo: info,
			StartTime:   time.Now(),
			Output:      make([]byte, 0),
		}
		//上锁
		jobLock = job.G_jobMamager.CreateJobLock(info.Job.JobName)
		//随机睡眠防止在分布式部署时总是一台主机在跑
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		//判断锁有没有被占用
		err = jobLock.Lock()
		defer jobLock.Unlock()
		//执行任务
		if err != nil {
			//上锁失败
			jobExecuteResult.Err = err
			jobExecuteResult.EndTime = time.Now()
		} else {
			jobExecuteResult.StartTime = time.Now()
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.ShellCmd)
			output, err = cmd.CombinedOutput()
			jobExecuteResult.Err = err
			jobExecuteResult.EndTime = time.Now()
			jobExecuteResult.Output = output
		}
		//添加执行结果
		scheduler.G_scheduler.PushJobResult(jobExecuteResult)
	}()
}

func InitExecutor() (err error) {
	var (
		executor *Executor
	)
	executor = &Executor{}
	G_executor = executor
	return
}
