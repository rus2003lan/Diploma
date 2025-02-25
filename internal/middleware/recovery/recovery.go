package recovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"go.uber.org/zap"
)

const panicStackBufSize = 2048

func Middleware(
	l *zap.SugaredLogger,
	provideAddAttrs func(ctx context.Context) []any,
) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				defer func(ctx context.Context) {
					if err := recover(); err != nil {
						buf := make([]byte, panicStackBufSize)
						buf = buf[:runtime.Stack(buf, false)]

						//nolint:goerr113
						err := fmt.Errorf(
							"recover from err: %s \n panic stack %s",
							err,
							buf,
						)

						attrs := []any{
							"err", err,
						}

						if provideAddAttrs != nil {
							attrs = append(attrs, provideAddAttrs(ctx)...)
						}

						l.Error(attrs...)

						writeError(w)
					}
				}(r.Context())

				h.ServeHTTP(w, r)
			},
		)
	}
}

func writeError(w http.ResponseWriter) {
	var (
		statusCode = http.StatusInternalServerError
		message    = http.StatusText(statusCode)
	)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	type errorResponse struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
	}

	_ = json.NewEncoder(w).Encode(
		errorResponse{
			Code:    int64(statusCode),
			Message: message,
		},
	)
}
