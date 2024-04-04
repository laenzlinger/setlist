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
	"path/filepath"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/laenzlinger/setlist/internal/gig"
	convert "github.com/laenzlinger/setlist/internal/html/pdf"
	tmpl "github.com/laenzlinger/setlist/internal/html/template"
	"github.com/laenzlinger/setlist/internal/repertoire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		For(gig).
		Render()

	data := tmpl.Data{
		Title:   gig.Name,
		Margin:  "0cm",
		Content: template.HTML(content), //nolint: gosec // not a web application
	}

	filename, err := tmpl.CreateSetlist(&data)
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, filepath.Join(config.Target(), fmt.Sprintf("Set List %s.pdf", gig.Name)))
}
