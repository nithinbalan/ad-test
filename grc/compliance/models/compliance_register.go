package models

type ComplianceRegister struct {
	ID                      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceCategoryID    uint   `json:"compliance_category_id" binding:"required" gorm:"not null"`
	ComplianceStatusID      uint   `json:"compliance_status_id" binding:"required" gorm:"not null"`
	ResponsibleUnitID       uint   `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID       uint   `json:"responsible_user_id" binding:"required" gorm:"not null"`
	ReviewFrequencyID       uint   `json:"review_frequency_id" binding:"required" gorm:"not null"`
	ReviewResponsibleUserID uint   `json:"review_responsible_user_id" gorm:"default:null"`
	SeverityLevelID         uint   `json:"severity_level_id" gorm:"default:null"`
	ReferenceNumber         string `json:"reference_number" gorm:"size:45"`
	StandardCode            string `json:"standard_code" gorm:"size:45"`
	TitleEn                 string `json:"title_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	TitleAr                 string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	DescriptionEn           string `json:"description_en" gorm:"type:varchar(max)"`
	DescriptionAr           string `json:"description_ar" gorm:"type:nvarchar(max)"`
	Version                 string `json:"version" gorm:"size:45"`
	// Deadline                string `json:"deadline" gorm:"type:date"`
	// IssueDate               string `json:"issue_date" gorm:"type:date"`
	// ExpiryDate              string `json:"expiry_date" gorm:"type:date"`
	// NextReviewDate          string `json:"next_review_date" gorm:"type:date"`
	// CreatedAt               string `json:"created_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	// UpdatedAt               string `json:"updated_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	CreatedBy         uint `json:"created_by" gorm:"not null"`
	UpdatedBy         uint `json:"updated_by" gorm:"default:null"`
	LifecycleStatusID uint `json:"lifecycle_status_id" gorm:"not null"`
}

// TableName overrides the table name used by ComplianceRegister to `compliance_register`
func (ComplianceRegister) TableName() string {
	return "compliance_register"
}
