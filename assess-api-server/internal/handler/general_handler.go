package handler

import (
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"time"

	"algo_assess/assess-api-server/internal/logic"
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var upGrader = websocket.Upgrader{
	HandshakeTimeout: 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GeneralHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GeneralReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		logx.Infof("get req:%+v", req)
		l := logic.NewGeneralLogic(r.Context(), svcCtx)
		// 升级为web socket
		if req.WebSocket == 1 { // 支持websocket
			ws, err := upGrader.Upgrade(w, r, nil)
			if err != nil {
				httpx.Error(w, err)
				return
			}
			defer func() {
				ws.Close()
			}()
			for {
				resp, _ := l.General(&req)
				ws.WriteJSON(resp)
				time.Sleep(time.Second * time.Duration(svcCtx.Config.WebSocket.DurationTime))
			}
		} else {
			resp, err := l.General(&req)
			if err != nil {
				httpx.Error(w, err)
			} else {
				httpx.OkJson(w, resp)
			}
		}
	}
}
