package main

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

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
	// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Add("Sec-Fetch-Mode", "navigate")
	request.Header.Add("Sec-Fetch-Site", "same-origin")
	request.Header.Add("Sec-Fetch-User", "?1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Add("Connection", "keep-alive")

	// Fetch Request
	httpResp, err := client.Do(request)
	if err != nil {
		fmt.Println("Failure : ", err)
		return ""
	}
	defer httpResp.Body.Close()
	data, _ := ioutil.ReadAll(czResp.Body)
	return string(data)
}
