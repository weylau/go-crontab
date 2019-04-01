package scheduler

import (
	"fmt"
	"mcrontab/common"
	"mcrontab/worker/executor"
	"time"
)

type Scheduler struct {
	jobEventChan      chan *common.JobEvent             //任务事件队列
	jobPlanTable      map[string]*common.JobPlan        //任务计划表
	jobExecutingTable map[string]*common.JobExecuteInfo // 任务执行表
	jobResultChan     chan *common.JobExecuteResult
}

var (
	G_scheduler *Scheduler
)

//添加任务事件
func (this *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	this.jobEventChan <- jobEvent
}

//添加执行结果记录到队列
func (this *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	this.jobResultChan <- jobResult
}

//处理任务事件
func (this *Scheduler) HandlerJobEvent(jobEvent *common.JobEvent) {
	var (
		jobPlan    *common.JobPlan
		jobExisted bool
		err        error
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: // 保存任务事件
		if jobPlan, err = common.BuildJobPlan(jobEvent.Job); err != nil {
			return
		}
		this.jobPlanTable[jobEvent.Job.JobName] = jobPlan
	case common.JOB_EVENT_DELETE: // 删除任务事件
		if jobPlan, jobExisted = this.jobPlanTable[jobEvent.Job.JobName]; jobExisted {
			delete(this.jobPlanTable, jobEvent.Job.JobName)
		}
	case common.JOB_EVENT_KILL: // 强杀任务事件
	}
}

//获取最近一次执行的任务时间
func (this *Scheduler) nextJobRuntime() {
	var (
		jobName string
		jobPlan *common.JobPlan
		now     time.Time
	)
	now = time.Now()
	for jobName, jobPlan = range this.jobPlanTable {
		if jobPlan.NextExcTime.Before(now) || jobPlan.NextExcTime.Equal(now) {
			this.startJob(jobPlan)
		}
	}
}

//执行任务计划
func (this *Scheduler) startJob(jobPlan *common.JobPlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	if jobExecuteInfo, jobExecuting = this.jobExecutingTable[jobPlan.Job.JobName]; jobExecuting {
		//同一个任务只能同时执行一次
		return
	}

	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	this.jobExecutingTable[jobPlan.Job.JobName] = jobExecuteInfo
	//执行任务
	fmt.Println("执行任务：", jobExecuteInfo.Job.JobName, "计划执行时间：", jobExecuteInfo.PlanTime, "实际执行时间：", jobExecuteInfo.RealTime)
	executor.G_executor.ExecuteJob(jobExecuteInfo)
}

func (this *Scheduler) Run() {
	var (
		jobEvent *common.JobEvent
	)
	//获取最近一次执行的任务时间
	this.nextJobRuntime()

	for {
		select {
		//监听并处理任务事件
		case jobEvent = <-this.jobEventChan:
			this.HandlerJobEvent(jobEvent)
		}
	}
}

//初始化调度器
func InitScheduler() (err error) {
	var (
		scheduler *Scheduler
	)

	scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
	}

	G_scheduler = scheduler
	go G_scheduler.Run()
	return
}
