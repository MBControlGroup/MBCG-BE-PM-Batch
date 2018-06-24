package models

// Soldier soldier model
type Soldier struct {
	SoldierID     int32
	Rank          string
	IDNum         string
	Name          string
	PhoneNum      uint64
	WechatOpenid  string
	CommanderID   int32
	ServeOfficeID int32
	ImUserID      int32
}

// TableName tn
func (Soldier) TableName() string {
	return "Soldiers"
}

// NewSoldier new soldier
func NewSoldier() *Soldier {
	return &Soldier{
		Rank:     "SD", // 默认军衔为基层民兵
		ImUserID: 1,    // 测试阶段暂定默认为1
	}
}
