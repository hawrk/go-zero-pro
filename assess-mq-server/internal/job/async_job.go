// Package job
/*
 Author: hawrkchen
 Date: 2023/5/18 9:43
 Desc:
*/
package job

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/global"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"strings"
)

func StartAsyncJob(c config.Config, svcContext *svc.ServiceContext) {
	logx.Info("into StartAsyncJob....")

	// 母单写入
	threading.GoSafe(func() {
		for o := range global.OrderChan {
			if err := PersistentOrder(c, svcContext, &o); err != nil {
				logx.Error("PersistentOrder error:", err)
			}
		}
	})

	// 子单异步写入
	threading.GoSafe(func() {
		for od := range global.OrderDetailChan {
			if err := PersistentOrderDetail(c, svcContext, &od); err != nil {
				logx.Error("PersistentOrderDetail error:", err)
			}
		}
	})

	// 画像明细写入
	threading.GoSafe(func() {
		for pf := range global.ProfileChan {
			if err := PersistentProfile(c, svcContext, &pf); err != nil {
				logx.Error("PersistentProfile error:", err)
			}
		}
	})

	// 算法汇总写入
	threading.GoSafe(func() {
		for ps := range global.ProfileSumChan {
			if err := PersistentProfileSum(c, svcContext, &ps); err != nil {
				logx.Error("PersistentProfileSum error:", err)
			}
		}
	})

	// 时间线写入
	threading.GoSafe(func() {
		for tl := range global.TlProfileSumChan {
			if err := PersistentTimeLine(c, svcContext, &tl); err != nil {
				logx.Error("PersistentTimeLine error:", err)
			}
		}
	})

	/*
		for {
			select {
			case profile := <- global.ProfileChan:
				PersistentProfile(c, svcContext, profile)
			case profileSum := <- global.ProfileSumChan:
				PersistentProfileSum(c, svcContext, profileSum)
			case tlProfileSum := <- global.TlProfileSumChan:
					PersistentTimeLine(c,svcContext, tlProfileSum)
			default:
				time.Sleep(time.Millisecond*100)
			}
		}
	*/

}

func PersistentOrder(c config.Config, svcContext *svc.ServiceContext, in *global.MAlgoOrder) error {
	ctx := context.Background()
	if err := svcContext.AlgoOrderRepo.CreateAlgoOrder(ctx, in); err != nil {
		return err
	}
	return nil
}

func PersistentOrderDetail(c config.Config, svcContext *svc.ServiceContext, in *global.ChildOrderData) error {
	ctx := context.Background()
	if err := svcContext.OrderDetailRepo.CreateOrderDetail(ctx, in); err != nil {
		return err
	}
	return nil
}

func PersistentProfile(c config.Config, svcContext *svc.ServiceContext, in *global.Profile) error {
	//logx.Infof("get pf %+v:", *in)
	ctx := context.Background()
	// 落算法画像表
	if strings.HasPrefix(in.SourcePrx, global.OrderSourceOri) || strings.HasPrefix(in.SourcePrx, global.OrderSourceImp) {
		if in.IndexCount > 1 {
			if err := svcContext.ProfileOrigRepo.UpdateAlgoProfile(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.ProfileOrigRepo.CreateAlgoProfile(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceNor) {
		if in.IndexCount > 1 {
			if err := svcContext.ProfileRepo.UpdateAlgoProfile(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.ProfileRepo.CreateAlgoProfile(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceFix) {
		if in.IndexCount > 2 {
			if err := svcContext.ProfileRepo.UpdateAlgoProfile(ctx, in); err != nil {
				return err
			}
		} else { // 数据修复为1时，不能判断DB中是否有数据
			record, err := svcContext.ProfileRepo.GetDataByProfileKey(ctx, in.Date, in.AccountId, in.AlgoId, in.SecId, in.AlgoOrderId)
			if err != nil {
				return err
			}
			if record.AccountId != "" { //有值
				if err := svcContext.ProfileRepo.UpdateAlgoProfile(ctx, in); err != nil {
					return err
				}
			} else {
				if err := svcContext.ProfileRepo.CreateAlgoProfile(ctx, in); err != nil {
					return err
				}
			}
		}
	} else {
		return errors.New("unknown order source")
	}
	return nil
}

func PersistentProfileSum(c config.Config, svcContext *svc.ServiceContext, in *global.ProfileSum) error {
	//logx.Infof("get ps:%+v", *in)
	ctx := context.Background()
	// 落地算法汇总表
	if strings.HasPrefix(in.SourcePrx, global.OrderSourceOri) || strings.HasPrefix(in.SourcePrx, global.OrderSourceImp) {
		if in.IndexCount > 1 {
			if err := svcContext.SummaryOrigRepo.UpdateAlgoSummary(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.SummaryOrigRepo.CreateAlgoSummary(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceNor) {
		if in.IndexCount > 1 {
			if err := svcContext.SummaryRepo.UpdateAlgoSummary(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.SummaryRepo.CreateAlgoSummary(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceFix) {
		if in.IndexCount > 2 {
			if err := svcContext.SummaryRepo.UpdateAlgoSummary(ctx, in); err != nil {
				return nil
			}
		} else { // 先查询
			record, err := svcContext.SummaryRepo.GetDataBySummaryKey(ctx, in.Date, in.AccountId, in.AlgoId)
			if err != nil {
				return err
			}
			if record.UserId != "" {
				if err := svcContext.SummaryRepo.UpdateAlgoSummary(ctx, in); err != nil {
					return nil
				}
			} else {
				if err := svcContext.SummaryRepo.CreateAlgoSummary(ctx, in); err != nil {
					return err
				}
			}
		}
	} else {
		return errors.New("unknown order source")
	}

	return nil
}

func PersistentTimeLine(c config.Config, svcContext *svc.ServiceContext, in *global.ProfileSum) error {
	//logx.Infof("get tl:%+v", *in)
	ctx := context.Background()
	// 落地时间线表
	if strings.HasPrefix(in.SourcePrx, global.OrderSourceOri) || strings.HasPrefix(in.SourcePrx, global.OrderSourceImp) {
		if in.IndexCount > 1 {
			if err := svcContext.TimeLineOrigRepo.UpdateAlgoTimeLine(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.TimeLineOrigRepo.CreateAlgoTimeLine(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceNor) {
		if in.IndexCount > 1 {
			if err := svcContext.TimeLineRepo.UpdateAlgoTimeLine(ctx, in); err != nil {
				return err
			}
		} else {
			if err := svcContext.TimeLineRepo.CreateAlgoTimeLine(ctx, in); err != nil {
				return err
			}
		}
	} else if strings.HasPrefix(in.SourcePrx, global.OrderSourceFix) {
		if in.IndexCount > 2 {
			if err := svcContext.TimeLineRepo.UpdateAlgoTimeLine(ctx, in); err != nil {
				return err
			}
		} else { // 先查询
			record, err := svcContext.TimeLineRepo.GetDataByTimeLineKey(ctx, in.AccountId, in.TransAt, in.AlgoId)
			if err != nil {
				return err
			}
			if record.AccountId != "" {
				if err := svcContext.TimeLineRepo.UpdateAlgoTimeLine(ctx, in); err != nil {
					return err
				}
			} else {
				if err := svcContext.TimeLineRepo.CreateAlgoTimeLine(ctx, in); err != nil {
					return err
				}
			}
		}
	} else {
		return errors.New("unknown order source")
	}
	return nil
}
