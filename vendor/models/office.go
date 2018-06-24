package models

// Office office model
type Office struct {
	OfficeID       int32
	OfficeLevel    string
	HigherOfficeID int32
	Name           string
}

// TableName tn
func (Office) TableName() string {
	return "Offices"
}
