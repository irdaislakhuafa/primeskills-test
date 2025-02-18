package user

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	mail "github.com/go-mail/gomail"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/irdaislakhuafa/go-sdk/auth"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/convert"
	"github.com/irdaislakhuafa/go-sdk/cryptography"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/header"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/operator"
	"github.com/irdaislakhuafa/go-sdk/smtp"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/business/domain"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/ctxkey"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/mailtemplates"
	"github.com/irdaislakhuafa/primeskills-test/src/validation"
)

type (
	Interface interface {
		Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error)
		Update(ctx context.Context, params validation.UpdateUserParams) (entity.User, error)
		List(ctx context.Context, params validation.ListUserParams) ([]entity.User, entity.Pagination, error)
		Login(ctx context.Context, params validation.LoginUserParams) (entity.User, string, error)
		RetrieveRegisterVerification(ctx context.Context, params validation.RetrieveRegisterVerificationParams) (string, error)
		RequestChangePassword(ctx context.Context, params validation.ChangePasswordParams) (string, error)
		VerifyChangePassword(ctx context.Context, params validation.VerifyChangePasswordParams) (string, error)
	}
	user struct {
		log        log.Interface
		dom        *domain.Domain
		val        *validator.Validate
		cfg        config.Config
		smtpGoMail smtp.GoMailInterface
	}
)

func Init(log log.Interface, dom *domain.Domain, val *validator.Validate, cfg config.Config, smtpGoMail smtp.GoMailInterface) Interface {
	return &user{
		log:        log,
		dom:        dom,
		val:        val,
		cfg:        cfg,
		smtpGoMail: smtpGoMail,
	}
}

func (u *user) Create(ctx context.Context, params validation.CreateUserParams) (entity.User, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	pwd, err := cryptography.NewBcrypt().Hash([]byte(params.Password))
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	params.Password = string(pwd)
	result, err := u.dom.User.Create(ctx, entity.CreateUserParams{
		Name:     params.Name,
		Password: params.Password,
		Email:    params.Email,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	m := mail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From": {u.cfg.Contacts.Email},
		"To":   {result.Email},
		"Subject": {strformat.TWE("{{ .AppName }} Verification - {{ .Email }}", map[string]string{
			"AppName": u.cfg.Meta.Title,
			"Email":   result.Email,
		})},
	})

	activationToken, err := cryptography.NewBcrypt().Hash([]byte(strformat.TWE("{{ .ID }}:{{ .Email }}", result)))
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeInternalServerError, "Cannot generate activation token")
	}

	mBody, err := mailtemplates.ReadAndParse(mailtemplates.REGISTER_VERIFICATION, map[string]any{
		"AppName": u.cfg.Meta.Title,
		"VerificationURL": strformat.TWE("{{ .Protocol }}://{{ .Host }}{{ .Port }}/api/v1/user/register/verify?uid={{ .UID }}&activation_token={{ .ActivationToken }}", map[string]any{
			"Protocol":        u.cfg.Meta.Protocol,
			"Host":            u.cfg.Meta.Host,
			"UID":             result.ID,
			"ActivationToken": string(activationToken),
			"Port":            operator.Ternary(u.cfg.Meta.Port == "", "", ":"+u.cfg.Meta.Port),
		}),
		"Contacts": u.cfg.Contacts,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}
	m.SetBody(header.ContentTypeHTML, mBody)
	if err := u.smtpGoMail.DialAndSend(m); err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	return result, nil
}

func (u *user) Update(ctx context.Context, params validation.UpdateUserParams) (entity.User, error) {
	if err := u.val.Struct(params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	result, err := u.dom.User.Update(ctx, entity.UpdateUserParams{
		Name:      params.Name,
		IsDeleted: params.IsDeleted,
		ID:        params.ID,
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return result, nil
}

func (u *user) List(ctx context.Context, params validation.ListUserParams) ([]entity.User, entity.Pagination, error) {
	if err := u.val.Struct(params); err != nil {
		err := validation.ExtractError(err, params)
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	results, pag, err := u.dom.User.List(ctx, entity.ListUserParams{
		CONCAT:    params.Search,
		CONCAT_2:  params.Search,
		IsDeleted: params.IsDeleted,
		Limit:     params.Limit,
		Offset:    params.Page,
	})
	if err != nil {
		return nil, entity.Pagination{}, errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	return results, pag, nil
}

func (u *user) Login(ctx context.Context, params validation.LoginUserParams) (entity.User, string, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err := validation.ExtractError(err, params)
		return entity.User{}, "", errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	user, err := u.dom.User.Get(ctx, entity.GetOneUserParams{
		Email:     params.Email,
		IsDeleted: 0,
	})
	if err != nil {
		if code := errors.GetCode(err); code == codes.CodeSQLRecordDoesNotExist {
			return entity.User{}, "", errors.NewWithCode(codes.CodeUnauthorized, "User not registered or already deleted!")
		} else {
			return entity.User{}, "", errors.NewWithCode(code, "%s", err.Error())
		}
	}

	if user.IsActive == 0 {
		return entity.User{}, "", errors.NewWithCode(codes.CodeUnauthorized, "Your account is not active, you need to verify your account first!")
	}

	if err := cryptography.NewBcrypt().Compare([]byte(params.Password), []byte(user.Password)); err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeUnauthorized, "Wrong password")
	}

	authJwt := auth.InitJWT([]byte(u.cfg.Secrets.Key), &entity.AuthJWTClaims{
		UID: fmt.Sprint(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: time.Now().Add(time.Minute * time.Duration(u.cfg.Token.ExpirationMinutes)),
			},
		},
	})
	token, err := authJwt.Generate(ctx)
	if err != nil {
		return entity.User{}, "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	user.Password = ""
	return user, token, nil
}

func (u *user) RetrieveRegisterVerification(ctx context.Context, params validation.RetrieveRegisterVerificationParams) (string, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		return "", errors.NewWithCode(codes.CodeBadRequest, "Request invalid")
	}
	user, err := u.dom.User.Get(ctx, entity.GetOneUserParams{
		ID:        params.UID,
		IsDeleted: 0,
	})
	if err != nil {
		if code := errors.GetCode(err); code == codes.CodeSQLRecordDoesNotExist {
			return "", errors.NewWithCode(codes.CodeBadRequest, "Request activation not found")
		}
		return "", errors.NewWithCode(codes.CodeInternalServerError, "Cannot verify your account")
	}

	token := strformat.TWE("{{.ID}}:{{.Email}}", user)
	err = cryptography.NewBcrypt().Compare([]byte(token), []byte(params.ActivationToken))
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "Cannot verify your account. Your activation token is invalid!")
	}

	if err := u.dom.User.UpdateActivationUser(ctx, entity.UpdateActivationUserParams{IsActive: 1, ID: user.ID}); err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	return strformat.TWE("Succesfully verify your account. Now you can login to {{ .Title }}", u.cfg.Meta), nil
}

func (u *user) RequestChangePassword(ctx context.Context, params validation.ChangePasswordParams) (string, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return "", errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	user, err := u.dom.User.Get(ctx, entity.GetOneUserParams{Email: params.Email})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeBadRequest, "User not registered or maybe has been deleted!")
	}

	otp, err := u.dom.Otp.Get(ctx, entity.GetOneOTPParams{
		UserID:     user.ID,
		IsUsed:     0,
		ExpirateAt: time.Now(),
	})
	if err != nil {
		if errCode := errors.GetCode(err); errCode != codes.CodeSQLRecordDoesNotExist {
			return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
		}
	} else {
		wib, _ := time.LoadLocation("Asia/Jakarta")
		return "", errors.NewWithCode(
			codes.CodeBadRequest,
			"You have otp code that has sent and not verified yet. You can send otp request again after %s for security reason!.",
			(operator.Ternary(wib == nil, otp.ExpirateAt, otp.ExpirateAt.In(wib))).Format(time.DateTime),
		)
	}

	code := fmt.Sprint(time.Now().UnixNano())
	expirateAt := time.Now().Add(time.Hour)
	code = code[len(code)-4:]
	plainCode := code
	if hashed, err := cryptography.NewBcrypt().Hash([]byte(code)); err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	} else {
		code = string(hashed)
	}

	_, err = u.dom.Otp.Create(ctx, entity.CreateOTPParams{
		UserID:     user.ID,
		Code:       code,
		ExpirateAt: expirateAt,
	})
	if err != nil {
		return "", errors.NewWithCode(errors.GetCode(err), "%s", err.Error())
	}

	m := mail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From": {u.cfg.Contacts.Email},
		"To":   {user.Email},
		"Subject": {strformat.TWE("{{ .AppName }} Change Password - {{ .Email }}", map[string]string{
			"AppName": u.cfg.Meta.Title,
			"Email":   user.Email,
		})},
	})

	wib, _ := time.LoadLocation("Asia/Jakarta")
	mBody, err := mailtemplates.ReadAndParse(mailtemplates.RESET_PASSWORD, map[string]any{
		"AppName":       u.cfg.Meta.Title,
		"Name":          user.Name,
		"ExpirateAt":    (operator.Ternary(wib == nil, expirateAt, expirateAt.In(wib))).Format(time.DateTime),
		"Code":          plainCode,
		"Contacts":      u.cfg.Contacts,
		"ExpirateAfter": "1 hour",
	})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	m.SetBody(header.ContentTypeHTML, mBody)

	if err := u.smtpGoMail.DialAndSend(m); err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	return "The otp code has been sent to your email, please use it to change your password.", nil
}

func (u *user) VerifyChangePassword(ctx context.Context, params validation.VerifyChangePasswordParams) (string, error) {
	if err := u.val.StructCtx(ctx, params); err != nil {
		err = validation.ExtractError(err, params)
		return "", errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	user, err := u.dom.User.Get(ctx, entity.GetOneUserParams{Email: params.Email})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeBadRequest, "%s", err.Error())
	}

	otp, err := u.dom.Otp.Get(ctx, entity.GetOneOTPParams{UserID: user.ID, IsUsed: 0, ExpirateAt: time.Now()})
	if err != nil {
		if errCode := errors.GetCode(err); errCode == codes.CodeSQLRecordDoesNotExist {
			return "", errors.NewWithCode(codes.CodeBadRequest, "You doesn't have any otp request!")
		}
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	bcrypt := cryptography.NewBcrypt()
	if err := bcrypt.Compare([]byte(params.OtpCode), []byte(otp.Code)); err != nil {
		return "", errors.NewWithCode(codes.CodeBadRequest, "Otp code is Invalid")
	}

	pwd, err := bcrypt.Hash([]byte(params.NewPassword))
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	err = u.dom.User.ChangePassword(ctx, entity.ChangePasswordUserParams{
		ID:        user.ID,
		Password:  string(pwd),
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		UpdatedBy: sql.NullString{Valid: true, String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID))},
	})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	_, err = u.dom.Otp.Update(ctx, entity.UpdateOTPParams{
		Code:      otp.Code,
		IsUsed:    1,
		UpdatedAt: sql.NullTime{Valid: true, Time: time.Now()},
		UpdatedBy: sql.NullString{Valid: true, String: convert.ToSafeValue[string](ctx.Value(ctxkey.USER_ID))},
		ID:        otp.ID,
	})
	if err != nil {
		return "", errors.NewWithCode(codes.CodeInternalServerError, "%s", err.Error())
	}

	return "Succesfully change your password!", nil
}
