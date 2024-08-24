package models

type Source struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn  string  `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr  string  `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Type     string  `json:"type" binding:"required,oneof=direct jethur-dm gdrive onedrive sharepoint dropbox" gorm:"type:varchar(20);not null;default:'direct';check:type IN ('direct','jethur-dm','gdrive','onedrive','sharepoint','dropbox')"`
	Status   string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (Source) TableName() string {
	return "sources"
}
