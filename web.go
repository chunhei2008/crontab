package main

import (
	"fmt"
	"net/http"
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
	fmt.Fprintf(w, "%s", "this is get")
}

func set(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is set")

}

func del(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is del")

}

func log(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is log")
}

func nxt(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is nxt")
}

func load(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is load")
}

func stop(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "this is stop")
}
