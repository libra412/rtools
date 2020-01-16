package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/astaxie/beego/toolbox"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"time"
)

var log = logrus.New()

//
func main() {
	log.Out = os.Stdout //日志标准输出
	file, err := os.OpenFile("wx.log", os.O_CREATE|os.O_WRONLY, 1)
	if err == nil {
		log.Out = file
	} else {
		log.Info("failed to log to file")
	}
	tk := toolbox.NewTask("tk", "0/2 * * * * ?", f)
	toolbox.AddTask("tk", tk)
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
	isRequest := false
	requestButton := ui.NewButton("开始接单")
	entryForm.Append("", requestButton, false)
	//
	supInput := ui.NewEntry()
	entryForm.Append("SUP商户号", supInput, false)
	supSecretInput := ui.NewEntry()
	entryForm.Append("SUP商户密钥", supSecretInput, false)
	productIdInput := ui.NewEntry()
	entryForm.Append("商品ID", productIdInput, false)

	secretInput := ui.NewEntry()
	entryForm.Append("输入支付密码", secretInput, false)
	accountInput := ui.NewEntry()
	entryForm.Append("输入充值账号", accountInput, false)
	moneyInput := ui.NewEntry()
	entryForm.Append("输入充值金额", moneyInput, false)
	orderIdInput := ui.NewEntry()
	entryForm.Append("输入订单号", orderIdInput, false)

	testButton := ui.NewButton("手动充值")
	entryForm.Append("", testButton, false)
	//
	testButton.OnClicked(func(*ui.Button) {
		account := accountInput.Text()
		secret := secretInput.Text()
		money := moneyInput.Text()
		orderId := orderIdInput.Text()
		if len(account) == 0 {
			ui.MsgBoxError(w, "错误提示", "账号不能为空")
			return
		}
		if len(secret) == 0 {
			ui.MsgBoxError(w, "错误提示", "密码不能为空")
			return
		}
		if len(money) == 0 {
			ui.MsgBoxError(w, "错误提示", "金额不能为空")
		}
		if len(orderId) == 0 {
			orderId = "123123123"
		}
		go rechargeDNF(account, secret, money, orderId)

	})
	//
	requestButton.OnClicked(func(*ui.Button) {
		if isRequest {
			isRequest = false
			toolbox.StopTask()
			requestButton.SetText("开始接单")
		} else {
			isRequest = true
			toolbox.StartTask()
			requestButton.SetText("停止接单")
		}
	})
	return hbox
}

// 自动接单
func f() error {
	fmt.Println("请求数据")
	data := getApiData()
	count := len(data)
	for i := 0; i < count; i++ {

	}
	return nil
}

//
func rechargeDNF(account, secret, money, orderId string) {
	begin := time.Now().Unix()
	//
	clickEmpty := "input tap 168 852"
	clickDnf := "input tap 10 1430"
	//clickAccount := "input tap 300 440 "
	//deleteAccount := "input keyevent --longpress 67 "
	oneAccount := "input tap 850 396"
	inputAccount := "input text " + account
	clickMoney := "input tap 168 952"
	inputMoney := "input text " + money
	clickPay := "input tap 650 1452"
	inputSecret := "input text " + secret
	keyBack := "input keyevent 4"
	fileName := orderId + ".png"
	screencapImage := "screencap -p /data/local/tmp/" + fileName
	copyImage := "/data/local/tmp/" + fileName
	desImage := "./" + fileName
	//
	execCommandRun(clickDnf)
	time.Sleep(time.Second)
	execCommandRun(oneAccount)
	execCommandRun(oneAccount)
	//for i := 0; i <= len(account)/2; i++ {
	//	execCommandRun(deleteAccount)
	//	// robotgo.KeyTap("del")
	//}
	execCommandRun(inputAccount)
	// robotgo.TypeString(account)
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
	execCommand(copyImage, desImage)
	res, err := execCommandSuccess(desImage)
	if err == nil {
		if strings.Contains(res, "成功") {
			for i := 0; i < 2; i++ {
				execCommandRun(keyBack)
			}
			log.Info("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，处理成功")
		} else {
			log.Error("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，错误信息：", "卡单，可疑")
		}
	} else {
		log.Error("订单号:", orderId, "，耗时：", time.Now().Unix()-begin, "，充值账号：", account, "，充值金额：", money, "，错误信息：", err)
	}
}

// shell执行命令
func execCommandRun(cmd string) error {
	c := exec.Command("adb", "shell", cmd)
	err := c.Run()
	fmt.Println(err)
	return err
}

// 执行命令
func execCommand(org, des string) error {
	c := exec.Command("adb", "pull", org, des)
	err := c.Start()
	fmt.Println("download image", err)
	return err
}

// 识别图片
func execCommandSuccess(fileName string) (string, error) {
	c := exec.Command("tesseract", fileName, "stdout", "-l", "chi_sim")
	res, err := c.Output()
	return string(res), err
}

//
type DnfData struct {
}

// 获取API数据
func getApiData() []DnfData {

	return nil
}
