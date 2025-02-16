package user

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/cryptography"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error)
		Update(ctx context.Context, params validation.UpdateUserParams) (entity.User, error)
		List(ctx context.Context, params validation.ListUserParams) ([]entity.User, error)
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

	pwd, err := cryptography.NewBcrypt().Hash([]byte(params.Password))
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	params.Password = string(pwd)
	result, err := u.dom.User.Create(ctx, entity.CreateUserParams{
		Name:     params.Name,
		Password: params.Password,
		Email:    params.Email,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}
	return result, nil
}

func (u *user) Update(ctx context.Context, params validation.UpdateUserParams) (entity.User, error) {
	if err := u.val.Struct(params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	result, err := u.dom.User.Update(ctx, entity.UpdateUserParams{
		Name:      params.Name,
		IsDeleted: params.IsDeleted,
		ID:        params.ID,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return result, nil
}

func (u *user) List(ctx context.Context, params validation.ListUserParams) ([]entity.User, error) {
	if err := u.val.Struct(params); err != nil {
		err := validation.ExtractError(err, params)
		return nil, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	results, err := u.dom.User.List(ctx, entity.ListUserParams{
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
		IsDeleted: params.IsDeleted,
		Limit:     params.Limit,
		Offset:    params.Page,
	})
	if err != nil {
		return nil, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}
	return results, nil
}
