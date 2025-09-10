package models

import "time"

type Workflow struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn           *string    `json:"title_en"`
	TitleAr           *string    `json:"title_ar"`
	ModuleIdentifier  string     `json:"module_identifier" gorm:"size:255;not null" binding:"required"`
	Type              string     `json:"type" binding:"required,oneof=document risk_register audit" gorm:"type:varchar(100);not null;check:type IN ('document')"`
	DocumentTypeID    *uint      `json:"document_type_id" gorm:"default:null"`
	Status            string     `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(20);not null;check:status IN ('active', 'inactive')"`
	WorkflowContextID uint       `json:"workflow_context_id"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	DocumentType    *DocumentType    `gorm:"foreignKey:DocumentTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowContext WorkflowContext  `gorm:"foreignKey:WorkflowContextID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	WorkflowLevels  []*WorkflowLevel `gorm:"foreignKey:WorkflowID"`
}

func (Workflow) TableName() string {
	return "workflows"
}
