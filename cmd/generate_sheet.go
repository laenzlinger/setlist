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
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	Short: "Generate a cheat sheet",
	Long: `Generates a cheat sheet for a Gig or for all songs.

Currently supports pdf sheets.
The pdf sheets are optionally generated for odf files.
`,
	Run: func(cmd *cobra.Command, args []string) {
		band := viper.GetString("band.name")
		all, err := cmd.Flags().GetBool("all")
		cobra.CheckErr(err)
		if all {
			err = sheet.AllForBand(band)
		} else {
			if len(args) == 0 {
				log.Fatal("gig name not provided")
			}
			gig, e := gig.New(band, args[0])
			cobra.CheckErr(e)
			err = sheet.ForGig(band, gig)
		}
		cobra.CheckErr(err)
	},
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	generateCmd.AddCommand(sheetCmd)

	sheetCmd.Flags().BoolP("all", "a", false, "Generate a cheat sheet out of all songs (ignores --gig).")
}
