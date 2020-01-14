package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	// . "github.com/libra/rtools"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

//
func main() {
	ui.Main(setUp)
}

//
func makeBasicControlsPage(w *ui.Window) ui.Control {
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
	entryForm.Append("请求头", httpHeader, true)
	entryForm.Append("请求体", httpBody, true)
	httpResponse.SetReadOnly(true)
	entryForm.Append("返回值", httpResponse, true)

	// 添加按钮事件
	requestButton.OnClicked(func(*ui.Button) {
		method := httpMethod.Text()
		if len(method) == 0 {
			method = "GET"
		}
		requestUrl := httpRequestUrl.Text()
		if len(requestUrl) == 0 {
			ui.MsgBoxError(w, "错误提示", "请求链接必须填写")
			return
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
	mainwin := ui.NewWindow("请求工具", 640, 480, false)
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
	tab.Append("HTTP请求", makeBasicControlsPage(mainwin))
	tab.SetMargined(0, true)
	// tab.Append("第二页", newBox())
	// tab.SetMargined(1, true)
	// 设置tab页
	mainwin.SetChild(tab)
	mainwin.SetMargined(true)

	// 最后显示
	mainwin.Show()
}

// 发送请求
func SendRequest(method, requestUrl, requestHeader, requestBody string) string {
	//Init jar
	j, _ := cookiejar.New(nil)
	// Create client
	client := &http.Client{Jar: j}
	//建立http请求对象
	request, _ := http.NewRequest(strings.ToUpper(method), requestUrl, strings.NewReader(requestBody))
	//这个一定要加，不加form的值post不过去
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Add("Sec-Fetch-Mode", "navigate")
	request.Header.Add("Sec-Fetch-Site", "same-origin")
	request.Header.Add("Sec-Fetch-User", "?1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Add("Connection", "keep-alive")
	if len(requestHeader) > 0 {
		headers := strings.Split(requestHeader, "/n")
		for i := 0; i < len(headers); i++ {
			headerList := strings.Split(headers[i], ":")
			if len(headerList) > 1 {
				request.Header.Set(headerList[0], headerList[1])
			}
		}
	}
	// Fetch Request
	httpResp, err := client.Do(request)
	if err != nil {
		fmt.Println("Failure : ", err)
		return ""
	}
	defer httpResp.Body.Close()
	data, _ := ioutil.ReadAll(httpResp.Body)
	return string(data)
}
