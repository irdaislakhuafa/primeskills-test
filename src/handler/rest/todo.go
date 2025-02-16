package rest

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

func (r *rest) CreateTodo(ctx *gin.Context) {
	body := validation.CreateTodoParams{}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
	}

	result, err := r.u.Todo.Create(ctx, body)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, result, nil)
}

func (r *rest) ListTodo(ctx *gin.Context) {
	userID, err := strconv.ParseInt(ctx.DefaultQuery("user_id", "0"), 10, 64)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "invalid user id"))
		return
	}

	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	if err != nil {
		page = 0
	}

	limit, err := strconv.ParseInt(ctx.DefaultQuery("limit", "15"), 10, 64)
	if err != nil {
		limit = 15
	}

	isDeleted, err := strconv.ParseInt(ctx.DefaultQuery("is_deleted", "0"), 10, 64)
	if err != nil {
		isDeleted = 0
	}

	query := validation.ListTodoParams{
		UserID:    userID,
		Status:    ctx.DefaultQuery("status", ""),
		Search:    ctx.DefaultQuery("search", ""),
		Page:      page,
		Limit:     limit,
		IsDeleted: int8(isDeleted),
	}

	results, err := r.u.Todo.List(ctx, query)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, results, nil)
}
