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
	"os"

	"github.com/laenzlinger/setlist/internal/config"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals // cobra is designed like this
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the target directory",
	Long: `The target directory and all its contents will be deleted.
`,
	Run: func(_ *cobra.Command, _ []string) {
		os.RemoveAll(config.Target())
	},
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	rootCmd.AddCommand(cleanCmd)
}
