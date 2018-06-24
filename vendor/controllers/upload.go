package controllers

import (
	"dao"
	"fmt"
	"log"
	"models"
	"net/http"
	"path"
	"protocol"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

// 检查文件类型，目前只支持上传xlsx
func validateFileType(filename string) bool {
	// 根据后缀名判断文件类型
	fileSuffix := path.Ext(filename)
	return fileSuffix == ".xlsx"
}

// 报告前端文件类型错误
func reportWrongType(ctx *gin.Context) {
	res := &protocol.ResponseMsg{
		Code:  http.StatusUnsupportedMediaType,
		EnMsg: "wrong type",
		CnMsg: "文件类型不支持",
		Data:  nil,
	}
	ctx.JSON(res.Code, res)
}

// 报告前端表格数据错误
func reportTableDataErr(ctx *gin.Context, errMsg string) {
	res := &protocol.ResponseMsg{
		Code:  http.StatusBadRequest,
		EnMsg: "table data error",
		CnMsg: errMsg,
		Data:  nil,
	}
	ctx.JSON(res.Code, res)
}

// 报告前端导入成功
func reportSuccess(ctx *gin.Context) {
	res := &protocol.ResponseMsg{
		Code:  http.StatusCreated,
		EnMsg: "import ok",
		CnMsg: "导入成功",
		Data:  nil,
	}
	ctx.JSON(res.Code, res)
}

var fieldsInTable = [...]string{
	"姓名",
	"手机号码",
	"身份证号码",
	"所属市",
	"所属区",
	"所属街道",
	"直属指挥官姓名",
}

const filedsCount = len(fieldsInTable)

// 将表格中的一行转成一个RawSoldier
func rowToRawSoldier(row *xlsx.Row) (res *models.RawSoldier, errMsg string) {
	var fields [filedsCount]string

	// 各列判空
	errCol := -1
	for col := 0; col < filedsCount; col++ {
		fields[col] = row.Cells[col].String()
		if fields[col] == "" || fields[col] == " " {
			errCol = col
			break
		}
	}

	if errCol == 0 {
		res = nil
		errMsg = ""
		return
	} else if errCol > 0 {
		res = nil
		errMsg = fmt.Sprintf("%s为空", fieldsInTable[errCol])
		return
	}

	// 判断手机号是否符合格式
	phoneNum, err := strconv.ParseUint(fields[1], 10, 64)
	if err != nil {
		res = nil
		errMsg = "手机号码数据格式错误"
		return
	}

	res = &models.RawSoldier{
		Name:          fields[0],
		PhoneNum:      phoneNum,
		IDNum:         fields[2],
		City:          fields[3],
		District:      fields[4],
		Street:        fields[5],
		CommanderName: fields[6],
	}
	errMsg = ""
	return
}

func getErrMsgInTable(rowIdx int, cellErrMsg string) string {
	return fmt.Sprintf("第%d行, %s", rowIdx+1, cellErrMsg)
}

// 取出上传表格中的士兵数据
func getRawSoldiers(filename string) (res []*models.RawSoldier, isClientErr bool, errMsg string) {
	// 读取xlsx文件
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		res = nil
		isClientErr = false
		errMsg = err.Error()
		return
	}

	// 检查表格中是否有数据
	sheetCount := len(xlFile.Sheets)
	if sheetCount == 0 {
		res = nil
		isClientErr = true
		errMsg = "文件中没有数据"
		return
	}

	// 读取表格数据
	sheet := xlFile.Sheets[0] // 目前仅支持一张表
	for rowIdx, row := range sheet.Rows {
		rawSoldier, cellErrMsg := rowToRawSoldier(row)
		if rawSoldier != nil {
			res = append(res, rawSoldier)
		} else if cellErrMsg == "" { // 到达空行，本sheet读取结束
			break
		} else {
			res = nil
			isClientErr = true
			errMsg = getErrMsgInTable(rowIdx, cellErrMsg)
			return
		}
	}

	isClientErr = false
	errMsg = ""
	return
}

// Upload upload batch
func Upload(ctx *gin.Context) {
	defer CrashHandler(ctx)

	// 获取上传的文件
	file, err := ctx.FormFile("upload_batch")
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusInternalServerError, "error when getting file\n")
		return
	}
	log.Println("got file:", file.Filename)

	if !validateFileType(file.Filename) {
		reportWrongType(ctx)
		return
	}

	// 存到本地
	newFilename := GetFilenameWithTimestamp()
	ctx.SaveUploadedFile(file, newFilename)

	// 读取数据
	rawSoldiers, isClientErr, errMsg := getRawSoldiers(newFilename)
	if rawSoldiers == nil {
		if isClientErr {
			reportTableDataErr(ctx, errMsg)
			return
		}
		panic(errMsg)
	}

	// 数据入库
	clientErrMsg := dao.InsertRawSoldiers(rawSoldiers)
	if clientErrMsg != "" {
		reportTableDataErr(ctx, clientErrMsg)
	} else {
		reportSuccess(ctx)
	}
}
