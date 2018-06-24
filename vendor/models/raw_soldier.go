package models

/*
字段依次为：
        1. 姓名
        2. 手机号码
        3. 身份证号码
        4. 所属市
        5. 所属区
        6. 所属街道
        7. 直属指挥官姓名
*/

// RawSoldier 批量导入时的民兵数据model
type RawSoldier struct {
	Name          string
	PhoneNum      uint64
	IDNum         string
	City          string
	District      string
	Street        string
	CommanderName string
}
