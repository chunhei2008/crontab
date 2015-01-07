package main

import (
	"encoding/json"
	"sync"
	"time"
)

/*
*  任务列表管理（添加，删除，更新）
 */

type job struct {
	Time    string   `json:"time"`    //任务执行时间
	Cmd     string   `json:"cmd"`     //可执行程序
	Args    []string `json:"args"`    //执行参数
	Out     string   `json:"out"`     //输出文件
	Comment string   `json:"comment"` //任务备注
	Start   string   `json:"start"`   //任务单次执行仅用作状态使用
	minute  []int
	hour    []int
	dom     []int
	month   []int
	dow     []int
}

func NewJobs() *Jobs {
	return &Jobs{mj: make(map[string]*job), lk: new(sync.RWMutex)}
}

type Jobs struct {
	mj map[string]*job
	lk *sync.RWMutex
}

func (jobs *Jobs) add(k string, v *job) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	jobs.mj[k] = v
}

func (jobs *Jobs) del(k string) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	delete(jobs.mj, k)
}

func (jobs *Jobs) json() ([]byte, error) {
	jobs.lk.RLock()
	defer jobs.lk.RUnlock()
	return json.Marshal(jobs.mj)
}

func (jobs *Jobs) getJobs() map[string]*job {
	jobs.lk.RLock()
	defer jobs.lk.RUnlock()
	return jobs.mj
}

func (jobs *Jobs) replaceJobs(mj map[string]*job) {
	jobs.lk.Lock()
	defer jobs.lk.Unlock()
	jobs.mj = mj
}

func (jobs *Jobs) runJobs() {
	t := time.Now()
	if t.Second() == 0 {
		jobs.lk.Lock()
		defer jobs.lk.Unlock()
		minute := t.Minute()
		hour := t.Hour()
		dom := t.Day()
		month := int(t.Month())
		dow := int(t.Weekday())
		for _, j := range jobs.mj {
			if inArray(j.minute, minute) &&
				inArray(j.hour, hour) &&
				inArray(j.dom, dom) &&
				inArray(j.month, month) &&
				inArray(j.dow, dow) {
				go runJob(*j)
			}
		}
	}
}
