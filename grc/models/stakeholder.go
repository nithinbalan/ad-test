package models

import "time"

type Stakeholder struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn               string     `json:"title_en" binding:"required"`
	TitleAr               string     `json:"title_ar" binding:"required"`
	TypeID                *uint      `json:"type_id"`
	StakeholderCategoryID *uint      `json:"stakeholder_category_id"`
	Role                  *string    `json:"role"`
	InfluenceLevelID      *uint      `json:"influence_level_id"`
	ImpactLevelID         *uint      `json:"impact_level_id"`
	ContactInformation    *string    `json:"contact_information"`
	PrimaryUnitID         *uint      `json:"organization_unit_id" gorm:"column:organization_unit_id"`
	ThirdPartyRegisterID  *uint      `json:"tird_party_register_id"`
	ClientID              *uint      `json:"client_id"`
	UserID                *uint      `json:"user_id"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null" binding:"omitempty"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null" binding:"omitempty"`

	// Relationships
	Type                *ProcessDependencyType `gorm:"foreignKey:TypeID"`
	StakeholderCategory *StakeholderCategory   `gorm:"foreignKey:StakeholderCategoryID"`
	InfluenceLevel      *SeverityLevel         `gorm:"foreignKey:InfluenceLevelID"`
	ImpactLevel         *SeverityLevel         `gorm:"foreignKey:ImpactLevelID"`
	OrganizationUnit    *OrganizationUnit      `gorm:"foreignKey:PrimaryUnitID"`
	ThirdPartyRegister  *ThirdPartyRegister    `gorm:"foreignKey:ThirdPartyRegisterID"`
	Client              *ThirdPartyRegister    `gorm:"foreignKey:ClientID"`
	User                *User                  `gorm:"foreignKey:UserID"`
	CreatedByUser       *User                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                  `gorm:"foreignKey:UpdatedBYID"`
}

func (Stakeholder) TableName() string {
	return "stakeholders"
}
