package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	hook "github.com/robotn/gohook"
	"golang.org/x/sys/windows/registry"
	"ikun-go/player"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var port = "9882"
var serviceType = flag.String("s", "", "支持install,uninstall")

func main() {
	flag.Parse()
	switch *serviceType {
	case "install":
		log.Println("install")
		// 获取程序路径
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		// 创建注册表项
		key, _, err := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
		if err != nil {
			panic(err)
		}
		defer key.Close()

		// 写入程序路径
		err = key.SetStringValue(filepath.Base(exePath), exePath)
		if err != nil {
			panic(err)
		}

		fmt.Println("程序已设置开机自启动")
	case "uninstall":
		log.Println("uninstall")
		// 获取程序路径
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		// 打开注册表项
		key, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.ALL_ACCESS)
		if err != nil {
			panic(err)
		}
		defer key.Close()

		// 删除程序路径
		err = key.DeleteValue(filepath.Base(exePath))
		if err != nil {
			panic(err)
		}

		fmt.Println("程序已取消开机自启动")
	default:
		log.Println("default")
		start()
	}

}

func start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "cxk")
	})
	go func() {
		err := r.Run(":" + port)
		if err != nil {
			panic(err)
		}
	}()

	jntm()
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
