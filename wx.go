package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/go-vgo/robotgo"
	"os/exec"
	"time"
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
	tab.Append("DNF", makeControl(mainwin))
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
func makeControl(w *ui.Window) ui.Control {
	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	hbox.Append(entryForm, false)
	// 设置发送请求按钮
	requestButton := ui.NewButton("开始接单")
	entryForm.Append("", requestButton, false)
	//
	secretInput := ui.NewEntry()
	entryForm.Append("输入密码", secretInput, false)
	accountInput := ui.NewEntry()
	entryForm.Append("输入测试账号", accountInput, false)
	testButton := ui.NewButton("测试")
	entryForm.Append("", testButton, false)
	//
	testButton.OnClicked(func(*ui.Button) {
		account := accountInput.Text()
		secret := secretInput.Text()
		if len(account) == 0 {
			ui.MsgBoxError(w, "错误提示", "账号不能为空")
			return
		}
		if len(secret) == 0 {
			ui.MsgBoxError(w, "错误提示", "密码不能为空")
			return
		}
		recharge(account, secret, "100", "123123123")
	})
	return hbox
}

//
func recharge(account, secret, money, orderId string) {
	begin := time.Now().Unix()
	//
	clickEmpty := "input tap 168 852"
	clickDnf := "input tap 10 1430"
	clickAccount := "input tap 300 440 "
	deleteAccount := "input keyevent 67 "
	inputAccount := "input text " + account
	clickMoney := "input tap 168 952"
	inputMoney := "input text " + money
	clickPay := "input tap 650 1452"
	inputSecret := "input text " + secret
	keyBack := "input keyevent 4"
	fileName := orderId + ".png"
	screencapImage := "/system/bin/screencap -p /data/local/tmp/" + fileName
	copyImage := "pull /data/local/tmp/" + fileName + " ./" + fileName
	//
	execCommandRun(clickDnf)
	time.Sleep(time.Second)
	execCommandRun(clickAccount)
	for i := 0; i <= len(account); i++ {
		// execCommandRun(deleteAccount)
		robotgo.KeyTap("backspace")
	}
	// execCommandRun(inputAccount)
	robotgo.TypeString(inputAccount)
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
	execCommandRun(screencapImage)
	execCommand(copyImage)
	for i := 0; i < 2; i++ {
		execCommandRun(keyBack)

	}
	fmt.Println(time.Now().Unix() - begin)
}

// shell执行命令
func execCommandRun(cmd string) error {
	c := exec.Command("adb", "shell", cmd)
	err := c.Run()
	fmt.Println(err)
	return err
}

// 执行命令
func execCommand(cmd string) error {
	c := exec.Command("adb", cmd)
	err := c.Run()
	return err
}
