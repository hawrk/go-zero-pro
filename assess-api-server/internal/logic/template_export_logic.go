package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/global"
	"bytes"
	"context"
	"encoding/csv"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type TemplateExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTemplateExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TemplateExportLogic {
	return &TemplateExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// TemplateExport 数据修复模板导出
func (l *TemplateExportLogic) TemplateExport(req *types.TemplateExpReq, w http.ResponseWriter) error {
	// todo: add your logic here and delete this line
	l.Logger.Infof("get req:%+v", *req)
	var fileName string
	var content []string
	if req.ExportType == 1 { // 行情模板
		return nil
	} else if req.ExportType == 2 { // 母单模板
		fileName, content = l.BuildAlgoOrderCsv()
	} else if req.ExportType == 3 { // 子单模板
		fileName, content = l.BuildChildOrderCsv()
	} else {
		l.Logger.Error("unsupported type:", req.ExportType)
		return nil
	}

	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "filename="+fileName)
	buff := new(bytes.Buffer)
	buff.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(buff)
	wStr.Write(content)
	wStr.Flush()
	_, _ = w.Write(buff.Bytes())

	return nil
}

// BuildQuoteCsv 行情模板数据
func (l *TemplateExportLogic) BuildQuoteCsv() []string {

	return nil
}

// BuildAlgoOrderCsv 母单模板数据
func (l *TemplateExportLogic) BuildAlgoOrderCsv() (string, []string) {
	return "algo_order.csv", global.AlgoHeader
}

func (l *TemplateExportLogic) BuildChildOrderCsv() (string, []string) {
	return "order.csv", global.ChildOrderHeader
}
