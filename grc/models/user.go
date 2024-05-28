package models

import (
	"time"
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
	ProfileImageInitial   string     `json:"profile_image_initial" gorm:"type:nvarchar(2);default:null"`
	ProfileImageColorCode string     `json:"profile_image_color_code" gorm:"type:varchar(7);default:null"`
	LastLoginAt           *time.Time `json:"last_login_at" gorm:"type:datetime;default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy             uint       `json:"created_by" gorm:"not null"`
	UpdatedBy             *uint      `json:"updated_by" gorm:"default:null"`

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
	CreatedByUser     *User               `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser     *User               `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}
