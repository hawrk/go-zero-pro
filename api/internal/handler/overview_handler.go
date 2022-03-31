package handler

import (
	"net/http"

	"algo_assess/api/internal/logic"
	"algo_assess/api/internal/svc"
	"algo_assess/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OverviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OverviewReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewOverviewLogic(r.Context(), svcCtx)
		resp, err := l.Overview(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
