package otp

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
		Create(ctx context.Context, params entity.CreateOTPParams) (entity.Otp, error)
		Update(ctx context.Context, params entity.UpdateOTPParams) (entity.Otp, error)
		Get(ctx context.Context, params entity.GetOneOTPParams) (entity.Otp, error)
	}
	otp struct {
		log     log.Interface
		db      *sql.DB
		queries *entity.Queries
	}
)

func Init(log log.Interface, db *sql.DB, queries *entity.Queries) Interface {
	return &otp{
		log:     log,
		db:      db,
		queries: queries,
	}
}

func (o *otp) Create(ctx context.Context, params entity.CreateOTPParams) (entity.Otp, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	queries := o.queries.WithTx(tx)
	params.CreatedAt = time.Now()
	params.CreatedBy = convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID))
	r, err := queries.CreateOTP(ctx, params)
	if err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	result := entity.Otp{
		UserID:    params.UserID,
		Code:      params.Code,
		IsUsed:    0,
		CreatedAt: params.CreatedAt,
		CreatedBy: params.CreatedBy,
	}
	if result.ID, err = r.LastInsertId(); err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	return result, nil
}

func (o *otp) Update(ctx context.Context, params entity.UpdateOTPParams) (entity.Otp, error) {
	tx, err := o.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxBegin, "%s", err.Error())
	}
	defer tx.Rollback()

	queries := o.queries.WithTx(tx)

	// ensure otp is exist
	prev, err := queries.GetOneOTP(ctx, entity.GetOneOTPParams{ID: params.ID})
	if err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}
	if prev.IsUsed == 1 {
		return entity.Otp{}, errors.NewWithCode(codes.CodeBadRequest, "Otp code already used!")
	}

	// fill update fields
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

	if _, err := queries.UpdateOTP(ctx, params); err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxExec, "%s", err.Error())
	}

	result := entity.Otp{
		ID:        params.ID,
		UserID:    prev.ID,
		Code:      prev.Code,
		IsUsed:    params.IsUsed,
		CreatedAt: prev.CreatedAt,
		CreatedBy: prev.CreatedBy,
		UpdatedAt: params.UpdatedAt,
		UpdatedBy: params.DeletedBy,
		DeletedAt: params.DeletedAt,
		DeletedBy: params.DeletedBy,
		IsDeleted: params.IsDeleted,
	}

	if err := tx.Commit(); err != nil {
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLTxCommit, "%s", err.Error())
	}

	return result, nil
}

func (o *otp) Get(ctx context.Context, params entity.GetOneOTPParams) (entity.Otp, error) {
	row, err := o.queries.GetOneOTP(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Otp{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "%s", err.Error())
		}
		return entity.Otp{}, errors.NewWithCode(codes.CodeSQLRead, "%s", err.Error())
	}

	result := entity.Otp{
		ID:        row.ID,
		UserID:    row.UserID,
		Code:      row.Code,
		IsUsed:    row.IsUsed,
		CreatedAt: row.CreatedAt,
		CreatedBy: row.CreatedBy,
		UpdatedAt: row.UpdatedAt,
		UpdatedBy: row.UpdatedBy,
		DeletedAt: row.DeletedAt,
		DeletedBy: row.DeletedBy,
		IsDeleted: row.IsDeleted,
	}
	return result, nil
}
