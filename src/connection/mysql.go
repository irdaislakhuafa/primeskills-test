package connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
)

func InitMySQL(cfg config.Config) (*sql.DB, error) {
	dsFormat := "{{ .Username }}:{{ .Password }}@tcp({{ .Host }}:{{ .Port }})/{{ .DBName }}?parseTime=true"
	dsn := strformat.TWE(dsFormat, cfg.DB.Master)
	db, err := sql.Open(DriverNameMySQL, dsn)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeSQL, "cannot connect to db: %v", err)
	}

	return db, nil
}
