package domain

import (
	"database/sql"

	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain/todo"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain/user"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
)

type (
	Domain struct {
		User user.Interface
		Todo todo.Interface
	}
)

func Init(log log.Interface, queries *entity.Queries, db *sql.DB) *Domain {
	return &Domain{
		User: user.Init(log, queries, db),
		Todo: todo.Init(log, queries, db),
	}
}
