package handler

import (
	"net/http"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GeneralHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GeneralReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGeneralLogic(r.Context(), svcCtx)
		resp, err := l.General(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
