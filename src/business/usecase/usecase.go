package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/smtp"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase/todo"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase/todo_histories"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase/user"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
)

type (
	Usecase struct {
		User        user.Interface
		Todo        todo.Interface
		TodoHistory todo_histories.Interface
	}
)

func Init(d *domain.Domain, log log.Interface, v *validator.Validate, cfg config.Config, smtpGoMail smtp.GoMailInterface) *Usecase {
	return &Usecase{
		User:        user.Init(log, d, v, cfg, smtpGoMail),
		Todo:        todo.Init(log, v, d),
		TodoHistory: todo_histories.Init(log, d, v),
	}
}
