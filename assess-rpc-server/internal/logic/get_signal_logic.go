package logic

import (
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSignalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSignalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSignalLogic {
	return &GetSignalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetSignal 高阶分析：信号分析
func (l *GetSignalLogic) GetSignal(in *proto.SignalReq) (*proto.SignalRsp, error) {
	l.Logger.Infof("in GetSignal,get req:%+v", in)
	// 直接查summary表
	out, err := l.svcCtx.SummaryRepo.GetSignalSummary(l.ctx, in.GetStartDay(), in.GetEndDay(),
		in.GetUserId(), in.GetUserType(), int(in.GetAlgoId()))
	if err != nil {
		l.Logger.Error("GetSignalSummary error:", err)
		return &proto.SignalRsp{
			Code: 370,
			Msg:  err.Error(),
			Info: nil,
		}, nil
	}
	if len(out) == 0 {
		l.Logger.Info("win ratio no data, return....")
		return &proto.SignalRsp{
			Code: 371,
			Msg:  "success",
			Info: nil,
		}, nil
	}

	var arr []*proto.SignalInfo
	strStartDay := cast.ToString(in.GetStartDay())
	y := cast.ToInt(strStartDay[:4])
	m := cast.ToInt(strStartDay[4:6])
	mk := time.Month(m)
	d := cast.ToInt(strStartDay[6:8])
	for i := 0; i <= 360; i++ {
		startDay := time.Date(y, mk, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, i).Format(global.TimeFormatDay)
		iDay := cast.ToInt(startDay)
		for _, v := range out {
			var a *proto.SignalInfo
			if cast.ToInt64(startDay) > in.GetEndDay() {
				break
			}
			if iDay < v.Date { // 当前日期在DB中无数据，填充0
				a = &proto.SignalInfo{
					Day:      TransDay(iDay), // 20220410 转成 2022.04.10 格式
					OrderNum: 0.00,
					Progress: 0.00,
				}
				arr = append(arr, a)
				break
			} else if iDay == v.Date {
				if v.EntrustQty == 0 {
					continue
				}
				a = &proto.SignalInfo{
					Day:      TransDay(v.Date),
					OrderNum: int32(v.OrderNum),
					Progress: float64(v.DealQty) / float64(v.EntrustQty) * 100,
				}
				arr = append(arr, a)
				break
			}
		}
	}

	return &proto.SignalRsp{
		Code: 200,
		Msg:  "success",
		Info: arr,
	}, nil
}
