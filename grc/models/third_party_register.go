package models

import "time"

type ThirdPartyRegister struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber    string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	Title              string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	ThirdPartyTypeID   uint       `json:"third_party_type_id" gorm:"not null" binding:"required"`
	CountryID          uint       `json:"country_id" gorm:"not null" binding:"required"`
	Address            *string    `json:"address,omitempty" gorm:"type:nvarchar(MAX);default:null"`
	Website            *string    `json:"website,omitempty" gorm:"type:nvarchar(255);default:null"`
	OrganizationUnitID *uint      `json:"organization_unit_id,omitempty"`
	RepresentativeName *string    `json:"representative_name,omitempty" gorm:"type:nvarchar(250);default:null"`
	RepresentativeRole *string    `json:"representative_role,omitempty" gorm:"type:nvarchar(250);default:null"`
	Email              *string    `json:"email,omitempty" gorm:"type:nvarchar(150);default:null"`
	PhoneNumber        *string    `json:"phone_number,omitempty" gorm:"type:varchar(20);default:null"`
	ResponsibleUserID  *uint      `json:"responsible_user_id,omitempty"`
	ResponsibleUnitID  *uint      `json:"responsible_unit_id"`
	OnboardingDate     *time.Time `json:"onboarding_date,omitempty" gorm:"type:date;default:null"`
	ServiceStartDate   *time.Time `json:"service_start_date,omitempty" gorm:"type:date;default:null"`
	ServiceEndDate     *time.Time `json:"service_end_date,omitempty" gorm:"type:date;default:null"`
	SeverityLevelID    *uint      `json:"severity_level_id,omitempty"`

	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty" gorm:"type:datetime"`

	// Relationships
	ThirdPartyType   *ThirdPartyType   `gorm:"foreignKey:ThirdPartyTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Country          *Country          `gorm:"foreignKey:CountryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel    *SeverityLevel    `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser  *User             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit  *OrganizationUnit `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (ThirdPartyRegister) TableName() string {
	return "third_party_register"
}
