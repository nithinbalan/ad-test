package models

type KPILibrary struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Category    string  `json:"category" gorm:"type:nvarchar(255);not null"`
	Title       string  `json:"title" gorm:"type:nvarchar(500);not null"`
	Description *string `json:"description" gorm:"type:nvarchar(max)"`
	Unit        string  `json:"unit" gorm:"type:nvarchar(50);not null"`
	Frequency   string  `json:"frequency" gorm:"type:nvarchar(50);not null"`
}

func (KPILibrary) TableName() string {
	return "kpi_library"
}
