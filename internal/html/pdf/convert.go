package convert

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/laenzlinger/setlist/internal/config"
)

func HTMLToPDF(in, out string) error {
	buf, err := HTMLToPDFBytes(in)
	if err != nil {
		return err
	}
	if err = os.WriteFile(out, buf, 0o600); err != nil {
		return fmt.Errorf("failed to write PDF: %w", err)
	}
	return nil
}

func HTMLToPDFBytes(in string) ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("file://%s/%s", wd, in)
	opts := []chromedp.ExecAllocatorOption{}
	if config.RunningInContainer() {
		opts = append(opts, chromedp.NoSandbox)
	}

	cctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(cctx)
	defer cancel()

	var buf []byte
	if err = chromedp.Run(ctx, printToPDF(url, &buf)); err != nil {
		return nil, fmt.Errorf("failed to print PDF: %w", err)
	}
	return buf, nil
}

func PageCount(pdfBytes []byte) (int, error) {
	// Count page objects by scanning for the /Type /Page pattern.
	// This is a lightweight approach that avoids full PDF parsing,
	// which can fail on Chrome-generated PDFs.
	count := 0
	needle := []byte("/Type /Page")
	notLeaf := []byte("/Type /Pages")
	for i := 0; i < len(pdfBytes); {
		idx := bytes.Index(pdfBytes[i:], needle)
		if idx < 0 {
			break
		}
		pos := i + idx
		// Skip /Type /Pages (non-leaf page tree nodes).
		if pos+len(notLeaf) <= len(pdfBytes) && bytes.Equal(pdfBytes[pos:pos+len(notLeaf)], notLeaf) {
			i = pos + len(notLeaf)
			continue
		}
		count++
		i = pos + len(needle)
	}
	return count, nil
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
				WithPrintBackground(true).
				WithPaperHeight(A4Height).
				WithPaperWidth(A4Width).
				WithMarginTop(0.2).
				WithMarginBottom(0.2).
				WithMarginLeft(0.2).
				WithMarginRight(0.2).
				WithLandscape(config.Landscape()).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
