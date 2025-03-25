package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	OrganizationUnit1ID   uint       `json:"organization_unit_1_id" binding:"required" gorm:"column:organization_unit_1_id;not null"`
	OrganizationUnit2ID   *uint      `json:"organization_unit_2_id" gorm:"column:organization_unit_2_id;default:null"`
	OrganizationUnit3ID   *uint      `json:"organization_unit_3_id" gorm:"column:organization_unit_3_id;default:null"`
	OrganizationUnit4ID   *uint      `json:"organization_unit_4_id" gorm:"column:organization_unit_4_id;default:null"`
	OrganizationUnit5ID   *uint      `json:"organization_unit_5_id" gorm:"column:organization_unit_5_id;default:null"`
	UserStatusID          uint       `json:"user_status_id" binding:"required" gorm:"not null"`
	DesignationID         uint       `json:"designation_id" binding:"required" gorm:"not null"`
	FirstNameEn           string     `json:"first_name_en" binding:"required, max:255" gorm:"type:varchar(255);not null"`
	LastNameEn            string     `json:"last_name_en" binding:"require, max:255" gorm:"type:varchar(255);not null"`
	FirstNameAr           string     `json:"first_name_ar" binding:"require, max:255" gorm:"type:nvarchar(255);not null"`
	LastNameAr            string     `json:"last_name_ar" binding:"require, max:255" gorm:"type:nvarchar(255);not null"`
	Email                 string     `json:"email" binding:"require, email, max:255" gorm:"type:varchar(255);not null"`
	MobileNumber          *string    `json:"mobile_number" binding:"require, max:45" gorm:"type:varchar(45);not null"`
	ProfileImagePath      *string    `json:"profile_image_path" gorm:"type:varchar(500);default:null"`
	Initial               string     `json:"initial" gorm:"type:nvarchar(2);default:null"`
	Color                 string     `json:"color" gorm:"type:varchar(10);default:null"`
	LastLoginAt           *time.Time `json:"last_login_at" gorm:"type:datetime;default:null"`
	MarkedAsAuditableItem bool       `json:"marked_as_auditable_item" gorm:"default:false" binding:"-"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID           *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	OrganizationUnit1 *OrganizationUnit   `gorm:"foreignKey:OrganizationUnit1ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit2 *OrganizationUnit   `gorm:"foreignKey:OrganizationUnit2ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit3 *OrganizationUnit   `gorm:"foreignKey:OrganizationUnit3ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit4 *OrganizationUnit   `gorm:"foreignKey:OrganizationUnit4ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit5 *OrganizationUnit   `gorm:"foreignKey:OrganizationUnit5ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	UserStatus        *UserStatus         `gorm:"foreignKey:UserStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Designation       *Designation        `gorm:"foreignKey:DesignationID"`
	UserRoles         *[]UserRole         `gorm:"foreignKey:UserID"`
	LoginHistory      *[]UserLoginHistory `gorm:"foreignKey:UserID"`
	UserSetting       *UserSetting        `gorm:"foreignKey:UserID"`
	CreatedByUser     *User               `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser     *User               `gorm:"foreignKey:UpdatedByID"`
	// Add this virtual field
	PrimaryOrganizationUnit *OrganizationUnit `gorm:"-"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}

// Scope function to preload all necessary data
func WithPrimaryOrganizationUnit() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Query to check if unit is primary
		primaryCheck := "EXISTS (SELECT 1 FROM organization_levels ol WHERE ol.id = organization_units.organization_level_id AND ol.is_primary = 1)"

		// Only load the unit if it's primary
		return db.Preload("OrganizationUnit1", primaryCheck).
			Preload("OrganizationUnit1.OrganizationLevel", "is_primary = 1").
			Preload("OrganizationUnit2", primaryCheck).
			Preload("OrganizationUnit2.OrganizationLevel", "is_primary = 1").
			Preload("OrganizationUnit3", primaryCheck).
			Preload("OrganizationUnit3.OrganizationLevel", "is_primary = 1").
			Preload("OrganizationUnit4", primaryCheck).
			Preload("OrganizationUnit4.OrganizationLevel", "is_primary = 1").
			Preload("OrganizationUnit5", primaryCheck).
			Preload("OrganizationUnit5.OrganizationLevel", "is_primary = 1")
	}
}

// Helper function to check if a unit is primary
func checkUnitIsPrimary(unit *OrganizationUnit) bool {
	return unit != nil &&
		unit.OrganizationLevel != nil &&
		unit.OrganizationLevel.IsPrimary
}

// Get the primary organization unit
func (u *User) GetPrimaryOrganizationUnit() *OrganizationUnit {
	if u == nil {
		return nil
	}

	if checkUnitIsPrimary(u.OrganizationUnit1) {
		return u.OrganizationUnit1
	}
	if checkUnitIsPrimary(u.OrganizationUnit2) {
		return u.OrganizationUnit2
	}
	if checkUnitIsPrimary(u.OrganizationUnit3) {
		return u.OrganizationUnit3
	}
	if checkUnitIsPrimary(u.OrganizationUnit4) {
		return u.OrganizationUnit4
	}
	if checkUnitIsPrimary(u.OrganizationUnit5) {
		return u.OrganizationUnit5
	}
	return nil
}

// AfterFind hook to automatically set primary unit
func (u *User) AfterFind(tx *gorm.DB) error {
	u.PrimaryOrganizationUnit = u.GetPrimaryOrganizationUnit()
	return nil
}
