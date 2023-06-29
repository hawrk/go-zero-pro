package handler

import (
	"algo_assess/global"
	"net/http"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ExportUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 加入token 校验
		lg := logic.NewLoginLogic(r.Context(), svcCtx)
		tokenPass, _ := lg.CheckToken(r)
		if !tokenPass {
			httpx.Error(w, global.TokenErr)
			return
		}
		// token 校验end
		l := logic.NewExportUserLogic(r.Context(), svcCtx)
		err := l.ExportUser(w)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
