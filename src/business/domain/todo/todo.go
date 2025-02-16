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
		Update(ctx context.Context, params entity.UpdateTodoParams) (entity.Todo, error)
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

	// fill data for result
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

	// create todo histories
	_, err = queries.CreateTodoHistory(ctx, entity.CreateTodoHistoryParams{
		TodoID:    result.ID,
		Message:   strformat.TWE("Todo '{{ .Title }}' created with status {{ .Status }}", result),
		CreatedAt: result.CreatedAt,
		CreatedBy: result.CreatedBy,
	})
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	// commit changes
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

func (t *todo) Update(ctx context.Context, params entity.UpdateTodoParams) (entity.Todo, error) {
	// begin tx
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	// prepare queries
	queries := t.queries.WithTx(tx)

	// fill necessary fields
	params.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	params.UpdatedBy = sql.NullString{
		Valid:  true,
		String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}
	if params.IsDeleted == 1 {
		params.DeletedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		params.DeletedBy = sql.NullString{
			Valid:  true,
			String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
		}
	}

	// ensure data todo is exists
	prev, err := queries.GetOneTodo(ctx, entity.GetOneTodoParams{
		ID:        params.ID,
		IsDeleted: params.IsDeleted,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Todo{}, errors.NewWithCode(codes.CodeBadRequest, "Todo not found!")
		}
		return entity.Todo{}, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	// update data todo
	_, err = queries.UpdateTodo(ctx, params)
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	// record changes to histories
	histories := []entity.CreateTodoHistoryParams{}
	historyParams := map[string]any{
		"Prev": prev,
		"New":  params,
	}
	history := entity.CreateTodoHistoryParams{
		TodoID:    params.ID,
		CreatedAt: params.UpdatedAt.Time,
		CreatedBy: params.UpdatedBy.String,
	}

	if prev.Title != params.Title {
		history.Message = strformat.TWE("Rename title '{{ .Prev.Title }}' to '{{ .New.Title }}'", historyParams)
		histories = append(histories, history)
	}

	if prev.Description != params.Description {
		history.Message = strformat.TWE("Change description of '{{ .New.Title }}'", historyParams)
		histories = append(histories, history)
	}

	if prev.Status != params.Status {
		history.Message = strformat.TWE("Change status from '{{ .Prev.Status }}' to {{ .New.Status }}", historyParams)
		histories = append(histories, history)
	}

	if prev.IsDeleted == 0 && params.IsDeleted == 1 {
		history.Message = strformat.TWE("Delete todo '{{ .Prev.Title }}'", historyParams)
		histories = append(histories, history)
	}

	for _, h := range histories {
		if _, err := queries.CreateTodoHistory(ctx, h); err != nil {
			return entity.Todo{}, errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
		}
	}

	// commit changes
	if err := tx.Commit(); err != nil {
		return entity.Todo{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	result := entity.Todo{
		ID:          params.ID,
		UserID:      prev.UserID,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
		CreatedAt:   prev.CreatedAt,
		CreatedBy:   prev.CreatedBy,
		UpdatedAt:   params.UpdatedAt,
		UpdatedBy:   params.UpdatedBy,
		DeletedAt:   params.DeletedAt,
		DeletedBy:   params.DeletedBy,
		IsDeleted:   params.IsDeleted,
	}
	return result, nil
}
