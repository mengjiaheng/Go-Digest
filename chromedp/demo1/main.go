package main

import (
	"context"
	"demo1/core"
	"demo1/utils"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var num int

var digest core.Digest

// var authCh = make(chan string)

// var ctx context.Context
// var cancel context.CancelFunc

func ListenTarget(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				auth := ev.Response.Headers["Www-Authenticate"].(string)

				auths := strings.SplitN(auth, " ", 2)
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

					// for k, v := range hashMap {
					// 	fmt.Println(k, ":", v)
					// }

					digest.Qop = hashMap["qop"]

					//产生客户端随机数与请求计数器
					digest.Cnonce = utils.RandomString()
					digest.Nc = "0000001"

					//设置账号密码
					digest.UserName = "mengjiaheng"
					digest.Password = "123456"
					digest.Response = utils.ResponseQop(hashMap, digest)

					auth = `Digest username="` + digest.UserName + `",realm="` + hashMap["realm"][1:len(hashMap["realm"])-1] + `",qop=` + hashMap["qop"][1:len(hashMap["qop"])-1] + `,nonce="` + hashMap["nonce"][1:len(hashMap["nonce"])-1] + `",uri="` + digest.Uri + `",nc=` + digest.Nc + `,cnonce="` + digest.Cnonce + `",response="` + digest.Response + `"`

					opts := append(Opts, chromedp.Flag("headless", false))
					ctxT, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
					ctxT, cancel = chromedp.NewContext(ctxT)
					err := core.RunResponse(ctxT, chromedp.Tasks{
						core.LoadHeaders(auth),
						chromedp.Navigate("http://127.0.0.1:8090/index"),
					})
					if err != nil {
						cancel()
					}
				}
			}
		case *network.EventRequestWillBeSent:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				digest.Method = ev.Request.Method

				str := strings.Split(ev.Request.URL, "/")
				for i := 3; 2 < i && i < len(str); i++ {
					digest.Uri = digest.Uri + "/" + str[i]
				}
				// fmt.Println("监听是否发送请求：", digest.Method, digest.Uri)
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

var Opts = append(chromedp.DefaultExecAllocatorOptions[:],
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
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		opts := append(Opts, chromedp.Flag("headless", true))
		ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel = chromedp.NewContext(ctx)

		ListenTarget(ctx)

		err := core.RunResponse(ctx, chromedp.Tasks{
			core.LoadHeaders(""),
			chromedp.Navigate("http://127.0.0.1:8090/index"),
			page.Close(),
		})
		if err != nil {
			cancel()
		}
		c.String(200, "OK")
	})
	r.Run(":9000")
}
