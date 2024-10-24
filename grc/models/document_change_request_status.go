package models

type DocumentChangeRequestStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=initiated rejected implementing completed cancelled" gorm:"type:varchar(10);not null;check:type IN ('initiated', 'rejected', 'implementing', 'completed', 'cancelled')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

// TableName overrides the table name used by DocumentChangeRequestStatus to `document_change_request_statuses`
func (DocumentChangeRequestStatus) TableName() string {
	return "document_change_request_statuses"
}
