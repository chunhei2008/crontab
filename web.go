package main

import (
	"bufio"
	"encoding/json"
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
	h = strings.TrimSpace(h)
	_, err = delJob(h)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", "success")
	}
}

func loger(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	d := r.FormValue("d")

	reg := regexp.MustCompile(`^[0-9]{8}$`)
	b := reg.MatchString(d)

	if !b {
		fmt.Fprintf(w, "%s", "invalid day")
		return
	}
	file := *logs + d + RUN_LOG_POSTFIX

	fp, err := os.Open(file)

	if err != nil {
		fmt.Fprintf(w, "%s", err)
		return
	}

	defer fp.Close()
	rd := bufio.NewReader(fp)
	rd.WriteTo(w)
}

func load(w http.ResponseWriter, r *http.Request) {
	loaded, loadErr := loadConf()
	if loaded {
		fmt.Fprintf(w, "%s", "success")
	} else {
		fmt.Fprintf(w, "%s", loadErr)
	}

}

func status(w http.ResponseWriter, r *http.Request) {
	brunning, err := json.Marshal(runnings)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", brunning)
	}
}

func stop(w http.ResponseWriter, r *http.Request) {
	stopCh <- true
	fmt.Fprintf(w, "%s", "success")
}

func start(w http.ResponseWriter, r *http.Request) {
	startCh <- true
	fmt.Fprintf(w, "%s", "success")
}
