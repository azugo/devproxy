package main

import (
	"os"

	"azugo.io/devproxy"

	"azugo.io/azugo"
	"go.uber.org/zap"
)

func main() {
	os.Setenv("ENVIRONMENT", string(azugo.EnvironmentDevelopment))

	a, err := devproxy.NewApp()
	if err != nil {
		panic(err)
	}

	if err := a.LoadConfig(); err != nil {
		a.Log().With(zap.Error(err)).Fatal("failed to load config")
	}

	azugo.Run(a)
}
