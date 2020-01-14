package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

//
func main() {
	ui.Main(setUp)
}

//
func setUp() {
	// 初始化窗口
	mainwin := ui.NewWindow("挂机软件", 640, 480, false)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	//tab
	tab := ui.NewTab()
	tab.Append("模拟点击", makeBasicControlsPage(mainwin))
	tab.SetMargined(0, true)
	// tab.Append("第二页", newBox())
	// tab.SetMargined(1, true)
	// 设置tab页
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	// 最后显示
	mainwin.Show()
}

//
func recharge() {
	begin := time.Now().Unix()
	//
	account := "748160424"
	clickEmpty := "input tap 168 852"
	clickDnf := "input tap 10 1430"
	clickAccount := "input tap 168 340 "
	deleteAccount := "input keyevent 67 "
	inputAccount := "input text " + account
	clickMoney := "input tap 168 952"
	inputMoney := "input text 100"
	clickPay := "input tap 650 1452"
	inputSecret := "input text 891210"
	keyBack := "input keyevent 4"
	//
	execCommandRun(clickDnf)
	time.Sleep(500)
	execCommandRun(clickAccount)
	for i := 0; i <= len(account); i++ {
		execCommandRun(deleteAccount)
	}
	execCommandRun(inputAccount)
	execCommandRun(clickEmpty)
	execCommandRun(clickMoney)
	execCommandRun(inputMoney)
	execCommandRun(clickEmpty)
	execCommandRun(clickPay)
	time.Sleep(3 * time.Second)
	fmt.Println("开始支付")
	execCommandRun(inputSecret)
	//
	time.Sleep(time.Second)
	for i := 0; i < 2; i++ {
		execCommandRun(keyBack)
	}
	fmt.Println(time.Now().Unix() - begin)
}

// 执行命令
func execCommandRun(cmd string) error {
	c := exec.Command("adb", "shell", cmd)
	err := c.Run()
	fmt.Println(err)
	return err
}
