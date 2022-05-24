package middleware

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type InterceptorMiddleware struct {
}

func NewInterceptorMiddleware() *InterceptorMiddleware {
	return &InterceptorMiddleware{}
}

func (m *InterceptorMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("Interceptor middle ware")
		next(w, r)
	}
}
