package config

import "os"

func RunningInContainer() bool {
	return os.Getenv("OS_ENV") == "container"
}

func UserHome() string {
	return os.Getenv("HOME")
}
