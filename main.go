package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/libra/rtools"
	// "strings"
)

//
func main() {
	ui.Main(setUp)
}

//
func makeBasicControlsPage() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	// 设置发送请求按钮
	requestButton := ui.NewButton("发送请求")
	hbox.Append(requestButton, false)
	// hbox.Append(ui.NewCheckbox("Checkbox"), false)
	// vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	// 请求参数部分
	group := ui.NewGroup("请求参数")
	group.SetMargined(true)
	vbox.Append(group, true)

	group.SetChild(ui.NewNonWrappingMultilineEntry())
	// 表单
	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
	// 定义输入框
	httpMethod := ui.NewEntry()
	httpRequestUrl := ui.NewEntry()
	httpHeader := ui.NewMultilineEntry()
	httpBody := ui.NewMultilineEntry()
	httpResponse := ui.NewNonWrappingMultilineEntry()

	entryForm.Append("请求方式", httpMethod, false)
	entryForm.Append("接口地址", httpRequestUrl, false)
	// entryForm.Append("Password Entry", ui.NewPasswordEntry(), false)
	// entryForm.Append("Search Entry", ui.NewSearchEntry(), false)
	entryForm.Append("请求头", httpHeader, true)
	entryForm.Append("请求体", httpBody, true)
	httpResponse.SetReadOnly(true)
	entryForm.Append("返回值", httpResponse, true)

	// 添加按钮事件
	requestButton.OnClicked(func() {
		method := httpMethod.Text()
		if len(method) == 0 {
			method = "GET"
		}
		requestUrl := httpRequestUrl.Text()
		if len(requestUrl) == 0 {
			ui.MsgBoxError("error", "错误提示", "请求链接必须填写")
		}
		requestBody := httpBody.Text()
		requestHeader := httpHeader.Text()
		result := SendRequest(method, requestUrl, requestHeader, requestBody)
		httpResponse.SetText(result)
	})

	return vbox
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
	tab.Append("HTTP请求", makeBasicControlsPage())
	tab.SetMargined(0, true)
	// tab.Append("第二页", newBox())
	// tab.SetMargined(1, true)
	// 设置tab页
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	// 最后显示
	mainwin.Show()
}
