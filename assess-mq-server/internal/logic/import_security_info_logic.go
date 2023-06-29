package logic

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/threading"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportSecurityInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportSecurityInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportSecurityInfoLogic {
	return &ImportSecurityInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ImportSecurityInfo 配置： 证券信息导入
func (l *ImportSecurityInfoLogic) ImportSecurityInfo(in *proto.ImportSecurityReq) (*proto.ImportSecurityRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in ImportSecurityInfo, process")
	if len(in.GetList()) == 0 {
		l.Logger.Error("empty list")
		return &proto.ImportSecurityRsp{
			Code:   300,
			Msg:    errors.New("empty list").Error(),
			Result: 2,
		}, nil
	}

	threading.GoSafe(func() {
		// 缓存一份数据，用来判断导入的数据是插入还是更新
		m := make(map[string]struct{})
		dao.GSecurityMap.RWMutex.RLock()
		for k, _ := range dao.GSecurityMap.SecurityBase {
			m[k] = struct{}{}
		}
		dao.GSecurityMap.RWMutex.RUnlock()

		for _, v := range in.GetList() {
			if _, exist := m[strings.TrimSpace(v.SecId)]; exist {
				if err := l.svcCtx.SecurityRepo.ImportSecurityUpdate(l.ctx, v); err != nil {
					l.Logger.Error("ImportSecurityUpdate error :", err)
					continue
				}
			} else {
				if err := l.svcCtx.SecurityRepo.ImportSecurityCreate(l.ctx, v); err != nil {
					l.Logger.Error("ImportSecurityCreate error:", err)
					continue
				}
			}
		}
	})

	return &proto.ImportSecurityRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
