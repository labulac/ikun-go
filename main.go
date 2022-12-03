package main

import (
	"flag"
	"github.com/kardianos/service"
	hook "github.com/robotn/gohook"
	"ikun-go/player"
	"log"
	"os/exec"
	"strings"
	"time"
)

// 服务管理
var serviceType = flag.String("s", "", "服务管理, 支持install, uninstall")

func main() {
	flag.Parse()

	switch *serviceType {
	case "install":
		installService()
	case "uninstall":
		uninstallService()
	default:

		s := getService()
		status, _ := s.Status()
		if status != service.StatusUnknown {
			// 以服务方式运行
			s.Run()
		} else {
			// 非服务方式运行
			switch s.Platform() {
			case "windows-service":
				log.Println("可使用 .\\ikun-go.exe -s install 安装服务运行")
			default:
				log.Println("可使用 sudo ./ikun-go -s install 安装服务运行")
			}
			jntm()
		}
	}

}

type program struct{}

func (q *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go q.run()
	return nil
}
func (q *program) run() {
	log.Println("服务运行")
	jntm()
}
func (q *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func getService() service.Service {
	options := make(service.KeyValue)
	if service.ChosenSystem().String() == "unix-systemv" {
		options["SysvScript"] = sysvScript
	}

	svcConfig := &service.Config{
		Name:        "ikun-go",
		DisplayName: "ikun-go",
		Description: "author:labulac@88.com",
		Option:      options,
		Arguments:   []string{},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}

// 卸载服务
func uninstallService() {
	s := getService()
	s.Stop()
	if service.ChosenSystem().String() == "unix-systemv" {
		if _, err := exec.Command("/etc/init.d/ikun-go", "stop").Output(); err != nil {
			log.Println(err)
		}
	}
	if err := s.Uninstall(); err == nil {
		log.Println("ikun-go 服务卸载成功!")
	} else {
		log.Printf("ikun-go 服务卸载失败, ERR: %s\n", err)
	}
}

// 安装服务
func installService() {
	s := getService()

	status, err := s.Status()
	if err != nil && status == service.StatusUnknown {
		// 服务未知，创建服务
		if err = s.Install(); err == nil {
			s.Start()

			log.Println("安装 ikun-go 服务成功!")
			if service.ChosenSystem().String() == "unix-systemv" {
				if _, err := exec.Command("/etc/init.d/ikun-go", "enable").Output(); err != nil {
					log.Println(err)
				}
				if _, err := exec.Command("/etc/init.d/ikun-go", "start").Output(); err != nil {
					log.Println(err)
				}
			}
			return
		}

		log.Printf("安装 ikun-go 服务失败, ERR: %s\n", err)
	}

	if status != service.StatusUnknown {
		log.Println("ikun-go 服务已安装, 无需再次安装")
	}
}

var lastNum int

func paintedEggShell(num int) {
	if num == 1 {
		lastNum = 0
	}
	if lastNum+1 == num {
		lastNum = num
		if lastNum == 4 {
			log.Println("小黑子,露出鸡脚了吧!")
			time.Sleep(time.Second * 1)
			go player.PlaySound("JNTM")
			lastNum = 0
		}
	} else {
		lastNum = 0
	}
}

func jntm() {
	// 监听键盘事件
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		key := strings.ToUpper(string(e.Keychar))
		log.Println("PRESS " + key)
		switch key {
		case "J":
			go player.PlaySound("J")
			paintedEggShell(1)
		case "N":
			go player.PlaySound("N")
			paintedEggShell(2)
		case "T":
			go player.PlaySound("T")
			paintedEggShell(3)
		case "M":
			go player.PlaySound("M")
			paintedEggShell(4)
		}

	})
	s := hook.Start()
	<-hook.Process(s)
}

const sysvScript = `#!/bin/sh /etc/rc.common
DESCRIPTION="{{.Description}}"
cmd="{{.Path}}{{range .Arguments}} {{.|cmd}}{{end}}"
name="ikun-go"
pid_file="/var/run/$name.pid"
stdout_log="/var/log/$name.log"
stderr_log="/var/log/$name.err"
START=99
get_pid() {
    cat "$pid_file"
}
is_running() {
    [ -f "$pid_file" ] && cat /proc/$(get_pid)/stat > /dev/null 2>&1
}
start() {
	if is_running; then
		echo "Already started"
	else
		echo "Starting $name"
		{{if .WorkingDirectory}}cd '{{.WorkingDirectory}}'{{end}}
		$cmd >> "$stdout_log" 2>> "$stderr_log" &
		echo $! > "$pid_file"
		if ! is_running; then
			echo "Unable to start, see $stdout_log and $stderr_log"
			exit 1
		fi
	fi
}
stop() {
	if is_running; then
		echo -n "Stopping $name.."
		kill $(get_pid)
		for i in $(seq 1 10)
		do
			if ! is_running; then
				break
			fi
			echo -n "."
			sleep 1
		done
		echo
		if is_running; then
			echo "Not stopped; may still be shutting down or shutdown may have failed"
			exit 1
		else
			echo "Stopped"
			if [ -f "$pid_file" ]; then
				rm "$pid_file"
			fi
		fi
	else
		echo "Not running"
	fi
}
restart() {
	stop
	if is_running; then
		echo "Unable to stop, will not attempt to start"
		exit 1
	fi
	start
}
`
