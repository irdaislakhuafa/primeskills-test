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
