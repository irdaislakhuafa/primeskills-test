package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

func (r *rest) CreateUser(ctx *gin.Context) {
	body := validation.CreateUserParams{}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	result, err := r.u.User.Create(ctx.Request.Context(), body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, result, nil)
}

func (r *rest) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "invalid id"))
		return
	}

	body := validation.UpdateUserParams{ID: id}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, err)
		return
	}
	r.log.Info(ctx, body)

	result, err := r.u.User.Update(ctx.Request.Context(), body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, result, nil)
}

func (r *rest) ListUser(ctx *gin.Context) {
	body := validation.ListUserParams{Limit: 15}
	if err := ctx.BindQuery(&body); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	results, pag, err := r.u.User.List(ctx.Request.Context(), body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, results, &pag)
}

func (r *rest) LoginUser(ctx *gin.Context) {
	body := validation.LoginUserParams{}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}

	result, token, err := r.u.User.Login(ctx.Request.Context(), body)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(errors.GetCode(err), "%s", err.Error()))
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, map[string]any{"user": result, "token": token}, nil)
}

func (r *rest) RetrieveRegisterVerification(ctx *gin.Context) {
	body := validation.RetrieveRegisterVerificationParams{}
	if err := ctx.BindQuery(&body); err != nil {
		ctx.String(200, "Invalid Request")
		return
	}

	msg, err := r.u.User.RetrieveRegisterVerification(ctx.Request.Context(), body)
	if err != nil {
		ctx.String(400, err.Error())
		return
	}

	ctx.String(200, msg)
}
