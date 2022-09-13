package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/Ollub/user_service/pkg/log"
)

const RequestIDKey = "requestID"

func SetupReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// https://github.com/opentracing/specification/blob/master/rfc/trace_identifiers.md
			requestID = RandBytesHex(16)
			r.Header.Set("X-Request-ID", requestID)
			r.Header.Set("trace-id", requestID)
			w.Header().Set("trace-id", requestID)
			w.Header().Set("X-Request-ID", requestID)
		}
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func InjectLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.New()
		logger.SetContext(log.Fields{"trace-id": RequestIDFromContext(r.Context())})

		ctx := context.WithValue(r.Context(), log.LoggerKey, logger)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func SetupAccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		log.Rlog(r).Info(
			r.URL.Path,
			log.Fields{
				"method":      r.Method,
				"remote_addr": r.RemoteAddr,
				"url":         r.URL.Path,
				"work_time":   time.Since(start),
			},
		)
	})
}

func RandBytesHex(n int) string {
	return fmt.Sprintf("%x", RandBytes(n))
}

func RandBytes(n int) []byte {
	res := make([]byte, n)
	rand.Read(res)
	return res
}

func RequestIDFromContext(ctx context.Context) string {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		return "-"
	}
	return requestID
}
