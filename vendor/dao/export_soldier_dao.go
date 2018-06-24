package dao

import (
	"models"
)

// GetExportSoldiers 获取批量导出的民兵数据
func GetExportSoldiers() []*models.ExportSoldier {
	// 检查连接
	MyDB.DB().Ping()

	// 获取原始民兵数据
	soldiers := GetSoldiers()

	// 修改字段
	var res []*models.ExportSoldier

	for _, soldier := range soldiers {
		city, district, street := GetOfficeNameForSoldier(soldier)
		commander := models.Soldier{}
		err := MyDB.Where("soldier_id = ?", soldier.CommanderID).First(&commander).Error
		if err != nil {
			panic(err)
		}
		var rank string
		if soldier.Rank == "CM" {
			rank = "指挥官"
		} else {
			rank = "民兵"
		}
		exportSoldier := &models.ExportSoldier{
			Name:          soldier.Name,
			IDNum:         soldier.IDNum,
			PhoneNum:      soldier.PhoneNum,
			Rank:          rank,
			WechatOpenid:  soldier.WechatOpenid,
			CommanderName: commander.Name,
			City:          city,
			District:      district,
			Street:        street,
		}
		res = append(res, exportSoldier)
	}

	return res
}
