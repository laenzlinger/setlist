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
	"log"

	"github.com/laenzlinger/setlist/internal/gig"
	"github.com/laenzlinger/setlist/internal/sheet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals // cobra is designed like this
var sheetCmd = &cobra.Command{
	Use:   "sheet",
	Short: "Generate a Cheat Sheet",
	Long: `Generates a Cheat Sheet for a Gig.

Currently supports pdf sheets.
The sheets can optionally be generated for odf files.
`,
	Run: func(cmd *cobra.Command, _ []string) {
		band := viper.GetString("band.name")
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}
		if all {
			err = sheet.AllForBand(band)
		} else {
			gigName := viper.GetString("gig.name")
			gig, e := gig.New(band, gigName)
			if e != nil {
				log.Fatal(e)
			}
			err = sheet.ForGig(band, gig)
		}
		if err != nil {
			log.Fatal(err)
		}

	},
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	rootCmd.AddCommand(sheetCmd)

	sheetCmd.Flags().BoolP("all", "a", false, "Generate a cheat sheet out of all songs (ignores --gig).")
}
