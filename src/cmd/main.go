package main

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase"
	"github.com/irdaislakhuafa/primeskills-test/src/connection"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/handler/rest"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
)

const (
	configFileJSON = "etc/cfg/conf.json"
)

func main() {
	// read config
	cfg, err := config.ReadFileJSON(configFileJSON)
	if err != nil {
		panic(err)
	}

	// initialize log
	l := log.Init(log.Config(cfg.Log))

	// initialize db
	db, err := connection.InitMySQL(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// initialize validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// initialize queries
	q := entity.New(db)

	// initialize domain
	d := domain.Init(l, q, db)

	// initialize usecase
	u := usecase.Init(d, l, v)

	// initialize api server
	r := rest.Init(cfg, l, u)
	l.Info(context.Background(), strformat.TWE("Listening at port {{ .Port }}", cfg.Gin))
	r.Run()
}
