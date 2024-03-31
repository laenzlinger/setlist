package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Band struct {
	Name   string
	Source string
}

func RunningInContainer() bool {
	return os.Getenv("OS_ENV") == "container"
}

func UserHome() string {
	return os.Getenv("HOME")
}

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
