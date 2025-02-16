package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error)
	}
	user struct {
		log log.Interface
		dom *domain.Domain
		val *validator.Validate
	}
)

func Init(log log.Interface, dom *domain.Domain, val *validator.Validate) Interface {
	return &user{
		log: log,
		dom: dom,
		val: val,
	}
}

func (u *user) Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	result, err := u.dom.User.Create(ctx, entity.CreateUserParams(params))
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}
	return result, nil
}
