package models

// ExportSoldier export data model for soldier
type ExportSoldier struct {
	Name          string
	IDNum         string
	PhoneNum      uint64
	Rank          string
	WechatOpenid  string
	CommanderName string
	City          string
	District      string
	Street        string
}
