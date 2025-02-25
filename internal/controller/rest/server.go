package webapi

import (
	"context"
	"fmt"
	"net/http"

	api "diploma-project/api/web/gen"

	"github.com/getkin/kin-openapi/openapi3filter"
	"go.uber.org/zap"

	"diploma-project/internal/middleware/recovery"

	"github.com/go-chi/chi/v5"
	validator "github.com/oapi-codegen/nethttp-middleware"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	ReportHandler
	DocHandler
	SQLMapHandler
}

type Opts struct {
	Port     int
	Handlers Handler

	ValidateErrorFunc    func(w http.ResponseWriter, r *http.Request, err error)
	ErrorRespFunc        func(w http.ResponseWriter, r *http.Request, err error)
	ValidateErrorHandler validator.ErrorHandler
}

type Server struct {
	r    *chi.Mux
	port int
}

func New(
	opts Opts,
	lg *zap.SugaredLogger,
) *Server {
	r := chi.NewRouter()

	ssi := Handler{
		ReportHandler: opts.Handlers.ReportHandler,
		DocHandler:    opts.Handlers.DocHandler,
		SQLMapHandler: opts.Handlers.SQLMapHandler,
	}

	api.HandlerWithOptions(
		api.NewStrictHandlerWithOptions(
			ssi,
			nil,
			api.StrictHTTPServerOptions{
				RequestErrorHandlerFunc:  opts.ValidateErrorFunc,
				ResponseErrorHandlerFunc: opts.ErrorRespFunc,
			},
		),
		api.ChiServerOptions{
			BaseURL:          "",
			BaseRouter:       r,
			Middlewares:      getMiddlewares(opts.ValidateErrorHandler, lg),
			ErrorHandlerFunc: opts.ValidateErrorFunc,
		},
	)

	return &Server{
		r:    r,
		port: opts.Port,
	}
}

func (s *Server) Start() error {
	//nolint:gosec, wrapcheck
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.r)
}

func getMiddlewares(
	validateErrorHandler validator.ErrorHandler,
	l *zap.SugaredLogger,
) []api.MiddlewareFunc {
	doc, _ := api.GetSwagger()

	openapi3filter.RegisterBodyDecoder("application/json", openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder("text/plain", openapi3filter.FileBodyDecoder)

	return []api.MiddlewareFunc{
		recovery.Middleware(l, provideLoggerAttrs),
		validator.OapiRequestValidatorWithOptions(doc, &validator.Options{
			ErrorHandler: validateErrorHandler,
		}),
	}
}

func provideLoggerAttrs(ctx context.Context) []any {
	attrs := []any{}

	attrs = append(attrs,
		"trace_id", trace.SpanContextFromContext(ctx).TraceID().String(),
	)

	return attrs
}
