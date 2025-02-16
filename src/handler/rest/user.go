package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/irdaislakhuafa/go-sdk/codes"
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
