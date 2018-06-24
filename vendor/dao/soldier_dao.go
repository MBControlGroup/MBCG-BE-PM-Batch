package dao

import (
	"models"
)

// GetSoldiers get soldiers in db
func GetSoldiers() []*models.Soldier {
	// 检查连接
	MyDB.DB().Ping()

	var soldiers []*models.Soldier
	err := MyDB.Find(&soldiers).Error
	if err != nil {
		panic(err)
	}
	return soldiers
}
