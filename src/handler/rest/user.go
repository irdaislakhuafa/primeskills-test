package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

func (r *rest) CreateUser(ctx *gin.Context) {
	body := validation.CreateUserParams{}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	result, err := r.u.User.Create(ctx, body)
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

	result, err := r.u.User.Update(ctx, body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, result, nil)
}

func (r *rest) ListUser(ctx *gin.Context) {
	isDeleted, err := strconv.ParseInt(ctx.DefaultQuery("is_deleted", "0"), 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}

	limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "15"), 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}

	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}
	search := ctx.DefaultQuery("search", "")
	body := validation.ListUserParams{
		CONCAT:    search,
		CONCAT_2:  search,
		IsDeleted: int8(isDeleted),
		Limit:     int32(limit),
		Offset:    int32(page),
	}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	results, err := r.u.User.List(ctx, body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, results, &entity.Pagination{})
}
