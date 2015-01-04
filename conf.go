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
var lock *sync.RWMutex = new(sync.RWMutex)
var jobs map[string]job = map[string]job{}
var regtime *regexp.Regexp = regexp.MustCompile(`^((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)\s+((\*(/[0-9]+)?)|[0-9\-\,/]+)$`)

func loadConf() (bool, error) {
	sysLog.Println("Load config start ...")
	tjobs := make(map[string]job, 20)
	fp, err := os.Open(*conf)
	if err != nil {
		sysLog.Printf("Err %s .\n", err)
		return false, err
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)

	for {
		line, rdErr := rd.ReadString('\n')

		if rdErr != nil && rdErr != io.EOF {
			sysLog.Println("Err %s.\n", rdErr)
			return false, rdErr
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if rdErr == io.EOF {
				break
			}
			continue
		}
		decode := json.NewDecoder(strings.NewReader(line))
		var j job
		if decErr := decode.Decode(&j); decErr != nil {
			sysLog.Printf("Err %s %s.\n", decErr, line)
			return false, decErr
		}
		_, parseErr := parseTime(&j)
		if parseErr != nil {
			sysLog.Printf("Err %s %s.\n", parseErr, line)
			return false, parseErr
		} else {
			h := md5.New()
			io.WriteString(h, line)
			hsum := fmt.Sprintf("%x", h.Sum(nil))
			tjobs[hsum] = j
		}
	}
	lock.Lock()
	defer lock.Unlock()
	jobs = tjobs
	sysLog.Println("Load config end.")
	return true, nil
}

func flushConf() (bool, error) {
	sysLog.Println("Flush config start ...")
	fp, err := os.Create(*conf)
	if err != nil {
		sysLog.Println(err)
		return false, err
	}
	defer fp.Close()
	for _, j := range jobs {
		b, _ := json.Marshal(j)
		fmt.Fprintf(fp, "%s\n", b)
	}
	sysLog.Println("Flush config end.")
	return true, nil
}

func parseTime(j *job) (bool, error) {
	if !regtime.MatchString(j.Time) {
		return false, errors.New("Time error")
	}

	respace := regexp.MustCompile(`\s+`)
	restar := regexp.MustCompile(`\*+`)
	r1 := restar.ReplaceAllString(j.Time, "*")
	r2 := respace.ReplaceAllString(r1, " ")
	r3 := strings.SplitN(r2, " ", -1)
	if len(r3) != 5 {
		return false, errors.New("Time error")
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
