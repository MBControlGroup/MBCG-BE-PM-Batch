package models

// OfficeRelationship office relationship model
type OfficeRelationship struct {
	OfrID          int32
	HigherOfficeID int32
	LowerOfficeID  int32
}

// TableName tn
func (OfficeRelationship) TableName() string {
	return "OfficeRelationships"
}
