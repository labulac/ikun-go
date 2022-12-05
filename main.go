package main

import (
	"fmt"
	"github.com/kardianos/service"
	hook "github.com/robotn/gohook"
	"ikun-go/player"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type program struct {
}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {

	fmt.Println("fsdfasdfsd")

	//code here
	jntm()
	wg.Done()
}

func (p *program) Stop(s service.Service) error {
	return nil
}

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	svcConfig := &service.Config{
		Name:        "test-service",
		DisplayName: "test-service",
		Description: "test-service",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			x := s.Install()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("service installed")
			return
		} else if os.Args[1] == "uninstall" {
			x := s.Uninstall()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("service uninstall success")
			return
		}
	}
	err = s.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	wg.Wait()
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
