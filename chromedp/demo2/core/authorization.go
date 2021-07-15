package core

import (
	"context"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

//设置请求头
func LoadHeaders(auth string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		return network.SetExtraHTTPHeaders(
			network.Headers{"Authorization": auth}).Do(ctx)

	}
}
