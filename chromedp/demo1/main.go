package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ListenTarget(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				str := ev.Response.Headers["Www-Authenticate"].(string)
				fmt.Println("测试请求头数据", str)
			}
		case *network.EventRequestWillBeSent:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				req := ev.Request
				reqUrl := req.URL
				fmt.Println("监听是否发送请求：", reqUrl)
			}
		case *page.EventWindowOpen:

			// default:
			// 	fmt.Println("hello world!")
		}
	})
}

func main() {

	// 禁用chrome headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ListenTarget(ctx)
	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	// navigate to a page, wait for an element, click
	var example string
	sel := `//*[@id="username"]`
	err := chromedp.Run(ctx,
		// chromedp.Navigate(`https://github.com/awake1t`),
		chromedp.Navigate(`http://127.0.0.1:8090/index`),
		chromedp.WaitVisible("body"),
		//缓一缓
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys(sel, "username", chromedp.BySearch), //匹配xpath

	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Go's time.After example:\n%s", example)

}
