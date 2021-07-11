package main

import (
	"context"
	"demo1/core"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var num int

func ListenTarget(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				auth := ev.Response.Headers["Www-Authenticate"].(string)
				fmt.Println("测试请求头数据", auth)

				auths := strings.SplitN(auth, " ", 2)

				// fmt.Println(auths[1])
				//如果认证方式为Digest
				if auths[0] == "Digest" {

					str := strings.Split(auths[1], ",")
					hashMap := make(map[string]string)
					for _, v := range str {
						split := strings.Split(v, ",")

						for _, s := range split {
							i := strings.Split(s, "=")
							hashMap[i[0]] = i[1]
						}
					}

					for k, v := range hashMap {
						fmt.Println(k, ":", v)
					}
				}
			}
		case *network.EventRequestWillBeSent:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				req := ev.Request
				reqUrl := req.URL
				fmt.Println("监听是否发送请求：", reqUrl)
				//数据库操作
			}
		case *page.EventWindowOpen:
			fmt.Println("窗口已打开")
			// default:
			// 	fmt.Println("hello world!")
		case *page.EventJavascriptDialogOpening:
			fmt.Println("不明白的代码")
			// default:
			// 	num++
			// 	fmt.Println(num)
		}
	})
}

func login(str *url.URL) chromedp.ActionFunc {

	var weakPasswords []Login
	weakPasswords = append(weakPasswords,
		Login{
			UserName: "mengjiaheng",
			Password: "123456",
		},
		Login{
			UserName: "asants",
			Password: "asants",
		})
	return func(ctx context.Context) error {

		if strings.Contains(str.Host, "login") {
			//等待页面加载完
			_ = chromedp.Sleep(2 * time.Second).Do(ctx)
			var htmlArray []core.Html
			chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(core.JSCode, 0), &htmlArray))
		Loop:
			for _, weakPassword := range weakPasswords {
				for _, re := range htmlArray {
					if re.ElType == "text" {
						_ = chromedp.SendKeys(re.Xpath, weakPassword.UserName).Do(ctx)
					}

					if re.ElType == "password" {
						_ = chromedp.SendKeys(re.Xpath, weakPassword.Password).Do(ctx)
					}

					if re.ElType == "submit" {
						_ = chromedp.Submit(re.Xpath).Do(ctx)

						//等待登录后跳转页面的时间
						_ = chromedp.Sleep(2 * time.Second).Do(ctx)

						var res string
						_ = chromedp.Evaluate("window.location.href", &res).Do(ctx)

						fmt.Println("res:", res)
						if str.String() == res {
							continue
						} else {
							break Loop
						}
					}
				}
			}
		}
		return nil
	}
}

//设置请求头
func loadHeaders(auth string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return network.SetExtraHTTPHeaders(
			network.Headers{"Authorization": auth}).Do(ctx)

	}
}

var opts = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.Flag("headless", false),
	chromedp.Flag("disable-gpu", true),
	chromedp.Flag("disable-web-security", true),
	chromedp.Flag("disable-xss-auditor", true),
	chromedp.Flag("no-sandbox", true),
	chromedp.Flag("disable-setuid-sandbox", true),
	chromedp.Flag("allow-running-insecure-content", true),
	chromedp.Flag("disable-webgl", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("block-new-web-contents", true),
	chromedp.Flag("blink-settings", "imagesEnabled=false"))

func main() {

	// 禁用chrome headless
	// opts = append(chromedp.DefaultExecAllocatorOptions[:],
	// 	chromedp.Flag("headless", false),
	// )

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// ctx, cancel = context.WithTimeout(ctx, time.Second*2)
	ctx, cancel = chromedp.NewContext(ctx)

	//待定
	defer cancel()

	// var site *url.URL
	ListenTarget(ctx)
	err := chromedp.Run(ctx, chromedp.Tasks{
		loadHeaders(""),
		chromedp.Navigate("http://127.0.0.1:8090/index"),
		chromedp.Sleep(2 * time.Second),
		// login(site),
	})
	// go core.Form(ctx)
	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	// navigate to a page, wait for an element, click
	// var example string
	// sel := `//*[@id="username"]`
	// err := chromedp.Run(ctx,
	// 	loadHeaders(),
	// 	chromedp.Navigate(`http://127.0.0.1:8090/index`),
	// 	chromedp.WaitVisible("body"),
	// 	//缓一缓
	// 	chromedp.Sleep(2*time.Second),

	// 	chromedp.SendKeys(sel, "username", chromedp.BySearch), //匹配xpath

	// )
	// time.Sleep(2 * time.Second)
	if err != nil {
		log.Fatal("未知报错：", err)
	}

	// log.Printf("Go's time.After example:\n%s", example)

}
