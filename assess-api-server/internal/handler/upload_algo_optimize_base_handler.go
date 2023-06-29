package handler

import (
	"net/http"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func uploadAlgoOptimizeBaseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewUploadAlgoOptimizeBaseLogic(r.Context(), svcCtx)
		resp, err := l.UploadAlgoOptimizeBase(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
