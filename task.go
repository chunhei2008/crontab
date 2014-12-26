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
	Time    string   `json:"time"`
	Cmd     string   `json:"cmd"`
	Args    []string `json:"args"`
	Out     string   `json:"out"`
	Comment string   `json:"comment"`
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
	return true, nil
}

func delJob(h string) (bool, error) {
	lock.Lock()
	defer lock.Unlock()
	if _, exists := jobs[h]; exists {
		delete(jobs, h)
		return true, nil
	}
	return false, errors.New("not exists")
}
