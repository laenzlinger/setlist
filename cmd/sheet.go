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
		band := cmd.Flag("band").Value.String()
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}
		if all {
			err = sheet.AllForBand(band)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			gigName := cmd.Flag("gig").Value.String()
			gig, err := gig.New(band, gigName)
			if err != nil {
				log.Fatal(err)
			}
			err = sheet.ForGig(band, gig)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	rootCmd.AddCommand(sheetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cheatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	sheetCmd.Flags().BoolP("all", "a", false, "Generate a cheat sheet out of all songs.")
}
