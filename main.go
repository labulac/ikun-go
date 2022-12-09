package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	hook "github.com/robotn/gohook"
	"ikun-go/player"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// 服务管理
var serviceType = flag.String("s", "", "服务管理, 支持install, uninstall")

var port = "9882"

func main() {
	flag.Parse()

	switch *serviceType {
	case "install":
		err := installService(true)
		if err != nil {
			log.Println(err)
		}

		return
	case "uninstall":
		err := installService(false)
		if err != nil {
			log.Println(err)
		}
		return
	default:
		log.Println("使用 -s install参数增加开机自启动")
		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		go func() {
			err := r.Run(":" + port)
			if err != nil {
				panic(err)
			}
		}()
		jntm()

	}

}

const (
	winBat = `start %s`
)

func installService(on bool) error {
	var err error
	var path, content string
	current, err := user.Current()
	if err != nil {
		return err
	}
	switch runtime.GOOS {

	case "windows":
		path = fmt.Sprintf("%s\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\ikun-go.bat", current.HomeDir)
		abs, _ := filepath.Abs(os.Args[0])
		//log.Println(abs)
		content = fmt.Sprintf(winBat, abs)
		if on {
			exec.Command(abs).Start()
		}

	default:
		return errors.New("不支持的系统")
	}
	return writer(on, path, content)

}

func writer(on bool, path, content string) error {
	var err error
	if on {
		stat, _ := os.Stat(path)
		if stat == nil {
			err = os.WriteFile(path, []byte(content), os.ModePerm)
		}
	} else {
		err = os.Remove(path)
		switch runtime.GOOS {

		case "windows":
			cmd := exec.Command("tasklist")
			output, err := cmd.Output()
			//log.Println(string(output))
			if err != nil {
				fmt.Println(err)
				return err
			}

			// 查找进程ID
			pid := findPID(output)
			if pid == -1 {
				fmt.Println("找不到使用该端口的进程")
				return nil
			}

			// 杀死进程
			cmd = exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
			err = cmd.Run()
			if err != nil {
				fmt.Println(err)
				return nil
			}
		default:
			command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
			execCmd(exec.Command("bash", "-c", command))
		}

	}
	return err
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
	log.Println("开始")
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

// Execute command and return exited code.
func execCmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Error during killing (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Port successfully killed (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}

func findPID(output []byte) int {
	// 在输出中查找进程ID
	// 这里假设输出是以空格分隔的列表，其中第一列是进程名称，第二列是进程ID
	// 例如：
	// process1  1234
	// process2  5678
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		if fields[0] == "ikun-go.exe" {
			pid, err := strconv.Atoi(fields[1])
			if err != nil {
				continue
			}
			return pid
		}
	}

	return -1
}
