package handler

import (
	"net/http"

	"account-auth/account-auth-server/internal/logic"
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckPasswordReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCheckPasswordLogic(r.Context(), svcCtx)
		resp, err := l.CheckPassword(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
