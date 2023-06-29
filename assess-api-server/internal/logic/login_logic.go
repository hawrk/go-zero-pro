package logic

import (
	"algo_assess/global"
	"algo_assess/mornano-rpc-server/mornanoservice"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpc"
	"github.com/zeromicro/go-zero/rest/token"
	"net/http"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(r *http.Request, req *types.LoginReq) (resp *types.LoginRsp, err error) {
	l.Logger.Infof("Login, get req:%+v", req)
	tokenKey := fmt.Sprintf("%s:%s", global.UserTokenKey, req.UserName)
	// 1. 到总线查询用户名和密码
	//var busAllow int
	lRsp, err := l.svcCtx.MornanoClient.LoginCheck(l.ctx, &mornanoservice.LoginReq{
		LoginName: req.UserName,
		Password:  req.Passwd,
	})
	if err != nil {
		l.Logger.Error("MornanoClient LoginCheck rpc error:", err)
		//busAllow = 2
		// 总线校验失败时，再校验一下绩效的
	}
	//2. 到用户权限表查询用户信息
	allow, userType, firstLogin := l.CheckAuthUser(req.UserName, req.Passwd, lRsp.GetPasswd(), lRsp.GetUserName(), lRsp.GetAllow())

	if allow == 0 {
		l.Logger.Error(" auth user and mornano both not found")
		return l.WrapperErrorRsp(300, "用户名/密码校验失败", 0), nil
	}
	// 校验一下密码校验规则
	if allow == 1 {
		if firstLogin == 0 {
			return l.WrapperErrorRsp(300, "首次登陆,请重新设置密码", 2), nil
		} else if firstLogin == 2 {
			return l.WrapperErrorRsp(300, "密码已过期,请重新设置密码", 2), nil
		}
	}
	// 到这里，应该是用户名密码校验通过了
	// 登陆态token校验
	token := r.Header.Get("Authorization")
	l.Logger.Info("get http head token:", token)
	var expire int64
	var accessToken string
	// TODO: 这里校验Token有点别扭，找个时间梳理一下
	// TODO： 因为既然已经校验了用户密码了，就没必要校验token了，但必须返回其他信息，还是得去权限服务查回来
	if token == "" || len(token) == 0 {
		// 3. 用户名密码验证通过,产生token
		accessToken, err = l.GenerateToken(req)
		if err != nil {
			l.Logger.Error("generate token error:", err)
			return l.WrapperErrorRsp(300, err.Error(), 0), nil
		}
		l.Logger.Info("generate token:", accessToken)
		// 3. token 和登陆状态信息写入redis
		m := make(map[string]string)
		m["token"] = accessToken
		m["login_status"] = "1" // 登录状态为1，退出登录为0
		_, err = l.svcCtx.HRedisClient.HSet(l.ctx, tokenKey, m).Result()
		//	time.Second*time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)).Result()
		if err != nil {
			l.Logger.Error("set redis key error:", err)
		}
		// 设置超时时间
		_, err = l.svcCtx.HRedisClient.Expire(l.ctx, tokenKey, time.Second*time.Duration(l.svcCtx.Config.JwtAuth.AccessExpire)).Result()
		if err != nil {
			l.Logger.Error("expire ", tokenKey, " error:", err)
		}
		expire = time.Now().Unix() + l.svcCtx.Config.JwtAuth.AccessExpire
	} else {
		// 未登陆，校验token
		_, err = l.CheckToken(r)
		if err != nil {
			return l.WrapperErrorRsp(300, err.Error(), 0), nil
		}

		// 4.用户名密码校验通过后，登陆状态改为1，已登陆
		_, err = l.svcCtx.HRedisClient.HSet(l.ctx, tokenKey, "login_status", "1").Result()
		if err != nil {
			return l.WrapperErrorRsp(209, err.Error(), 0), nil
		}
	}

	// 记录一下登陆次数
	l.SetLoginCount(req.UserName)

	return &types.LoginRsp{
		Code:     200,
		Msg:      "success",
		Allow:    1,
		Role:     int(lRsp.GetRole()),
		Token:    accessToken,
		Expire:   expire,
		UserType: userType,
	}, nil
}

func (l *LoginLogic) GenerateToken(req *types.LoginReq) (string, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessSecret := l.svcCtx.Config.JwtAuth.AccessSecret

	claims := make(jwt.MapClaims)
	claims["exp"] = now + accessExpire
	claims["iat"] = now
	claims[global.JwtUserId] = req.UserName
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(accessSecret))
}

func (l *LoginLogic) CheckToken(r *http.Request) (bool, error) {
	if l.svcCtx.Config.JwtAuth.TokenCheck {
		l.Logger.Info("token check....")
		parser := token.NewTokenParser()
		accessSecret := l.svcCtx.Config.JwtAuth.AccessSecret
		tok, err := parser.ParseToken(r, accessSecret, "")
		if err != nil {
			l.Logger.Error("parser token error :", err)
			return false, err
		}
		if tok.Valid {
			claims, ok := tok.Claims.(jwt.MapClaims) // 解析token中对内容
			if ok {
				userName, _ := claims[global.JwtUserId]
				l.Logger.Info("get username:", userName, ",token:", tok.Raw)
				// 与redis中的token比较
				tokenKey := fmt.Sprintf("%s:%s", global.UserTokenKey, userName)
				redisToken, err := l.svcCtx.HRedisClient.HGet(l.ctx, tokenKey, "token").Result()
				if err != nil {
					l.Logger.Error("get redis token key error:", err)
					return false, err
				}
				if redisToken != tok.Raw {
					l.Logger.Error("redis token not match request token:", redisToken, tok.Raw)
					return false, errors.New("redis token not match request token")
				}
			}
		}
		return true, nil
	}
	return true, nil
}

func (l *LoginLogic) WrapperErrorRsp(code int, msg string, allow int) *types.LoginRsp {
	return &types.LoginRsp{
		Code:  code,
		Msg:   msg,
		Allow: allow,
	}
}

func (l *LoginLogic) SetLoginCount(userId string) {
	LoginKey := global.LoginCount + ":" + time.Now().Format(global.TimeFormatDay) + ":" + userId
	l.Logger.Info("get LoginKey:", LoginKey)
	_, err := l.svcCtx.HRedisClient.Incr(l.ctx, LoginKey).Result()
	if err != nil {
		l.Logger.Error("incr login count error:", err)
	}
	_, _ = l.svcCtx.HRedisClient.Expire(l.ctx, LoginKey, time.Second*86400).Result() // 设置过期时间为一天
}

// CheckAuthUser 调用用户权限的api接口校验登陆账户
func (l *LoginLogic) CheckAuthUser(userId, passwd string, busPasswd, busUserName string, match int32) (int, int, int) {
	type AuthLoginRequest struct {
		Header      string `header:"X-Header"`
		UserId      string `json:"user_id"`
		Password    string `json:"password"`
		ChanType    int    `json:"chan_type"`
		BusPasswd   string `json:"bus_passwd"`
		BusUserName string `json:"bus_user_name"`
		BusAllow    int    `json:"bus_allow"`
	}
	type AuthLoginRsp struct {
		Code       int    `json:"code"`
		Msg        string `json:"msg"`
		Allow      int    `json:"allow"` // 1-允许登陆， 其他不允许
		UserType   int    `json:"user_type"`
		FirstLogin int    `json:"first_login"` // 是否首次登陆， 0-首次登陆 1-非首次登陆 2-密码过期
	}
	//
	authReq := AuthLoginRequest{
		Header:      "foo-header",
		UserId:      userId,
		Password:    passwd,
		ChanType:    1, // 绩效渠道
		BusPasswd:   busPasswd,
		BusUserName: busUserName,
		BusAllow:    int(match),
	}
	url := l.svcCtx.Config.AccountAuth.UrlPrefix + "/algo-assess/v1/auth/login"
	resp, err := httpc.Do(l.ctx, http.MethodPost, url, authReq)
	if err != nil {
		l.Logger.Error("httpc request error:", err)
		return 0, 0, 0
	}
	//l.Logger.Infof("get resp:%+v", resp)
	var aRsp AuthLoginRsp
	httpc.Parse(resp, &aRsp)
	l.Logger.Infof("get resp body:%+v", aRsp)

	return aRsp.Allow, aRsp.UserType, aRsp.FirstLogin
}
