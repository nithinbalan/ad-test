package models

import (
	"time"
)

// AssetPolicyChecklist represents the asset_policy_checklists table
type AssetPolicyChecklist struct {
	ID                         uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	AssetPolicyID              uint       `json:"asset_policy_id" gorm:"not null"`
	DocumentRequirementGroupID uint       `json:"document_requirement_group_id" gorm:"not null"`
	QuestionGroupTitle         string     `json:"question_group_title" gorm:"type:nvarchar(500);not null"`
	QuestionTitle              string     `json:"question_title" gorm:"type:nvarchar(500);not null"`
	Justification              *string    `json:"justification" gorm:"type:nvarchar(max)"`
	ComplianceStatusID         *uint      `json:"compliance_status_id"`
	AnsweredBYID               *uint      `json:"answered_by_id"`
	AnsweredAt                 *time.Time `json:"answered_at" gorm:"type:datetime"`
	CreatedAt                  time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                  *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBYID                uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                *uint      `json:"updated_by_id"`

	// Relationships
	AssetPolicy              AssetPolicy              `gorm:"foreignKey:AssetPolicyID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	DocumentRequirementGroup DocumentRequirementGroup `gorm:"foreignKey:DocumentRequirementGroupID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceStatus         *ComplianceStatus        `gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AnsweredByUser           *User                    `gorm:"foreignKey:AnsweredBYID"`
	CreatedByUser            User                     `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser            *User                    `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the default table name
func (AssetPolicyChecklist) TableName() string {
	return "asset_policy_checklists"
}
