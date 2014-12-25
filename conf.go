package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

/*
* 任务配置文件，读取&更新
 */
var lock *sync.RWMutex
var jobs map[string]job

func loadConf() {
	jobs = make(map[string]job)

	fmt.Println("load conf")
	fp, err := os.Open(*conf)
	if err != nil {
		fmt.Print(err)
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)

	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		decode := json.NewDecoder(strings.NewReader(line))
		var j job
		if decerr := decode.Decode(&j); decerr != nil {
			break
		}
		h := md5.New()
		io.WriteString(h, line)
		hsum := fmt.Sprintf("%x", h.Sum(nil))
		jobs[hsum] = j
	}

	// read file
	fmt.Println("end load")

}

func flushConf() {
	//taskList := getTask()
	//TODO write into file
}
