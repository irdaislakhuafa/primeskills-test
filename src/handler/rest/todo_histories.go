package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

func (r *rest) ListTodoHistories(ctx *gin.Context) {
	query := validation.ListTodoHistories{Limit: 15}
	if err := ctx.BindQuery(&query); err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}

	results, err := r.u.TodoHistory.List(ctx, query)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(errors.GetCode(err), "%s", err.Error()))
		return
	}
	r.httpRespSuccess(ctx, codes.CodeSuccess, results, nil)
}
