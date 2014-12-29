package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
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

func getJobs() ([]byte, error) {
	lock.RLock()
	defer lock.RUnlock()
	allJobs, err := json.Marshal(jobs)
	if err != nil {
		return nil, err
	}
	return allJobs, nil
}

func setJob(h string, j string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()
	decode := json.NewDecoder(strings.NewReader(j))
	var jj job
	if decerr := decode.Decode(&jj); decerr != nil {
		return false, decerr
	}
	parseTime(&jj)
	md5er := md5.New()
	io.WriteString(md5er, j)
	hsum := fmt.Sprintf("%x", md5er.Sum(nil))
	jobs[hsum] = jj
	return flushConf()
}

func delJob(h string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := jobs[h]; exists {
		delete(jobs, h)
		return flushConf()
	}
	return false, errors.New("not exists")
}
