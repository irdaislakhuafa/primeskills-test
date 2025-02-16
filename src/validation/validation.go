package validation

import (
	"database/sql"
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
)

type (
	CreateUserParams struct {
		Name      string    `db:"name" json:"name" validate:"required,min=1,max=255"`
		Password  string    `db:"password" json:"password" validate:"required,min=8,max=255"`
		Email     string    `db:"email" json:"email" validate:"required,email,max=255"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		CreatedBy string    `db:"created_by" json:"created_by"`
	}
	UpdateUserParams struct {
		Name      string         `db:"name" json:"name" validate:"required,min=1,max=255"`
		UpdatedAt sql.NullTime   `db:"updated_at" json:"updated_at"`
		UpdatedBy sql.NullString `db:"updated_by" json:"updated_by"`
		IsDeleted int8           `db:"is_deleted" json:"is_deleted" validate:"number"`
		ID        int64          `db:"id" json:"id" validate:"required,number"`
	}
	ListUserParams entity.ListUserParams

	CreateTodoParams struct {
		UserID      int64   `db:"user_id" json:"user_id" validate:"required"`
		Title       string  `db:"title" json:"title" validate:"required,min=1,max=255"`
		Description *string `db:"description" json:"description" validate:""`
		Status      string  `db:"status" json:"status" validate:"required,oneof=complete cancel hold todo"`
	}

	ListTodoParams struct {
		UserID    int64  `db:"user_id" json:"user_id" validate:"number,required"`
		Status    string `db:"status" json:"status" validate:""`
		Search    string `json:"search" validate:""`
		Page      int64  `json:"page" validate:""`
		Limit     int64  `json:"limit" validate:""`
		IsDeleted int8   `json:"is_deleted" validate:""`
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
