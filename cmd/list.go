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
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Generate a Setlist",
	Long: `Generates a Setlist for a Gig.
`,
	Run: func(cmd *cobra.Command, args []string) {
		generateSetlist(cmd.Flag("band").Value.String(), cmd.Flag("gig").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

}

func generateSetlist(band, gigName string) {
	rep := repertoire.New()
	gig := gig.New(band, gigName)

	content := rep.Filter(gig).
		RemoveColumns("Lead", "Copyright", "Key").
		Render()

	data := tmpl.TemplateData{
		Title:   gig.Name,
		Content: template.HTML(content),
	}
	filename := tmpl.CreateSetlist(&data)
	convert.HtmlToPdf(filename, fmt.Sprintf("out/Setlist %s.pdf", gig.Name))
}
