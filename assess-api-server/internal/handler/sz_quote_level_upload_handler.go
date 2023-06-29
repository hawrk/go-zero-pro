package handler

import (
	"net/http"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SzQuoteLevelUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewSzQuoteLevelUploadLogic(r.Context(), svcCtx)
		resp, err := l.SzQuoteLevelUpload(r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
