package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/*
*  get	获取任务列表
*  set	设置任务/添加任务
*  del	删除任务
*  log	任务执行日志
*  nxt	下一分钟要执行的任务
*  load 重新加载任务列表
*  stop 停止任务执行，等待正在执行的任务退出
 */

func get(w http.ResponseWriter, r *http.Request) {
	allJobs, err := getJobs()
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", allJobs)
	}
}

func set(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	h := r.FormValue("h")
	j := r.FormValue("j")
	j = strings.TrimSpace(j)
	if j == "" {
		fmt.Fprintf(w, "%s", "job empty")
	}
	_, err = setJob(h, j)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", "success")
	}
}

func del(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	h := r.FormValue("h")
	_, err = delJob(h)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", "success")
	}
}

func log(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	d := r.FormValue("d")

	reg := regexp.MustCompile(`^[0-9]{8}$`)
	b := reg.MatchString(d)

	if !b {
		fmt.Println("time err")
	}

	fp, err := os.Open("./web.go")
	if err != nil {
		fmt.Println("ERR")
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)
	rd.WriteTo(w)
}

func load(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is load")
	//TODO reload crontab.conf
}

func stop(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is stop")
	// TODO stop all jobs
}
