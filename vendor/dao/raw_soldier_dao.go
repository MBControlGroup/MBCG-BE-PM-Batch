package dao

import (
	"fmt"
	"models"

	"github.com/jinzhu/gorm"
)

func getCommanderIDByCommanderName(transactionDB *gorm.DB, commanderName string) (commanderID int32, err error) {
	commanderID = -1

	commander := models.Soldier{}
	err = transactionDB.Where("name = ?", commanderName).First(&commander).Error
	if err != nil {
		return
	}
	commanderID = commander.SoldierID
	return
}

func getServeOfficeIDByOfficeName(transactionDB *gorm.DB, cityName, districtName, streetName string) (serveOfficeID int32, err error) {
	serveOfficeID = -1

	cityOffice := models.Office{}
	err = transactionDB.Where("name = ?", cityName).First(&cityOffice).Error
	if err != nil {
		return
	}

	districtOffice := models.Office{}
	err = transactionDB.Where(&models.Office{
		HigherOfficeID: cityOffice.OfficeID,
		Name:           districtName,
	}).First(&districtOffice).Error
	if err != nil {
		return
	}

	streetOffice := models.Office{}
	err = transactionDB.Where(&models.Office{
		HigherOfficeID: districtOffice.OfficeID,
		Name:           streetName,
	}).First(&streetOffice).Error
	if err != nil {
		return
	}

	serveOfficeID = streetOffice.OfficeID
	return
}

// 将RawSoldier替换成Soldier
func rawSoldierToSoldier(transactionDB *gorm.DB, rawSoldier *models.RawSoldier) (realSoldier *models.Soldier, errMsg string, dbErr error) {
	// 检查字段是否合法
	// 获取民兵直属指挥官ID
	commanderID, dbErr := getCommanderIDByCommanderName(transactionDB, rawSoldier.CommanderName)
	if dbErr != nil {
		realSoldier = nil
		if dbErr.Error() == RecordNotFoundErrMsg {
			errMsg = "民兵直属指挥官未在数据库中"
			dbErr = nil
		}
		return
	}

	// 获取民兵所属机关ID
	serveOfficeID, dbErr := getServeOfficeIDByOfficeName(transactionDB, rawSoldier.City, rawSoldier.District, rawSoldier.Street)
	if dbErr != nil {
		realSoldier = nil
		if dbErr.Error() == RecordNotFoundErrMsg {
			errMsg = "民兵所属机关未在数据库中"
			dbErr = nil
		}
		return
	}

	realSoldier = models.NewSoldier()
	realSoldier.IDNum = rawSoldier.IDNum
	realSoldier.Name = rawSoldier.Name
	realSoldier.PhoneNum = rawSoldier.PhoneNum
	realSoldier.CommanderID = commanderID
	realSoldier.ServeOfficeID = serveOfficeID
	return
}

// InsertRawSoldiers 批量导入民兵数据
func InsertRawSoldiers(rawSoldiers []*models.RawSoldier) (clientErrMsg string) {
	defer func() {
		// 重启自动提交
		MyDB.Exec("SET AUTOCOMMIT=1")
	}()

	// 检查连接
	MyDB.DB().Ping()

	// 关闭自动提交
	if dbErr := MyDB.Exec("SET AUTOCOMMIT=0").Error; dbErr != nil {
		panic(dbErr)
	}

	// 开启批量导入事务
	transactionDB := MyDB.Begin()
	for rowIdx, rawSoldier := range rawSoldiers {
		// 根据身份证号码查重
		existSoldier := models.Soldier{}
		if dbErr := transactionDB.Where("id_num = ?", rawSoldier.IDNum).First(&existSoldier).Error; dbErr == nil {
			clientErrMsg = fmt.Sprintf("第%d行, 该民兵已存在", rowIdx+1)
			transactionDB.Rollback()
			return
		} else if dbErr.Error() != RecordNotFoundErrMsg {
			transactionDB.Rollback()
			panic(dbErr)
		}

		// 将原始民兵数据转成实际民兵条目数据
		realSoldier, errMsg, dbErr := rawSoldierToSoldier(transactionDB, rawSoldier)
		if realSoldier == nil {
			if errMsg != "" {
				clientErrMsg = fmt.Sprintf("第%d行, %s", rowIdx+1, errMsg)
				transactionDB.Rollback()
				return clientErrMsg
			}
			if dbErr != nil {
				transactionDB.Rollback()
				panic(dbErr)
			}
		}

		// 执行插入
		dbErr = transactionDB.Create(realSoldier).Error
		if dbErr != nil {
			transactionDB.Rollback()
			panic(dbErr)
		}
	}
	transactionDB.Commit()
	return
}
