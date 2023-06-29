package handler

import (
	"net/http"

	"account-auth/account-auth-server/internal/logic"
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthUserListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewAuthUserListLogic(r.Context(), svcCtx)
		resp, err := l.AuthUserList(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
