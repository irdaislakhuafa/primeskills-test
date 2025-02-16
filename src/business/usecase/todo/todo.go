package todo

import (
	"context"
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		Create(ctx context.Context, params validation.CreateTodoParams) (entity.Todo, error)
		List(ctx context.Context, params validation.ListTodoParams) ([]entity.Todo, error)
	}
	todo struct {
		log log.Interface
		val *validator.Validate
		dom *domain.Domain
	}
)

func Init(log log.Interface, val *validator.Validate, dom *domain.Domain) Interface {
	return &todo{
		log: log,
		val: val,
		dom: dom,
	}
}

func (t *todo) Create(ctx context.Context, params validation.CreateTodoParams) (entity.Todo, error) {
	if err := t.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.Todo{}, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	result, err := t.dom.Todo.Create(ctx, entity.CreateTodoParams{
		UserID: params.UserID,
		Title:  params.Title,
		Description: sql.NullString{
			Valid:  params.Description != nil,
			String: convert.ToSafeValue[string](params.Description),
		},
		Status: params.Status,
	})
	if err != nil {
		return entity.Todo{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return result, nil
}

func (t *todo) List(ctx context.Context, params validation.ListTodoParams) ([]entity.Todo, error) {
	if err := t.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return nil, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	results, err := t.dom.Todo.List(ctx, entity.ListTodoParams{
		UserID:    params.UserID,
		Status:    params.Status,
		IsDeleted: params.IsDeleted,
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Page),
	})
	if err != nil {
		return nil, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return results, nil
}
