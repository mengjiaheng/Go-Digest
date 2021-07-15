package utils

import (
	"context"

	"github.com/chromedp/chromedp"
)

func RunResponse(ctx context.Context, tasks chromedp.Tasks) error {
	_, err := chromedp.RunResponse(ctx, tasks)
	return err
}
