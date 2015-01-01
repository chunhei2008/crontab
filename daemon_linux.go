package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var d *bool = flag.Bool("d", true, "daemon")

func fork() (int, syscall.Errno) {
	r1, r2, err := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return 0, err
	}

	if runtime.GOOS == "darwin" && r2 == 1 {
		r1 = 0
	}

	return int(r1), 0
}

func daemon(chdir bool) {
	pid, err := fork()
	if err != 0 {
		fmt.Println("Daemon error!")
		os.Exit(1)
	}

	if pid < 0 {
		fmt.Println("Daemon error!")
		os.Exit(1)
	} else if pid > 0 {
		os.Exit(0)
	} else if pid == 0 {
		syscall.Setsid()
		if chdir {
			os.Chdir("/")
		}
		syscall.Umask(0)
		f, _ := os.Open("/dev/null")
		devnull := f.Fd()

		syscall.Dup2(int(devnull), int(os.Stdin.Fd()))
		syscall.Dup2(int(devnull), int(os.Stdout.Fd()))
		syscall.Dup2(int(devnull), int(os.Stderr.Fd()))
	}
}
