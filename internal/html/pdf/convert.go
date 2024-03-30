package convert

import (
	"context"
	"fmt"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func HTMLToPDF(in, out string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("file://%s/%s", wd, in)
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoSandbox,
	}
	cctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(cctx)
	defer cancel()

	var buf []byte
	if err = chromedp.Run(ctx, printToPDF(url, &buf)); err != nil {
		return fmt.Errorf("failed to print PDF: %w", err)
	}

	if err = os.WriteFile(out, buf, 0o600); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}
	return nil
}

const (
	A4Height = 11.7
	A4Width  = 8.3
)

func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(false).
				WithPaperHeight(A4Height).
				WithPaperWidth(A4Width).
				WithLandscape(true).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
