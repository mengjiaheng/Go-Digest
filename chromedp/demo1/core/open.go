package core

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type Digest struct {
	UserName string
	Password string
	Method   string
	Realm    string
	Qop      string
	Nonce    string
	Uri      string
	Nc       string
	Cnonce   string
	Response string
}

//设置请求头
func LoadHeaders(auth string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return network.SetExtraHTTPHeaders(
			network.Headers{"Authorization": auth}).Do(ctx)

	}
}
func RunResponse(ctx context.Context, tasks chromedp.Tasks) error {
	// for {
	_, err := chromedp.RunResponse(ctx, tasks)
	return err
	// }
}
