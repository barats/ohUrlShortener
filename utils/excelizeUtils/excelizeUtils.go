package excelizeutils

import (
	"errors"
	"ohurlshortener/core"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func AccessLogToExcel(logs []core.AccessLog) ([]byte, error) {
	if logs == nil {
		return nil, errors.New("数据为空")
	}
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	//填充表头
	f.SetCellValue("Sheet1", "A1", "短链接")
	f.SetCellValue("Sheet1", "B1", "访问时间")
	f.SetCellValue("Sheet1", "C1", "访问IP")
	f.SetCellValue("Sheet1", "D1", "UserAgent")
	f.SetColWidth("Sheet1", "B", "C", 15)
	for i := 0; i < len(logs); i++ {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), logs[i].ShortUrl)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), logs[i].AccessTime)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), logs[i].Ip.String)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), logs[i].UserAgent.String)
	}
	f.SetActiveSheet(index)
	if excellBytes, erorrW := f.WriteToBuffer(); erorrW != nil {
		return nil, erorrW
	} else {
		return excellBytes.Bytes(), nil
	}

}
