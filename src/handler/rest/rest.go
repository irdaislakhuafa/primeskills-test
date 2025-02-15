package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/irdaislakhuafa/go-sdk/appcontext"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
	"github.com/irdaislakhuafa/go-sdk/header"
	"github.com/irdaislakhuafa/go-sdk/language"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/irdaislakhuafa/go-sdk/strformat"
	"github.com/irdaislakhuafa/primeskills-test/src/entity"
	"github.com/irdaislakhuafa/primeskills-test/src/utils/config"
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
	}
)

func Init(cfg config.Config, log log.Interface) Interface {
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
	return
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
