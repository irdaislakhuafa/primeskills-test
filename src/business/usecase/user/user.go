package user

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/irdaislakhuafa/go-sdk/auth"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/cryptography"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error)
		Update(ctx context.Context, params validation.UpdateUserParams) (entity.User, error)
		List(ctx context.Context, params validation.ListUserParams) ([]entity.User, entity.Pagination, error)
		Login(ctx context.Context, params validation.LoginUserParams) (entity.User, string, error)
	}
	user struct {
		log log.Interface
		dom *domain.Domain
		val *validator.Validate
		cfg config.Config
	}
)

func Init(log log.Interface, dom *domain.Domain, val *validator.Validate, cfg config.Config) Interface {
	return &user{
		log: log,
		dom: dom,
		val: val,
		cfg: cfg,
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

func (u *user) List(ctx context.Context, params validation.ListUserParams) ([]entity.User, entity.Pagination, error) {
	if err := u.val.Struct(params); err != nil {
		err := validation.ExtractError(err, params)
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	results, pag, err := u.dom.User.List(ctx, entity.ListUserParams{
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
		IsDeleted: params.IsDeleted,
		Limit:     params.Limit,
		Offset:    params.Page,
	})
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}
	return results, pag, nil
}

func (u *user) Login(ctx context.Context, params validation.LoginUserParams) (entity.User, string, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, "", errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	user, err := u.dom.User.Get(ctx, entity.GetOneUserParams{
		Email:     params.Email,
		IsDeleted: 0,
	})
	if err != nil {
		if code := errors.GetCode(err); code == codes.CodeSQLRecordDoesNotExist {
			return entity.User{}, "", errors.NewWithCode(codes.CodeUnauthorized, "User not registered or already deleted!")
		} else {
			return entity.User{}, "", errors.NewWithCode(code, "%s", err.Error())
		}
	}

	if err := cryptography.NewBcrypt().Compare([]byte(params.Password), []byte(user.Password)); err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeUnauthorized, "Wrong password")
	}

	authJwt := auth.InitJWT([]byte(u.cfg.Secrets.Key), &entity.AuthJWTClaims{
		UID: fmt.Sprint(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * time.Duration(u.cfg.Token.ExpirationMinutes)),
			},
		},
	})
	token, err := authJwt.Generate(ctx)
	if err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	return user, token, nil
}
