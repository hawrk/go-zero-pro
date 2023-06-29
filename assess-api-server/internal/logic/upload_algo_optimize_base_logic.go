package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadAlgoOptimizeBaseLogic {
	return &UploadAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadAlgoOptimizeBaseLogic) UploadAlgoOptimizeBase(r *http.Request) (resp *types.OptimizeBaseRsp, err error) {
	file, header, fileError := r.FormFile("file")
	if fileError != nil {
		resp = &types.OptimizeBaseRsp{
			Code: 10000,
			Msg:  "文件上传错误",
		}
		return resp, fileError
	}
	filename := header.Filename
	fmt.Println("上传文件名:", filename)
	resp = &types.OptimizeBaseRsp{
		Code: 0,
		Msg:  "文件上传成功",
	}
	list := make([]*assessservice.OptimizeBase, 0)
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	records, csvErr := reader.ReadAll()
	if csvErr != nil {
		fmt.Println("文件读取错误:", csvErr)
		fmt.Println("文件读取:", records)
		resp = &types.OptimizeBaseRsp{
			Code: 10000,
			Msg:  "文件读取错误",
		}
		return resp, csvErr
	}
	for i, record := range records {
		fmt.Println("文件行数据:", record)
		if i == 0 {
			continue
		}
		var parseErr error
		var Id int64
		var ProviderId int32
		var ProviderName string
		var SecId string
		var SecName string
		var AlgoId int32
		var AlgoType int32
		var AlgoName string
		var OpenRate float64
		var IncomeRate float64
		var BasisPoint float64
		var CreateTime string
		var UpdateTime string
		Id, parseErr = strconv.ParseInt(record[0], 10, 64)
		fmt.Printf("Id:%v\n", Id)
		providerId, parseErr := strconv.ParseInt(record[1], 10, 32)
		ProviderId = int32(providerId)
		fmt.Printf("ProviderId:%v\n", ProviderId)
		ProviderName = record[2]
		fmt.Printf("ProviderName:%v\n", ProviderName)
		SecId = record[3]
		SecName = record[4]
		algoId, parseErr := strconv.ParseInt(record[5], 10, 32)
		AlgoId = int32(algoId)
		algoType, parseErr := strconv.ParseInt(record[6], 10, 32)
		AlgoType = int32(algoType)
		AlgoName = record[7]
		OpenRate, parseErr = strconv.ParseFloat(record[8], 64)
		IncomeRate, parseErr = strconv.ParseFloat(record[9], 64)
		BasisPoint, parseErr = strconv.ParseFloat(record[10], 64)
		CreateTime = record[11]
		UpdateTime = record[12]
		if parseErr != nil {
			resp = &types.OptimizeBaseRsp{
				Code: 10000,
				Msg:  "文件解析错误",
			}
			return resp, parseErr
		}
		optimizeBase := &assessservice.OptimizeBase{
			Id:           Id,
			ProviderId:   ProviderId,
			ProviderName: ProviderName,
			SecId:        SecId,
			SecName:      SecName,
			AlgoId:       AlgoId,
			AlgoType:     AlgoType,
			AlgoName:     AlgoName,
			OpenRate:     OpenRate,
			IncomeRate:   IncomeRate,
			BasisPoint:   BasisPoint,
			CreateTime:   CreateTime,
			UpdateTime:   UpdateTime,
		}
		list = append(list, optimizeBase)
	}
	if len(list) == 0 {
		resp = &types.OptimizeBaseRsp{
			Code: 10000,
			Msg:  "文件解析，合规的数据为0",
		}
	}
	in := &assessservice.UploadOptimizeBaseReq{List: list}
	_, err = l.svcCtx.AssessClient.UploadOptimizeBase(l.ctx, in)
	if err != nil {
		resp = &types.OptimizeBaseRsp{
			Code: 10000,
			Msg:  "文件导入失败",
		}
	}
	return
}
