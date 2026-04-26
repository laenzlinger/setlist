/*
Copyright © 2024 Christof Laenzlinger <christof@laenzlinger.net>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/laenzlinger/setlist/internal/gig"
	convert "github.com/laenzlinger/setlist/internal/html/pdf"
	tmpl "github.com/laenzlinger/setlist/internal/html/template"
	"github.com/laenzlinger/setlist/internal/repertoire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultPages = 2

//nolint:gochecknoglobals // cobra is designed like this
var setlistCmd = &cobra.Command{
	Use:   "list",
	Short: "Generate a set list",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Long: `Generates a setlist for a gig.
`,
	Run: func(_ *cobra.Command, args []string) {
		gig, err := config.GigName(args)
		cobra.CheckErr(err)
		err = generateSetlist(gig)
		cobra.CheckErr(err)
	},
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	generateCmd.AddCommand(setlistCmd)

	setlistCmd.Flags().StringSliceP("include-columns", "i", []string{"Title", "Year", "Description"},
		"defines the repertoire columns to include in the output")

	err := viper.BindPFlag("setlist.include-columns", setlistCmd.Flags().Lookup("include-columns"))
	cobra.CheckErr(err)

	setlistCmd.Flags().StringP("font-size", "f", "", "set the main font size (css values are supported)")
	err = viper.BindPFlag("generate.list.font-size", setlistCmd.Flags().Lookup("font-size"))
	cobra.CheckErr(err)

	setlistCmd.Flags().IntP("pages", "p", defaultPages, "target number of pages (auto-fits font size)")
	err = viper.BindPFlag("generate.list.pages", setlistCmd.Flags().Lookup("pages"))
	cobra.CheckErr(err)
}

func generateSetlist(gigName string) error {
	include := viper.GetStringSlice("setlist.include-columns")
	band := config.NewBand()

	rep, err := repertoire.New(band)
	if err != nil {
		return err
	}

	gig, err := gig.New(band, gigName)
	if err != nil {
		return err
	}

	content := rep.NoHeader().
		IncludeColumns(include...).
		Merge(gig).
		Render()

	outPath := filepath.Join(config.Target(), fmt.Sprintf("Set List %s.pdf", gig.Name))

	fontSize := viper.GetString("generate.list.font-size")
	if fontSize != "" {
		return renderSetlist(content, gig.Name, fontSize, outPath)
	}

	targetPages := viper.GetInt("generate.list.pages")
	return autoFitSetlist(content, gig.Name, targetPages, outPath)
}

func renderSetlist(content, title, fontSize, outPath string) error {
	filename, err := createSetlistHTML(content, title, fontSize)
	if err != nil {
		return err
	}
	defer os.Remove(filename)

	return convert.HTMLToPDF(filename, outPath)
}

func autoFitSetlist(content, title string, targetPages int, outPath string) error {
	const (
		minFontSize = 10.0
		maxFontSize = 60.0
		precision   = 0.5
	)

	lo, hi := minFontSize, maxFontSize

	for hi-lo > precision {
		mid := (lo + hi) / 2 //nolint:mnd // binary search midpoint
		fontSize := fmt.Sprintf("%.1fpx", mid)

		pages, err := countPages(content, title, fontSize)
		if err != nil {
			return err
		}

		if pages <= targetPages {
			lo = mid
		} else {
			hi = mid
		}
	}

	fontSize := fmt.Sprintf("%.1fpx", lo)
	log.Printf("auto-fit font size: %s for %d pages", fontSize, targetPages)

	return renderSetlist(content, title, fontSize, outPath)
}

func countPages(content, title, fontSize string) (int, error) {
	filename, err := createSetlistHTML(content, title, fontSize)
	if err != nil {
		return 0, err
	}
	defer os.Remove(filename)

	pdfBytes, err := convert.HTMLToPDFBytes(filename)
	if err != nil {
		return 0, err
	}

	return convert.PageCount(pdfBytes)
}

func createSetlistHTML(content, title, fontSize string) (string, error) {
	data := tmpl.Data{
		Title:    title,
		FontSize: fontSize,
		Margin:   "0cm",
		Content:  template.HTML(content), //nolint: gosec // not a web application
	}
	return tmpl.CreateSetlist(&data)
}
