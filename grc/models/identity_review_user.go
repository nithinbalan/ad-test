package models

import "time"

type IdentityReviewUser struct {
	ID            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Username      string     `json:"username" gorm:"type:nvarchar(255);not null"`
	DepartmentID  uint       `json:"department_id" gorm:"not null"`
	PrivilegeID   *uint      `json:"privilege_id" gorm:"default:null"`
	AccountTypeID *uint      `json:"account_type_id" gorm:"default:null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID   uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID   *uint      `json:"updated_by_id" gorm:"default:null"`

	Department    *OrganizationUnit  `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Privilege     *PrivilegedAccount `json:"privilege" gorm:"foreignKey:PrivilegeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	AccountType   *AccountType       `json:"account_type" gorm:"foreignKey:AccountTypeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser *User              `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	UpdatedByUser *User              `json:"updated_by" gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (IdentityReviewUser) TableName() string {
	return "identity_review_users"
}
