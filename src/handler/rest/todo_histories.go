package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

func (r *rest) ListTodoHistories(ctx *gin.Context) {
	query := validation.ListTodoHistories{}
	if err := ctx.BindQuery(&query); err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error()))
		return
	}
	r.httpRespSuccess(ctx, codes.CodeSuccess, query, nil)
}
