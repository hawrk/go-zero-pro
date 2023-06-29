package handler

import (
	"net/http"

	"account-auth/account-auth-server/internal/logic"
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewUserAuthLogic(r.Context(), svcCtx)
		resp, err := l.UserAuth(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
