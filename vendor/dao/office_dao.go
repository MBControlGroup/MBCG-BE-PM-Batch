package dao

import (
	"models"
)

// GetOfficeNameForSoldier 根据民兵信息获取其所属机关的市、区、街道名称
func GetOfficeNameForSoldier(soldier *models.Soldier) (cityName, districtName, streetName string) {
	// 检查连接
	MyDB.DB().Ping()

	var err error
	streetOffice := models.Office{}
	cityOffice := models.Office{}
	districtOffice := models.Office{}

	err = MyDB.Where("office_id = ?", soldier.ServeOfficeID).First(&streetOffice).Error
	if err != nil {
		panic(err)
	}

	err = MyDB.Where("office_id = ?", streetOffice.HigherOfficeID).First(&districtOffice).Error
	if err != nil {
		panic(err)
	}

	err = MyDB.Where("office_id = ?", districtOffice.HigherOfficeID).First(&cityOffice).Error
	if err != nil {
		panic(err)
	}

	return cityOffice.Name, districtOffice.Name, streetOffice.Name
}
