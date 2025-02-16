package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase/user"
)

type (
	Usecase struct {
		User user.Interface
	}
)

func Init(d *domain.Domain, log log.Interface, v *validator.Validate) *Usecase {
	return &Usecase{
		User: user.Init(log, d, v),
	}
}
