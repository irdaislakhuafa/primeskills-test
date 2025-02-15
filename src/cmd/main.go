package main

import (
	"context"

	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/handler/rest"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
)

const (
	configFileJSON = "etc/cfg/conf.json"
)

func main() {
	cfg, err := config.ReadFileJSON(configFileJSON)
	if err != nil {
		panic(err)
	}

	l := log.Init(log.Config(cfg.Log))

	r := rest.Init(cfg, l)
	l.Info(context.Background(), strformat.TWE("Listening at port {{ .Port }}", cfg.Gin))
	r.Run()
}
