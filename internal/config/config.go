package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Band struct {
	Name   string
	Source string
}

func GigName(args []string) (string, error) {
	var result string
	if len(args) == 0 {
		result = viper.GetString("gig.name")
	} else {
		result = args[0]
	}
	if len(result) == 0 {
		return result, errors.New("gig name not provided")
	}
	return result, nil
}

func RunningInContainer() bool {
	return os.Getenv("OS_ENV") == "container"
}

func UserHome() string {
	return os.Getenv("HOME")
}

func Target() string {
	return viper.GetString("generate.target")
}

func Landscape() bool {
	return viper.GetBool("generate.landscape")
}

func PlaceholderDir() string {
	return filepath.Join(Target(), "placeholder")
}

// returns configuration of source band.
func NewBand() Band {
	band := Band{
		Name:   viper.GetString("band.name"),
		Source: viper.GetString("band.src"),
	}

	if band.Name == "" {
		cobra.CheckErr("Band name is mandatory")
	}

	if band.Source == "" {
		band.Source = filepath.Join(".", band.Name)
	}
	return band
}
