package core

import (
	"context"
	"crypto/md5"
	"demo1/utils"
	"fmt"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//计算客户端response
func ResponseQop(hashMap map[string]string, digest Digest) (res string) {

	var A1, A2 string
	if len(digest.Qop) != 0 && digest.Qop[1:len(digest.Qop)-1] == "auth" {
		A1 = digest.UserName + ":" + hashMap["realm"][1:len(hashMap["realm"])-1] + ":" + digest.Password
		A2 = digest.Method + ":" + digest.Uri
	}

	res = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%x", md5.Sum([]byte(A1)))+":"+hashMap["nonce"][1:len(hashMap["nonce"])-1]+":"+digest.Nc+":"+digest.Cnonce+":"+digest.Qop[1:len(digest.Qop)-1]+":"+fmt.Sprintf("%x", md5.Sum([]byte(A2))))))
	return
}

func ListenTarget(ctx context.Context) {

	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventResponseReceived:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {

				if ev.Response.Status == 401 {
					auth := ParseAuthorization(ev.Response.Headers)
					AgainChrome(auth, "http://127.0.0.1:8090/index")
				}

			}
		case *network.EventRequestWillBeSent:
			if ev.Type == network.ResourceTypeDocument || ev.Type == network.ResourceTypeXHR {
				digest.Method = ev.Request.Method
				digest.Uri = utils.ParseURI(ev.Request.URL)

				//数据库操作
			}
		case *page.EventWindowOpen:
			fmt.Println("page.EventWindowOpen")

		case *page.EventJavascriptDialogOpening:
			fmt.Println("page.EventJavascriptDialogOpening")
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

func NewCrawler() {
	NewChrome("", "http://127.0.0.1:8090/index")
}

func NewChrome(headStr string, url string) {
	opts := append(Opts, chromedp.Flag("headless", true))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ListenTarget(ctx)

	_ = utils.RunResponse(ctx, chromedp.Tasks{
		LoadHeaders(headStr),
		chromedp.Navigate(url),
		page.Close(),
	})
}

func AgainChrome(headStr string, url string) {

	opts := append(Opts, chromedp.Flag("headless", false))
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ListenTarget(ctx)

	_ = utils.RunResponse(ctx, chromedp.Tasks{
		LoadHeaders(headStr),
		chromedp.Navigate(url),
	})
}
