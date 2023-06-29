package handler

import (
	"algo_assess/global"
	"net/http"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserProfileExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProfileExportReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		// 加入token 校验
		lg := logic.NewLoginLogic(r.Context(), svcCtx)
		tokenPass, _ := lg.CheckToken(r)
		if !tokenPass {
			httpx.Error(w, global.TokenErr)
			return
		}
		// token 校验end
		l := logic.NewUserProfileExportLogic(r.Context(), svcCtx)
		err := l.UserProfileExport(w, &req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
