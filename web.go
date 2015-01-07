package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
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
*  load 重新加载任务列表
*  stop 停止任务触发，正在执行的任务正常执行
*  start开始任务触发，
*  status获取正在执行的任务
 */

func get(w http.ResponseWriter, r *http.Request) {
	allJobs, err := configJobs.json()
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
	decode := json.NewDecoder(strings.NewReader(j))
	var jj job
	if decerr := decode.Decode(&jj); decerr != nil {
		fmt.Fprintf(w, "%s", decerr)
		return
	}
	parseTime(&jj)
	if h == "" {
		md5er := md5.New()
		io.WriteString(md5er, j)
		h = fmt.Sprintf("%x", md5er.Sum(nil))
	}

	configJobs.add(h, &jj)
	_, err = flushConf()
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
	configJobs.del(h)
	_, err = flushConf()
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
	file := *logs + d + "_" + RUN_LOG_POSTFIX

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
	brunning, err := runningJobs.json()
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
