package todo_histories

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		List(ctx context.Context, params validation.ListTodoHistories) ([]entity.TodoHistory, error)
	}

	todoHistory struct {
		log log.Interface
		dom *domain.Domain
		val *validator.Validate
	}
)

func Init(log log.Interface, dom *domain.Domain, val *validator.Validate) Interface {
	return &todoHistory{
		log: log,
		dom: dom,
		val: val,
	}
}

func (th *todoHistory) List(ctx context.Context, params validation.ListTodoHistories) ([]entity.TodoHistory, error) {
	if err := th.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return nil, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	results, err := th.dom.TodoHistory.List(ctx, entity.ListTodoHistoriesParams{
		TodoID:    params.TodoID,
		IsDeleted: params.IsDeleted,
		Limit:     int32(params.Limit),
		Offset:    int32(params.Page),
	})
	if err != nil {
		return nil, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return results, nil
}
