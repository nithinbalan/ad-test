package models

type RiskLibrary struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Category    string  `json:"category" gorm:"type:nvarchar(255);not null"`
	Title       string  `json:"title" gorm:"type:nvarchar(500);not null"`
	Description *string `json:"description" gorm:"type:nvarchar(max)"`
	Likelihood  string  `json:"likelihood" gorm:"type:nvarchar(45);not null"`
	Impact      string  `json:"impact" gorm:"type:nvarchar(45);not null"`    
}

func (RiskLibrary) TableName() string {
	return "risk_library"
}
