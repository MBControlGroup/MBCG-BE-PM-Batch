package controllers

import (
	"dao"
	"models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

var headRowNames = [...]string{
	"姓名",
	"身份证号",
	"手机号码",
	"军衔",
	"微信openid",
	"直属指挥官姓名",
	"所属市",
	"所属区",
	"所属街道",
}

const colCount = len(headRowNames)

// 根据民兵数据生成xlsx文件，返回文件名
func getFileForSoldiers(exportSoldiers []*models.ExportSoldier) string {
	// 新建文件
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	// 添加表头
	headRow := sheet.AddRow()
	for col := 0; col < colCount; col++ {
		cell := headRow.AddCell()
		cell.Value = headRowNames[col]
	}

	// 将数据插入表格中
	for _, exportSoldier := range exportSoldiers {
		row := sheet.AddRow()
		var cells []*xlsx.Cell
		for col := 0; col < colCount; col++ {
			cells = append(cells, row.AddCell())
		}
		cells[0].Value = exportSoldier.Name
		cells[1].Value = exportSoldier.IDNum
		cells[2].Value = strconv.FormatUint(exportSoldier.PhoneNum, 10)
		cells[3].Value = exportSoldier.Rank
		cells[4].Value = exportSoldier.WechatOpenid
		cells[5].Value = exportSoldier.CommanderName
		cells[6].Value = exportSoldier.City
		cells[7].Value = exportSoldier.District
		cells[8].Value = exportSoldier.Street
	}

	// 将文件写入硬盘
	filename := GetFilenameWithTimestamp()
	err = file.Save(filename)
	if err != nil {
		panic(err)
	}

	return filename
}

// Download download batch
func Download(ctx *gin.Context) {
	defer CrashHandler(ctx)

	// 获取DB中所有的民兵数据（经过导出处理）
	exportSoldiers := dao.GetExportSoldiers()

	// 生成文件
	filename := getFileForSoldiers(exportSoldiers)

	// 返回下载文件
	ctx.Header("Content-Disposition", "attachment; filename=\"民兵数据.xlsx\"")
	ctx.File(filename)
}
