package core

import (
	"context"
	"demo2/utils"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ListenTarget(ctx context.Context) {

	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				fmt.Println("监听响应数据：", ev.Response)
			}
		case *network.EventRequestWillBeSent:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				fmt.Println("监听请求数据：", ev.Request)
				//数据库操作
			}
		case *page.EventWindowOpen:
			fmt.Println("page.EventWindowOpen")

		case *page.EventJavascriptDialogOpening:
			fmt.Println("page.EventJavascriptDialogOpening")
		}
	})
}

var opts = append(chromedp.DefaultExecAllocatorOptions[:],
	chromedp.Flag("headless", false),
	chromedp.Flag("disable-gpu", true),          //禁用gpu
	chromedp.Flag("disable-web-security", true), //禁用网络安全
	chromedp.Flag("disable-xss-auditor", true),  //禁用xss审核
	chromedp.Flag("no-sandbox", true),
	chromedp.Flag("disable-setuid-sandbox", true),
	chromedp.Flag("allow-running-insecure-content", true),
	chromedp.Flag("disable-webgl", true),
	chromedp.Flag("disable-popup-blocking", true),
	chromedp.Flag("block-new-web-contents", true),
	chromedp.Flag("blink-settings", "imagesEnabled=false"))

func NewChrome(url, auth string) {
	ctx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ = chromedp.NewContext(ctx)

	_ = utils.RunResponse(ctx, chromedp.Tasks{
		LoadHeaders(auth),
		chromedp.Navigate(url),
	})
}
