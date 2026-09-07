package models

import "time"

type ProjectSecurityGate struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProjectTitleEn    string     `json:"project_title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ProjectTitleAr    string     `json:"project_title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	StageID           uint       `json:"stage_id" binding:"required" gorm:"not null"`
	TemplateVersionID uint       `json:"template_version_id" gorm:"not null"`
	BusinessOwnerID   uint       `json:"business_owner_id" binding:"required" gorm:"not null"`
	ReadinessStatusID uint       `json:"readiness_status_id" gorm:"not null"`
	Notes             *string    `json:"notes" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Stage           *ProjectSecurityGateStage                    `json:"stage" gorm:"foreignKey:StageID"`
	TemplateVersion *ProjectSecurityGateChecklistTemplateVersion `json:"template_version" gorm:"foreignKey:TemplateVersionID"`
	BusinessOwner   *User                                        `json:"business_owner" gorm:"foreignKey:BusinessOwnerID"`
	ReadinessStatus *ProjectSecurityGateReadinessStatus          `json:"readiness_status" gorm:"foreignKey:ReadinessStatusID"`
	CreatedByUser   *User                                        `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User                                        `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID"`
	Classifications []ProjectSecurityGateClassificationMapping   `json:"classifications" gorm:"foreignKey:SecurityGateID"`
}

func (ProjectSecurityGate) TableName() string {
	return "project_security_gates"
}
