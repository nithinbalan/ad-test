package models

type DocumentChangeRequestWorkflowStatus struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,oneof=comment revert reject approve pending resubmit" gorm:"type:varchar(10);not null;check:type IN ('comment', 'revert', 'reject', 'approve', 'pending', 'resubmit')"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (DocumentChangeRequestWorkflowStatus) TableName() string {
	return "document_change_request_workflow_statuses"
}
