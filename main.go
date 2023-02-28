package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
	hook "github.com/robotn/gohook"
	"ikun-go/player"
	"log"
	"os"
	"strings"
	"time"
)

var port = "9882"
var logger = service.ConsoleLogger

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	log.Println("使用 install参数增加开机自启动")
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

func (p *program) Stop(s service.Service) error {
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "akun", //服务显示名称
		DisplayName: "akun", //服务名称
		Description: "微服务",  //服务描述
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error(err)
	}

	if err != nil {
		logger.Error(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			s.Install()
			logger.Info("服务安装成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "start":
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "stop":
			s.Stop()
			logger.Info("服务关闭成功!")
			break
		case "restart":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Start()
			logger.Info("服务启动成功!")
			break
		case "remove":
			s.Stop()
			logger.Info("服务关闭成功!")
			s.Uninstall()
			logger.Info("服务卸载成功!")
			break
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
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
