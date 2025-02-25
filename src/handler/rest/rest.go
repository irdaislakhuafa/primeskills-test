package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/irdaislakhuafa/go-sdk/appcontext"
	"github.com/irdaislakhuafa/go-sdk/auth"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/header"
	"github.com/irdaislakhuafa/go-sdk/language"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/business/usecase"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/ctxkey"
)

var once = &sync.Once{}

type (
	Interface interface {
		Run()
	}
	rest struct {
		cfg config.Config
		svr *gin.Engine
		log log.Interface
		u   *usecase.Usecase
	}
)

func Init(cfg config.Config, log log.Interface, u *usecase.Usecase) Interface {
	r := &rest{}
	once.Do(func() {
		modes := map[string]string{
			gin.DebugMode:   gin.DebugMode,
			gin.TestMode:    gin.TestMode,
			gin.ReleaseMode: gin.ReleaseMode,
		}

		gin.SetMode(modes[cfg.Gin.Mode])

		svr := gin.New()
		r = &rest{
			cfg: cfg,
			svr: svr,
			log: log,
			u:   u,
		}

		if cfg.Gin.Cors.Mode == "allowall" {
			r.svr.Use(cors.New(cors.Config{
				AllowAllOrigins: true,
				AllowHeaders:    []string{"*"},
				AllowMethods: []string{
					http.MethodPost,
					http.MethodDelete,
					http.MethodGet,
					http.MethodOptions,
					http.MethodPut,
				},
			}))
		} else {
			r.svr.Use(cors.New(cors.DefaultConfig()))
		}

		// add metadata fields to context
		r.svr.Use(r.addFieldsToContext)

		// enable gin recovery on panic app
		r.svr.Use(gin.Recovery())

		// set timeout gin server
		r.svr.Use(r.SetTimeout)

		// register route
		r.Register()
	})

	return r
}

func (r *rest) Run() {
	if r.cfg.Gin.Port != "" {
		if err := r.svr.Run(strformat.TWE(":{{ .Port }}", r.cfg.Gin)); err != nil {
			r.log.Fatal(context.Background(), err.Error())
		}
	} else {
		if err := r.svr.Run(":8000"); err != nil {
			r.log.Fatal(context.Background(), err.Error())
		}
	}
}

func (r *rest) Register() {
	// server health and testing purpose
	r.svr.GET("/ping", r.Ping)

	api := r.svr.Group("/api")
	v1 := api.Group("/v1")
	{

		v1.POST("/user/register", r.CreateUser)
		v1.GET("/user/register/verify", r.RetrieveRegisterVerification)
		v1.POST("/user/login", r.LoginUser)
		v1.POST("/user/change/password", r.ChangePasswordUser)
		v1.POST("/user/change/password/verify", r.VerifyChangePasswordUser)

		user := v1.Group("/users", r.addJwtAuth)
		{
			user.POST("/:id", r.UpdateUser)
			user.GET("/", r.ListUser)
		}

		todo := v1.Group("/todos", r.addJwtAuth)
		{
			todo.POST("/", r.CreateTodo)
			todo.GET("/", r.ListTodo)
			todo.POST("/:id", r.UpdateTodo)
		}

		todoHistory := v1.Group("/todo/histories", r.addJwtAuth)
		{
			todoHistory.GET("/", r.ListTodoHistories)
		}
	}
}

func (r *rest) SetTimeout(ctx *gin.Context) {
	// wrap context with timeout
	to := time.Duration(time.Second * time.Duration(r.cfg.Gin.TimeoutSeconds))
	c, cancel := context.WithTimeout(ctx.Request.Context(), to)
	defer func() {
		// if context timeout was reached then write response and abort the request
		if c.Err() == context.DeadlineExceeded {
		}

		cancel()
	}()
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}

func (r *rest) addJwtAuth(ctx *gin.Context) {
	headerAuth := ctx.GetHeader(header.KeyAuthorization)
	if headerAuth == "" {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeUnauthorized, "Unauthorized"))
		return
	}
	const BEARIER = "Bearer "
	if !strings.HasPrefix(headerAuth, BEARIER) {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeUnauthorized, "Unauthorized"))
		return
	}

	token := strings.ReplaceAll(headerAuth, BEARIER, "")

	authJwt := auth.InitJWT([]byte(r.cfg.Secrets.Key), &entity.AuthJWTClaims{})
	jToken, err := authJwt.Validate(ctx, token)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeUnauthorized, "%s", err.Error()))
		return
	}

	claims, err := authJwt.ExtractClaims(ctx, jToken)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(errors.GetCode(err), "%s", err.Error()))
		return
	}

	ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), ctxkey.USER_ID, claims.UID))
	ctx.Next()
}

func (r *rest) addFieldsToContext(ctx *gin.Context) {
	reqID := ctx.GetHeader(header.KeyRequestID)
	if reqID == "" {
		reqID = uuid.New().String()
	}

	// override context with fields
	c := ctx.Request.Context()
	c = appcontext.SetRequestID(c, reqID)
	c = appcontext.SetAcceptLanguage(c, language.Language(ctx.Request.Header.Get(header.KeyAcceptLanguage)))
	c = appcontext.SetUserAgent(c, ctx.Request.Header.Get(header.KeyUserAgent))
	c = appcontext.SetServiceVersion(c, r.cfg.Meta.Version)

	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}

func (r *rest) httpRespError(ctx *gin.Context, err error) {
	httpStatusCode, displayErr := errors.Compile(err, appcontext.GetAcceptLanguage(ctx))
	statusStr := http.StatusText(httpStatusCode)

	c := ctx.Request.Context()
	res := entity.HTTPRes{
		Message: entity.HTTPMessage{
			Title: displayErr.Title,
			Body:  displayErr.Body,
		},
		Meta: entity.Meta{
			Path:       r.cfg.Meta.Host + ctx.Request.URL.String(),
			StatusCode: httpStatusCode,
			StatusStr:  statusStr,
			Message: strformat.TmplWithoutErr("{{ .Method }} {{ .URI }} [{{ .StatusCode }}] {{ .StatusStr }}", map[string]any{
				"Method":     ctx.Request.Method,
				"URI":        ctx.Request.URL.RequestURI(),
				"StatusCode": httpStatusCode,
				"statusStr":  statusStr,
			}),
			Timestamp: time.Now().Format(time.RFC3339),
			Error: &entity.MetaError{
				Code:    int(displayErr.Code),
				Message: err.Error(),
			},
			RequestID: appcontext.GetRequestID(c),
		},
		Data:       nil,
		Pagination: nil,
	}

	r.log.Error(c, err)
	ctx.Header(header.KeyRequestID, appcontext.GetRequestID(c))
	ctx.AbortWithStatusJSON(httpStatusCode, res)
}

func (r *rest) httpRespSuccess(ctx *gin.Context, code codes.Code, data any, p *entity.Pagination) {
	successApp := codes.Compile(code, appcontext.GetAcceptLanguage(ctx))
	c := ctx.Request.Context()
	res := entity.HTTPRes{
		Message: entity.HTTPMessage{
			Title: successApp.Title,
			Body:  successApp.Body,
		},
		Meta: entity.Meta{
			Path:       r.cfg.Meta.Host + ctx.Request.URL.String(),
			StatusCode: successApp.StatusCode,
			StatusStr:  http.StatusText(successApp.StatusCode),
			Message: strformat.TmplWithoutErr("{{ .Method }} {{ .URI }} [{{ .StatusCode }}] {{ .StatusStr }}", map[string]any{
				"Method":     ctx.Request.Method,
				"URI":        ctx.Request.URL.RequestURI(),
				"StatusCode": successApp.StatusCode,
				"statusStr":  http.StatusText(successApp.StatusCode),
			}),
			Timestamp: time.Now().Format(time.RFC3339),
			Error:     nil,
			RequestID: appcontext.GetRequestID(ctx),
		},
		Data:       data,
		Pagination: p,
	}

	raw, err := json.Marshal(res)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, "Cannot marshal response"))
		return
	}

	ctx.Header(header.KeyRequestID, appcontext.GetRequestID(c))
	ctx.Data(successApp.StatusCode, header.ContentTypeJSON, raw)
}

/**
 * @Summary Heatlh Check
 * @Description This endpoint will hit the server
 * @Tags server
 * @Produce JSON
 * @Success 200 string example="PONG!"
 * @Router /ping [GET]
 */
func (r *rest) Ping(ctx *gin.Context) {
	r.httpRespSuccess(ctx, codes.CodeSuccess, "PONG!", nil)
}
