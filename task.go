package main

/*
*  任务列表管理（添加，删除，更新）
 */
type task struct {
	hash        string //hash值
	prevRuntime int    //上次执行时间
	commit      string //备注
	time        string //crontab时间 * * * * *
	cmd         string //命令行
}

func getTask() {

}

func setTask() {

}

func delTask() {

}

func nxtTask() {

}
