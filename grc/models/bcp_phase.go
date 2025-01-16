package models

type BcpPhase struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TimeFrame string `json:"time_frame" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by BcpPhase to `bcp_phases`
func (BcpPhase) TableName() string {
	return "bcp_phases"
}
