package server

import (
	"go.uber.org/zap"
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,PUT,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "application/json")
		//log request example
		zapLogger, _ := zap.NewProduction()
		defer zapLogger.Sync()
		zapLogger.Info(r.URL.String(),
			zap.String("method", r.Method),
			zap.String("host", r.Host),
		)

		next.ServeHTTP(w, r)
		//log response

	})
}
