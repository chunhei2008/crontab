package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
* 任务配置文件，读取&更新
 */

func loadConf() {
	fmt.Println("load conf")
	fp, err := os.Open(*conf)
	if err != nil {
		fmt.Print(err)
	}
	defer fp.Close()
	rd := bufio.NewReader(fp)

	for {
		line, err := rd.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		} else {
			fmt.Println(line)
			fmt.Println("##############")
			cmds := strings.Split(line, " ")
			if len(cmds) < 6 {
				os.Exit(1)
			}
			f := false
			for k, v := range cmds {
				switch {
				case k < 5:
					//校验时间设置
				case f == true:
					//组合 ""参数
				default:
					if v[0] == '"' {
						f = true
					}
				}
			}
		}
	}

	// read file
	fmt.Println("end load")

}

func flushConf() {
	//taskList := getTask()
	//TODO write into file
}
