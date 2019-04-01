package common

var (
	// 任务保存目录
	JOB_SAVE_DIR = "/cron/jobs/"
	// 任务锁目录
	JOB_LOCK_DIR = "/cron/lock/"
	/*************job事件定义**************/
	// 保存任务事件
	JOB_EVENT_SAVE = 1
	// 删除任务事件
	JOB_EVENT_DELETE = 2
	// 强杀任务事件
	JOB_EVENT_KILL = 3
)
