package models

type DocumentChangeRequestStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=draft in-review rejected approved published archived reverted" gorm:"type:varchar(10);not null;default:'draft';check:type IN ('draft', 'in-review', 'rejected', 'approved', 'published', 'archived', 'reverted')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (DocumentChangeRequestStatus) TableName() string {
	return "document_change_request_statuses"
}
