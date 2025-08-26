package main

import (
	"chera_khube/handler/app"
	"chera_khube/internal/constant"
	"context"
	"flag"
)

func main() {
	path := flag.String("e", constant.DefaultEnvPath, "env file path")
	flag.Parse()
	config, err := app.SetupViper(*path)
	if err != nil {
		panic(err)
	}

	application := app.NewApplication(config)
	ctx := context.Background()

	err = application.Setup(ctx)
	if err != nil {
		panic(err)
	}
}
