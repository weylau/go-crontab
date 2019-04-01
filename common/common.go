package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"golang.org/x/net/context"
	"time"
)

//任务
type Job struct {
	JobName  string `json:"job_name"`
	ShellCmd string `json:"shell_cmd"`
	CronExpr string `json:"cron_expr"`
}

//任务事件
type JobEvent struct {
	EventType int
	Job       *Job
}

//任务计划
type JobPlan struct {
	Job         *Job                 // 要调度的任务信息
	Expr        *cronexpr.Expression // 解析好的cronexpr表达式
	NextExcTime time.Time            // 下次调度时间
}

// 任务执行状态
type JobExecuteInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    // 任务command的context
	CancelFunc context.CancelFunc //  用于取消command执行的cancel函数
}

// 任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo // 执行状态
	Output      []byte          // 脚本输出
	Err         error           // 脚本错误原因
	StartTime   time.Time       // 启动时间
	EndTime     time.Time       // 结束时间
}

// HTTP接口应答
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	// 1, 定义一个response
	var (
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	// 2, 序列化json
	resp, err = json.Marshal(response)
	return
}

// 反序列化Job
func UnpackJob(value []byte) (ret *Job, err error) {
	var (
		job *Job
	)

	job = &Job{}
	if err = json.Unmarshal(value, job); err != nil {
		return
	}
	ret = job
	return
}

// 构造任务执行计划
func BuildJobPlan(job *Job) (jobPlan *JobPlan, err error) {
	var (
		expr *cronexpr.Expression
	)
	// 解析JOB的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}
	// 生成任务调度计划对象
	jobPlan = &JobPlan{
		Job:         job,
		Expr:        expr,
		NextExcTime: expr.Next(time.Now()),
	}
	return
}

// 构造执行状态信息
func BuildJobExecuteInfo(jobPlan *JobPlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobPlan.Job,
		PlanTime: jobPlan.NextExcTime, // 计算调度时间
		RealTime: time.Now(),          // 真实调度时间
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
