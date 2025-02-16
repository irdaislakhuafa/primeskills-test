package todo

import (
	"context"
	"database/sql"
	"time"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/ctxkey"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/entutils"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateTodoParams) (entity.Todo, error)
		List(ctx context.Context, params entity.ListTodoParams) ([]entity.Todo, error)
	}
	todo struct {
		log     log.Interface
		queries *entity.Queries
		db      *sql.DB
	}
)

func Init(log log.Interface, queries *entity.Queries, db *sql.DB) Interface {
	return &todo{
		log:     log,
		queries: queries,
		db:      db,
	}
}

func (t *todo) Create(ctx context.Context, params entity.CreateTodoParams) (entity.Todo, error) {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	queries := t.queries.WithTx(tx)

	// ensure user is exists
	_, err = queries.GetOneUser(ctx, entity.GetOneUserParams{ID: params.UserID, IsDeleted: 0})
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Todo{}, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "User not found!")
		}
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	// create todo
	params.CreatedAt = time.Now()
	params.CreatedBy = convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID))
	params.Status = entutils.TODO_STATUS_TODO
	r, err := queries.CreateTodo(ctx, params)
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	result := entity.Todo{
		UserID:      params.UserID,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		CreatedAt:   params.CreatedAt,
		CreatedBy:   params.CreatedBy,
	}

	if result.ID, err = r.LastInsertId(); err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	return result, nil
}

func (t *todo) List(ctx context.Context, params entity.ListTodoParams) ([]entity.Todo, error) {
	params.Status = strformat.TWE("%{{ .Status }}%", params)
	rows, err := t.queries.ListTodo(ctx, params)
	if err != nil {
		return nil, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	results := []entity.Todo{}
	for _, row := range rows {
		results = append(results, entity.Todo{
			ID:          row.ID,
			UserID:      row.UserID,
			Title:       row.Title,
			Description: row.Description,
			Status:      row.Status,
			CreatedAt:   row.CreatedAt,
			CreatedBy:   row.CreatedBy,
			UpdatedAt:   row.UpdatedAt,
			UpdatedBy:   row.UpdatedBy,
			DeletedAt:   row.DeletedAt,
			DeletedBy:   row.DeletedBy,
			IsDeleted:   row.IsDeleted,
		})
	}

	return results, nil
}
