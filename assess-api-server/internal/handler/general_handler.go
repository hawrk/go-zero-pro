package handler

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"bytes"
	"encoding/binary"
	"errors"
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
		if svcCtx.Config.HRedis.NeedCheck {
			if err := checkLoginToken(svcCtx, r); err != nil {
				httpx.Error(w, errors.New("登陆鉴权失败"))
				return
			}
		}

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

func parseRedisSession(str string) (session global.Session, err error) {
	buf := bytes.NewReader([]byte(str))
	if err := binary.Read(buf, binary.LittleEndian, &session); err != nil {
		return global.Session{}, err
	}
	return session, nil
}

func checkLoginToken(svcCtx *svc.ServiceContext, r *http.Request) error {
	// logx.Infof("get header:%+v", r.Header)
	// 校验一下token
	out, err := svcCtx.HRedisClient.HGet(r.Context(), global.TokenKey, r.Header.Get("Id")).Result()
	if err != nil {
		logx.Error("get redis error:", err)
		return err
	}
	session, err := parseRedisSession(out)
	if err != nil {
		logx.Error("parse redis session error:", err)
		return err
	}
	//logx.Infof("get redis hash session:%+v", session)
	token := tools.Bytes2String(session.Token[:])
	//serverIp := tools.Bytes2String(session.ServerIp[:])
	//logx.Info("get token:", token, ", serverIP:", tools.RMu0000(serverIp))
	if token != r.Header.Get("Token") {
		logx.Error("parse redis session error:", "web token:", r.Header.Get("Token"), ",cache token:", token)
		return err
	}
	return nil
}
