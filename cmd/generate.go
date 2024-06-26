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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//nolint:gochecknoglobals // cobra is designed like this
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output",
	Long: `Generate a set list or a cheat sheet.
`,
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringP("target", "t", "out", "the target directory")
	err := viper.BindPFlag("generate.target", generateCmd.PersistentFlags().Lookup("target"))
	cobra.CheckErr(err)

	generateCmd.PersistentFlags().BoolP("landscape", "l", false, "generate landscape document")
	err = viper.BindPFlag("generate.landscape", generateCmd.PersistentFlags().Lookup("landscape"))
	cobra.CheckErr(err)
}
