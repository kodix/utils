// Copyright 2018 Kodix LLC. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package mw

import (
	"net/http"
	"github.com/abramd/log"
	"context"
)

type ctxKey string
var loggerKey ctxKey = "logger"

// DefLogsMw - default LogsMw() middleware with requestID = 'X-Request-Id'
func DefLogsMw(next http.Handler) http.HandlerFunc {
	return LogsMw("X-Request-Id", next)
}

// LogsMw creates logger with specified prefix (requestID)
func LogsMw(requestIDKey string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDKey)
		logger := log.Copy()
		logger.AddPrefix(reqID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), loggerKey, logger)))
	}
}

// LoggerFromCtx - logger getter from request context
func LoggerFromCtx(r *http.Request) *log.Logger {
	logger, ok := r.Context().Value(loggerKey).(*log.Logger)
	if !ok {
		return log.Copy()
	}
	return logger
}