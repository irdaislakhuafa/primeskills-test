package user

import (
	"context"
	"database/sql"
	"time"

	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/ctxkey"
)

type (
	Interface interface {
		Create(ctx context.Context, params entity.CreateUserParams) (entity.User, error)
		Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error)
		List(ctx context.Context) ([]entity.User, error)
	}
	user struct {
		log     log.Interface
		queries *entity.Queries
		db      *sql.DB
	}
)

func Init(log log.Interface, queries *entity.Queries, db *sql.DB) Interface {
	return &user{
		log:     log,
		queries: queries,
		db:      db,
	}
}

func (u *user) Create(ctx context.Context, params entity.CreateUserParams) (entity.User, error) {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	queries := u.queries.WithTx(tx)

	params.CreatedAt = time.Now()
	params.CreatedBy = convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID))

	r, err := queries.CreateUser(ctx, params)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	result := entity.User{
		Name:      params.Name,
		Email:     params.Email,
		CreatedAt: params.CreatedAt,
		CreatedBy: params.CreatedBy,
		IsDeleted: 0,
	}
	result.ID, err = r.LastInsertId()
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	return result, nil
}

func (u *user) Update(ctx context.Context, params entity.UpdateUserParams) (entity.User, error) {
	tx, err := u.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	queries := u.queries.WithTx(tx)

	params.UpdatedAt = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	params.UpdatedBy = sql.NullString{
		Valid:  true,
		String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID)),
	}

	_, err = queries.UpdateUser(ctx, params)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	row, err := queries.GetOneUser(ctx, entity.GetOneUserParams{
		ID:        params.ID,
		IsDeleted: params.IsDeleted,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	return entity.User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
		CreatedBy: row.CreatedBy,
		UpdatedAt: row.UpdatedAt,
		UpdatedBy: row.UpdatedBy,
		DeletedAt: row.DeletedAt,
		DeletedBy: row.DeletedBy,
		IsDeleted: row.IsDeleted,
	}, nil
}

func (u *user) List(ctx context.Context) ([]entity.User, error) {
	panic("")
}
