package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"strconv"
)

type DownloadAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadAlgoOptimizeBaseLogic {
	return &DownloadAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadAlgoOptimizeBaseLogic) DownloadAlgoOptimizeBase(w http.ResponseWriter) error {
	// todo: add your logic here and delete this line
	r, err := l.svcCtx.AssessClient.DownloadOptimizeBase(l.ctx, &assessservice.DownloadOptimizeBaseReq{})
	if err != nil {
		return err
	}
	optimizeBases := r.GetList()
	//内容先写入buffer缓存
	buff := new(bytes.Buffer)
	//写入UTF-8 BOM,此处如果不写入就会导致写入的汉字乱码
	buff.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(buff)
	wStr.Write([]string{"Id", "ProviderId", "ProviderName", "SecId", "SecName", "AlgoId", "AlgoType", "AlgoName", "OpenRate", "IncomeRate", "BasisPoint", "CreateTime", "UpdateTime"})
	for _, a := range optimizeBases {
		wStr.Write([]string{
			strconv.FormatInt(a.Id, 10),
			fmt.Sprintf("%d", a.ProviderId),
			a.ProviderName,
			a.SecId,
			a.SecName,
			fmt.Sprintf("%d", a.AlgoId),
			fmt.Sprintf("%d", a.AlgoType),
			a.AlgoName,
			fmt.Sprintf("%f", a.OpenRate),
			fmt.Sprintf("%f", a.IncomeRate),
			fmt.Sprintf("%f", a.BasisPoint),
			a.CreateTime,
			a.UpdateTime},
		)
	}
	wStr.Flush()
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "attachment;filename="+"algoOptimize.csv")
	// 返回[]byte数据
	_, _ = w.Write(buff.Bytes())
	return nil
}

type AlgoOptimizeBases struct {
	AlgoOptimizeBase []types.OptimizeBase `xml:"AlgoOptimizeBase"`
}
