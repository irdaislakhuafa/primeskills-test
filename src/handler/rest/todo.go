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

	result, err := r.u.Todo.Create(ctx.Request.Context(), body)
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

	results, pag, err := r.u.Todo.List(ctx.Request.Context(), query)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, results, &pag)
}

func (r *rest) UpdateTodo(ctx *gin.Context) {
	body := validation.UpdateTodoParams{}
	if err := ctx.BindJSON(&body); err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 16)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "invalid id"))
		return
	}
	body.ID = id

	result, err := r.u.Todo.Update(ctx.Request.Context(), body)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(errors.GetCode(err), "%s", err.Error()))
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, result, nil)
}
