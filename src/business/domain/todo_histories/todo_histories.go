package todo_histories

import (
	"context"
	"database/sql"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
)

type (
	Interface interface {
		List(ctx context.Context, params entity.ListTodoHistoriesParams) ([]entity.TodoHistory, error)
	}

	todoHistory struct {
		log     log.Interface
		db      *sql.DB
		queries *entity.Queries
	}
)

func Init(log log.Interface, db *sql.DB, queries *entity.Queries) Interface {
	return &todoHistory{
		log:     log,
		db:      db,
		queries: queries,
	}
}

func (th *todoHistory) List(ctx context.Context, params entity.ListTodoHistoriesParams) ([]entity.TodoHistory, error) {
	rows, err := th.queries.ListTodoHistories(ctx, params)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	results := []entity.TodoHistory{}
	for _, row := range rows {
		results = append(results, entity.TodoHistory{
			ID:        row.ID,
			TodoID:    row.TodoID,
			Message:   row.Message,
			CreatedAt: row.CreatedAt,
			CreatedBy: row.CreatedBy,
			UpdatedAt: row.UpdatedAt,
			UpdatedBy: row.UpdatedBy,
			DeletedAt: row.DeletedAt,
			DeletedBy: row.DeletedBy,
			IsDeleted: row.IsDeleted,
		})
	}
	return results, nil
}
