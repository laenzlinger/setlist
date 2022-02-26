package convert

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func HtmlToPdf(in, out string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("file://%s/%s", wd, in)

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(url, &buf)); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(out, buf, 0o644); err != nil {
		log.Fatal(err)
	}

}

func printToPDF(urlstr string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(false).
				WithPaperHeight(11.7).
				WithPaperWidth(8.3).
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
