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
	"errors"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const AppName = "setlist"

//nolint:gochecknoglobals // cobra is designed like this
var cfgFile string

//nolint:gochecknoglobals // cobra is designed like this
var rootCmd = &cobra.Command{
	Use:   AppName,
	Short: "CLI to maintain a repertoire for artists and bands.",
	Long: `Generate Cheat Sheet or Setlist out of repertoire based on Markdown and PDF files.
`,
	DisableAutoGenTag: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//nolint:gochecknoinits // cobra is desigend like this
func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .setlist.yaml)")

	rootCmd.PersistentFlags().StringP("band-name", "b", "", "the name of the band")
	err := viper.BindPFlag("band.name", rootCmd.PersistentFlags().Lookup("band-name"))
	cobra.CheckErr(err)

	rootCmd.PersistentFlags().StringP("band-src", "s", "", "the source directory (default is band-name)")
	err = viper.BindPFlag("band.src", rootCmd.PersistentFlags().Lookup("band-src"))
	cobra.CheckErr(err)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	if len(cfgFile) > 0 {
		viper.SetConfigFile(cfgFile)
	} else if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		viper.AddConfigPath(path.Join(xdg, AppName))
		viper.SetConfigName("config")
	} else {
		viper.AddConfigPath(path.Join(home, ".config", AppName))
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	var cnf = viper.ConfigFileNotFoundError{}
	if errors.As(err, &cnf) {
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName("." + AppName)
		err = viper.ReadInConfig()
		if errors.As(err, &cnf) {
			err = nil // It's ok if the config file does not exist
		}
	}
	cobra.CheckErr(err)
}

func Instance() *cobra.Command {
	return rootCmd
}
