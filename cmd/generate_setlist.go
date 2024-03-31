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
	Short: "Generate a Setlist",
	Long: `Generates a Setlist for a Gig.
`,
	Run: func(_ *cobra.Command, _ []string) {
		err := generateSetlist()
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

func generateSetlist() error {
	band := viper.GetString("band.name")
	gigName := viper.GetString("gig.name")
	include := viper.GetStringSlice("setlist.include-columns")

	rep, err := repertoire.New(band)
	if err != nil {
		return err
	}

	gig, err := gig.New(band, gigName)
	if err != nil {
		return err
	}

	content := rep.Filter(gig).
		IncludeColumns(include...).
		Render()

	data := tmpl.Data{
		Title:   gig.Name,
		Content: template.HTML(content), //nolint: gosec // not a web application
	}

	filename, err := tmpl.CreateSetlist(&data)
	if err != nil {
		return err
	}

	return convert.HTMLToPDF(filename, fmt.Sprintf("out/Setlist %s.pdf", gig.Name))
}