package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

)

type MySpider struct {
	indexUrl string
}

//实现GET请求函数
func (this MySpider) readUrlBody() (string, error) {
	resp, err := http.Get(this.indexUrl)
	if err != nil {
		return "err", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "err", err
	}
	return string(body), err
}

//抓取登录的域名地址
func (this MySpider) catchCategoryUrl() []string {
	body, _ := this.readUrlBody()
	rcg := regexp.MustCompile(`class="catalog_sec_btn clearfix mgt15">(?sU:.*)<a tts="link_31301" href="(https://www.changtu.com/login/login.html)" rel="nofollow" class="lgi left"`)
	urls := rcg.FindAllStringSubmatch(body, -1)

	cateUrl := make([]string, len(urls))
	for i, u := range urls {
		cateUrl[i] = u[1]
	}
	return cateUrl
}

//通过GET请求，尝试登录请求，并且输出登录是否成功
func (this MySpider) getData() string {
	this.indexUrl = "https://www.changtu.com/login/ttslogin.htm?callback=0&username=17621396282&password=mys123456"
	body, _ := this.readUrlBody()
	rcg := regexp.MustCompile(`{"failCount":"0","flag":"true"}`)
	result := rcg.FindAllStringSubmatch(body, -1)
	if result != nil {
		this.indexUrl = "https://www.changtu.com"
		body1, _ := this.readUrlBody()
		//		fmt.Println(string(body1))
		rcg1 := regexp.MustCompile(`<p class="fontW font14 catalog_sec_hi mgt5 catalog_sec_login clearfix afterLogin hide">(?sU:.*)<span class="left">(.*?)</span>(?sU:.*)</p>`)
		result1 := rcg1.FindAllStringSubmatch(body1, -1)
		for i := range result1 {
			line1 := result1[i]
			fmt.Println("<<======>>", line1[1])
		}
		//		fmt.Println(result1)
		return ""
	}
	//	fmt.Println(string(body))
	fmt.Println("<<======>>登录失败")
	return ""
}

//验证目前是否处于未登录状态
func (this MySpider) catchLoginInfo() string {
	body, _ := this.readUrlBody()
	rcg := regexp.MustCompile(`<div class="nloginBox">(?sU:.*)<h3 class="fontW font16">(.*?)</h3>(?sU:.*)</div>`)
	result := rcg.FindAllStringSubmatch(body, -1)
	//	fmt.Println("body=", body)
	for i := range result {
		line := result[i]
		fmt.Println("<<======>>", line[1])
	}
	return ""
}

func (this MySpider) run() string {
	cateUrls := this.catchCategoryUrl()
	for _, u := range cateUrls {
		this.indexUrl = u
		this.catchLoginInfo()
		this.getData()
		break
	}
	return ""
}
func main() {
	ms := new(MySpider)
	ms.indexUrl = "https://www.changtu.com"
	ms.run()
}