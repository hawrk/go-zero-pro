package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/pkg/tools"
	"context"
	"encoding/csv"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"path"
	"strings"
)

type ShQuoteLevelUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShQuoteLevelUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShQuoteLevelUploadLogic {
	return &ShQuoteLevelUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ShQuoteLevelUpload 沪市行情数据修复文件上传
func (l *ShQuoteLevelUploadLogic) ShQuoteLevelUpload(r *http.Request) (resp *types.BaseRsp, err error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("form file error:", err)
		resp = &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}
	}
	fileSuffix := path.Ext(header.Filename)
	l.Logger.Info("get filename:", header.Filename, ", suffix:", fileSuffix)
	// 不带后缀名
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.Comma = '|' // 指定分隔符
	reader.LazyQuotes = true
	lines, err := reader.ReadAll()
	if err != nil {
		l.Logger.Error("read file error:", err)
		return &types.BaseRsp{
			Code: 10000,
			Msg:  "文件解析失败",
		}, nil
	}
	//today := time.Now().Format(global.TimeFormatDay)
	// SecID|pTimestamp|OrigTime|LastPrice|AskPrice|AskVol|BidPrice|BidVol|TotalTradeNum|TotalTradeVol|TotalTradeValue|TradeStatus|MsgID
	l.Logger.Info("get record len:", len(lines))
	if len(lines) > 100 { // 行情的数据量比较大时，先回包，否则页面会等待超时
		threading.GoSafe(func() {
			l.ProcessQuoteFile(lines)
		})
		return &types.BaseRsp{
			Code: 200,
			Msg:  "文件已接收，请等待处理结果",
		}, nil
	} else {
		l.ProcessQuoteFile(lines)
		return &types.BaseRsp{
			Code: 200,
			Msg:  "文件上传成功",
		}, nil
	}
}

// GetOriginTime 行情源文件时间格式为HHMMSSssss  1410060272(超过10点为10位， 10点之前为9位)
func GetOriginTime(in string, date string) int64 {
	var out string
	if len(in) == 9 {
		out = date + "0" + in[:3]
	} else if len(in) == 10 {
		out = date + in[:4]
	}
	t := tools.TimeMoveForward(out)
	return t
}

// GetPriceStr 把原*1000000的字符串转换成真实价格的字符串
func GetPriceStr(in string) string {
	arr := strings.Split(in, ",")
	var build strings.Builder
	for _, v := range arr {
		k := tools.DivTenThousand(cast.ToUint32(v))
		build.WriteString(k)
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// GetVolStr 把原*100 申买量转成真实申买量的字符串
func GetVolStr(in string) string {
	arr := strings.Split(in, ",")
	var build strings.Builder
	for _, v := range arr {
		k := tools.DivHundred(cast.ToUint64(v))
		build.WriteString(k)
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

func (l *ShQuoteLevelUploadLogic) ProcessQuoteFile(lines [][]string) {
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
			_, err := l.svcCtx.AssessClient.PushShQuoteLevel(l.ctx, &assessservice.ReqPushShLevel{
				Quote: quotes,
			})
			if err != nil {
				l.Logger.Error("rpc PushShQuoteLevel error:", err)
			}
		}
	}
}
