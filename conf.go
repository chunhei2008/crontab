package main

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

/*
* 任务配置文件，读取&更新
 */
var lock sync.RWMutex
var jobs map[string]job

func loadConf() {
	sysLog.Println("load config start ...")
	lock.Lock()
	defer lock.Unlock()
	jobs = make(map[string]job)
	fp, err := os.Open(*conf)
	if err != nil {
		sysLog.Println(err)
		os.Exit(1)
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
			sysLog.Printf("%s:%s", decerr, line)

			break
		}
		_, parseErr := parseTime(&j)
		if parseErr != nil {
			sysLog.Printf("%s:%s", parseErr, line)
		} else {
			h := md5.New()
			io.WriteString(h, line)
			hsum := fmt.Sprintf("%x", h.Sum(nil))
			jobs[hsum] = j
		}
	}
	sysLog.Println("load config end .")
}

func flushConf() {
	sysLog.Println("flush config start ...")
	lock.RLock()
	defer lock.RUnlock()
	fp, err := os.Create(*conf)
	if err != nil {
		sysLog.Println(err)
	}
	defer fp.Close()
	for _, j := range jobs {
		b, _ := json.Marshal(j)
		fmt.Fprintf(fp, "%s\n", b)
	}
	sysLog.Println("flush config end .")
}

func parseTime(j *job) (bool, error) {
	regtime := regexp.MustCompile(`^((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)$`)

	if !regtime.MatchString(j.Time) {
		return false, errors.New("time error")
	}

	respace := regexp.MustCompile(`\s+`)
	restar := regexp.MustCompile(`\*+`)
	r1 := restar.ReplaceAllString(j.Time, "*")
	r2 := respace.ReplaceAllString(r1, " ")
	r3 := strings.SplitN(r2, " ", -1)
	if len(r3) != 5 {
		return false, errors.New("time error")
	} else {
		j.minute = parseNumber(r3[0], 0, 59)
		j.hour = parseNumber(r3[1], 0, 23)
		j.dom = parseNumber(r3[2], 1, 31)
		j.month = parseNumber(r3[3], 1, 12)
		j.dow = parseNumber(r3[4], 0, 6)
	}
	return true, nil
}

func parseNumber(s string, min, max int) []int {
	v := strings.SplitN(s, ",", -1)
	result := make([]int, 0)
	for _, vv := range v {
		if vv == "" {
			continue
		}
		vvv := strings.SplitN(vv, "/", -1)
		var step int
		if len(vvv) < 2 || vvv[1] == "" {
			step = 1
		} else {
			step, _ = strconv.Atoi(vvv[1])
		}
		vvvv := strings.SplitN(vvv[0], "-", -1)
		var _min, _max int
		if len(vvvv) == 2 {
			_min, _ = strconv.Atoi(vvvv[0])
			_max, _ = strconv.Atoi(vvvv[1])
		} else {
			if vvv[0] == "*" {
				_min = min
				_max = max
			} else {
				_min, _ = strconv.Atoi(vvv[0])
				_max = _min
			}
		}

		for i := _min; i <= _max; i += step {
			result = append(result, i)
		}
	}
	return result
}
