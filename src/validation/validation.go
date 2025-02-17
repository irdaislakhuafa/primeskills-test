package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/strformat"
)

type (
	CreateUserParams struct {
		Name     string `json:"name" validate:"required,min=1,max=255"`
		Password string `json:"password" validate:"required,min=8,max=255"`
		Email    string `json:"email" validate:"required,email,max=255"`
	}
	UpdateUserParams struct {
		Name      string `json:"name" validate:"required,min=1,max=255"`
		IsDeleted int8   `json:"is_deleted" validate:"number"`
		ID        int64  `json:"id" validate:"required,number"`
	}

	ListUserParams struct {
		Search    string `json:"search" form:"search"`
		IsDeleted int8   `json:"is_deleted" form:"is_deleted"`
		Limit     int32  `json:"limit" form:"limit"`
		Page      int32  `json:"page" form:"page"`
	}
	LoginUserParams struct {
		Email    string `json:"email" validate:"required,email,max=255,min=0"`
		Password string `json:"password" validate:"required,min=8,max=255"`
	}

	RetrieveRegisterVerificationParams struct {
		UID             int64  `json:"uid" form:"uid" validate:"required,number"`
		ActivationToken string `json:"activation_token" form:"activation_token" validate:"required"`
	}

	CreateTodoParams struct {
		UserID      int64  `json:"user_id" validate:"required"`
		Title       string `json:"title" validate:"required,min=1,max=255"`
		Description string `json:"description" validate:""`
		Status      string `json:"status" validate:"required,oneof=complete cancel hold todo"`
	}

	ListTodoParams struct {
		UserID    int64  `json:"user_id" form:"user_id" validate:"number,required"`
		Status    string `json:"status" form:"status" validate:""`
		Search    string `json:"search" form:"search" validate:""`
		Page      int64  `json:"page" form:"page" validate:""`
		Limit     int64  `json:"limit" form:"limit" validate:""`
		IsDeleted int8   `json:"is_deleted" form:"is_deleted" validate:""`
	}

	UpdateTodoParams struct {
		Title       string `db:"title" json:"title"`
		Description string `db:"description" json:"description"`
		Status      string `db:"status" json:"status"`
		IsDeleted   int8   `db:"is_deleted" json:"is_deleted"`
		ID          int64  `db:"id" json:"id"`
	}

	ListTodoHistories struct {
		TodoID    int64 `json:"todo_id" form:"todo_id" validate:"required,number"`
		IsDeleted int8  `json:"is_deleted" form:"is_deleted" validate:"number"`
		Limit     int   `json:"limit" form:"limit" validate:"min=0"`
		Page      int   `json:"page" form:"page" validate:"min=0"`
	}
)

var customMessages = map[string]string{
	"required": "Field '{{ .Field }}' is required",
	"email":    "Field '{{ .Field }}' must be a valid email format",
	"max":      "Field '{{ .Field }}' cannot exceed {{ .Param }} characters",
	"min":      "Field '{{ .Field }}' must be at least {{ .Param }} characters",
	"gte":      "Field '{{ .Field }}' must be greater than or equal to {{ .Param }}",
	"lte":      "Field '{{ .Field }}' must be less than or equal to {{ .Param }}",
	"oneof":    "Field '{{ .Field }}' must be one of [{{ .Param }}]",
}

func ExtractError(err error, val any) error {
	if err, isOk := err.(validator.ValidationErrors); isOk {
		err := err[0]
		msg := customMessages[err.Tag()]
		msg = strformat.TWE(msg, map[string]string{
			"Field": getJSONTag(reflect.TypeOf(val), err.Field()),
			"Param": err.Param(),
		})
		return errors.NewWithCode(codes.CodeBadRequest, "%s", msg)
	}
	return errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
}

func getJSONTag(ref reflect.Type, fieldName string) string {
	sf, isOk := ref.FieldByName(fieldName)
	if !isOk {
		return fieldName
	}

	tagVal := sf.Tag.Get("json")
	if tagVal == "" {
		return fieldName
	}

	return tagVal
}
