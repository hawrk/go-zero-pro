// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"algo_assess/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/algo-demo",
				Handler: DemoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/algo-overview",
				Handler: OverviewHandler(serverCtx),
			},
		},
		rest.WithPrefix("/algo-assess/v1"),
	)
}
