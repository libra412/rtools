package services

import (
	"encoding/json"
	"time"

	//"errors"
	// "flag"
	"fmt"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

//
type NetEase struct {
	OrderNo      string
	UserAccount  string
	BuyAmount    string
	RechargeMode string
}

var myClient *NetEaseLogin

//
type NetEaseLogin struct {
	MyClient  *http.Client
	LoginTime int64
}

const (
	user = "aini63750081"
	pwd  = "wu63750081"
)

//
func Recharge(userAccount, buyAmount, rechargeType string) (string, string) {
	// 登陆成功返回请求客户端
	//client, result := login()
	getMyClient()
	if myClient.MyClient != nil {
		ok, err := checkUserAccount(userAccount, myClient.MyClient)
		if err != nil {
			fmt.Println("检测账号===", err)
			myClient = nil
			return "500", "账号密码登陆失败"
		}
		if ok {
			csrf := getCsrfValue(refreshContent(myClient.MyClient, user))
			if csrf == "" {
				return "500", "登陆失败"
			}
			fmt.Println(csrf)
			res := rechargeUserAccount(userAccount, buyAmount, rechargeType, csrf, myClient.MyClient)
			if strings.Index(res, "成功") > -1 {
				return "200", "充值成功"
			}
			fmt.Println("充值结果======", res)
			return "500", "充值失败"
		}
		return "500", "账号错误"
	}
	return "500", "账号密码登陆失败"
}

func getMyClient() {
	if myClient == nil {
		client := login()
		myClient = &NetEaseLogin{client, time.Now().Unix()}
	} else {
		diff := time.Now().Unix() - myClient.LoginTime
		if diff/60 > 10 { // 超过 10分钟 重新登录
			fmt.Println(diff/60, "重新登陆")
			client := login()
			myClient.MyClient = client
			myClient.LoginTime = time.Now().Unix()
		}
	}
}

// 登陆
func login() *http.Client {
	//Init jar
	j, _ := cookiejar.New(nil)
	// Create client
	client := &http.Client{Jar: j}
	// Create request
	v := url.Values{}
	v.Set("seller_name", user)
	v.Set("password", pwd)
	v.Set("oper_name", "")
	v.Set("a", "seller_login")
	v.Set("otp", "")
	body := strings.NewReader(v.Encode()) //把form数据编下码
	//建立http请求对象
	request, _ := http.NewRequest("POST", "https://esales.163.com/script/seller/login", body)
	//这个一定要加，不加form的值post不过去
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	request.Header.Add("Host", "esales.163.com")
	request.Header.Add("Origin", "https://esales.163.com")
	request.Header.Add("Referer", "https://esales.163.com/")
	request.Header.Add("Sec-Fetch-Mode", "navigate")
	request.Header.Add("Sec-Fetch-Site", "same-origin")
	request.Header.Add("Sec-Fetch-User", "?1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Add("Connection", "keep-alive")
	// Fetch Request
	httpResp, err := client.Do(request)
	if err != nil {
		fmt.Println("Failure : ", err)
		return nil
	}
	defer httpResp.Body.Close()
	//全局保存
	fmt.Println("jar======", j.Cookies(request.URL))
	return client
}

//
func refreshContent(client *http.Client, user string) string {
	chongzhi, _ := http.NewRequest("GET", "https://esales.163.com/script/seller/direct_credit/?a=show", nil)
	czResp, err := client.Do(chongzhi)
	if err != nil {
		fmt.Println("Failure2 : ", err)
	}
	defer czResp.Body.Close() //关闭resp.Body
	data2, _ := ioutil.ReadAll(czResp.Body)
	decode := mahonia.NewDecoder("gbk")
	result := decode.ConvertString(string(data2))
	idx := strings.Index(result, user)
	if idx != -1 {
		return result
	}
	return ""
}

// 充值
func rechargeUserAccount(userAccount, buyAmount, rechargeType, csrf string, client *http.Client) string {
	v := url.Values{}
	v.Set("_csrf", csrf)
	v.Set("a", "trans")        //trans
	v.Set("pts_shortcut", "0") //0其他
	v.Set("pts", buyAmount)    //其他购买数
	v.Set("urs", userAccount)
	v.Set("urs_repeat", userAccount)
	v.Set("reason", rechargeType)         //1是直冲，2是寄售
	body := strings.NewReader(v.Encode()) //把form数据编下码
	fmt.Println("request body", body)
	request, _ := http.NewRequest("POST", "https://esales.163.com/script/seller/direct_credit/", body)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.97 Safari/537.36")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3")
	// Fetch Request
	httpResp, err := client.Do(request)
	if err != nil {
		fmt.Printf("httpResp.Header=%s", httpResp.Header)
		fmt.Println("rechargeUserAccount===Failure : ", err)
		return ""
	}
	//
	defer httpResp.Body.Close() //关闭resp.Body
	data, _ := ioutil.ReadAll(httpResp.Body)
	decode := mahonia.NewDecoder("gbk")
	result := decode.ConvertString(string(data))
	return result
}

//获取 csrf的值
func getCsrfValue(content string) string {
	reg := regexp.MustCompile(`<input(\s)type="hidden"(\s)name="_csrf"(\s)value="(.+)"(\s*?)(/?)>`)
	list := reg.FindStringSubmatch(content)
	// fmt.Println(content, list)
	if len(list) > 4 {
		return list[4]
	}
	return ""
}

// 检查账号是否有效
func checkUserAccount(userAccount string, client *http.Client) (bool, error) {
	checkUrl, _ := http.NewRequest("GET", "https://esales.163.com/script/seller/direct_credit?a=check_urs&urs="+userAccount, nil)
	resp, err := client.Do(checkUrl)
	if err != nil {
		return false, err
	}
	var res struct {
		UrsStatus bool `json:"urs_status"`
		IsMainUrs bool `json:"is_main_urs"`
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return false, err
	}
	return res.UrsStatus, nil
}
