package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"context"
	"encoding/csv"
	"github.com/spf13/cast"
	"net/http"
	"path"

	"github.com/zeromicro/go-zero/core/logx"
)

type SzQuoteLevelUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSzQuoteLevelUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SzQuoteLevelUploadLogic {
	return &SzQuoteLevelUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SzQuoteLevelUploadLogic) SzQuoteLevelUpload(r *http.Request) (resp *types.BaseRsp, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("form file error:", err)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}, nil
	}
	fileSuffix := path.Ext(header.Filename)
	l.Logger.Info("get filename:", header.Filename, ", suffix:", fileSuffix)

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.Comma = ',' // 指定分隔符
	reader.LazyQuotes = true
	lines, err := reader.ReadAll()
	if err != nil {
		l.Logger.Error("read file error:", err)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}, nil
	}
	l.Logger.Info("get record len:", len(lines))
	l.ProcessDBQuoteFile(lines)

	return &types.BaseRsp{
		Code: 200,
		Msg:  "文件上传成功",
	}, nil

	/*
		if len(lines) > 100 { // 行情的数据量比较大时，先回包，否则页面会等待超时
			threading.GoSafe(func() {
				l.ProcessSzQuoteFile(lines)
			})
			return &types.BaseRsp{
				Code: 200,
				Msg:  "文件已接收，请等待处理结果",
			}, nil
		} else { // 否则处理完后就回包
			l.ProcessSzQuoteFile(lines)
			return &types.BaseRsp{
				Code: 200,
				Msg:  "文件上传成功",
			}, nil
		}
	*/

}

func (l *SzQuoteLevelUploadLogic) ProcessDBQuoteFile(lines [][]string) {
	var quotes []*assessservice.QuoteLevel
	for k, line := range lines {
		if len(line) < 13 {
			l.Logger.Info("format error, line:", line)
			continue
		}
		if k == 0 { //跳过页眉
			continue
		}
		//l.Logger.Info("get line:", line)
		//l.Logger.Info("get line[4]", line[4])
		//l.Logger.Info("get line[5]", line[5])
		quote := &assessservice.QuoteLevel{
			//Id:            line[0],
			SeculityId:    line[1],
			OrgiTime:      cast.ToInt64(line[2]),
			LastPrice:     cast.ToInt64(line[3]),
			AskPrice:      line[4],
			AskVol:        line[5],
			BidPrice:      line[6],
			BidVol:        line[7],
			TotalTradeVol: cast.ToInt64(line[8]),
			TotalAskVol:   cast.ToInt64(line[9]),
			TotalBidVol:   cast.ToInt64(line[10]),
			MkVwap:        0,
		}
		quotes = append(quotes, quote)
		l.Logger.Infof("get quotes:%+v", quotes)
	}
	_, err := l.svcCtx.AssessClient.PushSzQuoteLevel(l.ctx, &assessservice.ReqPushSzLevel{
		Quote: quotes,
	})
	if err != nil {
		l.Logger.Error("rpc PushSzQuoteLevel error:", err)
	}
}

func (l *SzQuoteLevelUploadLogic) ProcessSzQuoteFile(lines [][]string) {
	var quotes []*assessservice.QuoteLevel
	for i, line := range lines {
		//l.Logger.Info("get line:", line)
		if len(line) < 13 {
			l.Logger.Info("format error, line:", line)
			continue
		}
		quote := &assessservice.QuoteLevel{
			//Id:            line[0],
			SeculityId:    line[0],
			OrgiTime:      cast.ToInt64(line[2]),
			LastPrice:     cast.ToInt64(line[3]),
			AskPrice:      line[4],
			AskVol:        line[5],
			BidPrice:      line[6],
			BidVol:        line[7],
			TotalTradeVol: cast.ToInt64(line[9]),
			TotalAskVol:   0,
			TotalBidVol:   0,
			MkVwap:        0,
		}
		if len(quotes) < 100 {
			quotes = append(quotes, quote)
		} else {
			quotes = make([]*assessservice.QuoteLevel, 0, 100)
		}
		// 限制100条数据才发rpc 请求
		if len(quotes) == 100 || i == len(lines)-1 {
			_, err := l.svcCtx.AssessClient.PushSzQuoteLevel(l.ctx, &assessservice.ReqPushSzLevel{
				Quote: quotes,
			})
			if err != nil {
				l.Logger.Error("rpc PushSzQuoteLevel error:", err)
			}
		}
	}
}
