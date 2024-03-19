package models

import (
	"jethur/grc/settings/models"
	"time"
)

type User struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	OrganizationUnit1ID *uint     `gorm:"column:organization_unit_1_id"`
	OrganizationUnit2ID *uint     `gorm:"column:organization_unit_2_id"`
	OrganizationUnit3ID *uint     `gorm:"column:organization_unit_3_id"`
	OrganizationUnit4ID *uint     `gorm:"column:organization_unit_4_id"`
	OrganizationUnit5ID *uint     `gorm:"column:organization_unit_5_id"`
	FirstNameEN         string    `gorm:"column:first_name_en;size:100;not null"`
	LastNameEN          string    `gorm:"column:last_name_en;size:100;not null"`
	FirstNameAR         string    `gorm:"column:first_name_ar;size:45;not null"`
	LastNameAR          string    `gorm:"column:last_name_ar;size:45;not null"`
	Email               string    `gorm:"column:email;size:100;not null"`
	Otp                 string    `gorm:"column:otp;size:10;null"`
	SentAt              time.Time `gorm:"column:sent_at;null"`
	PasswordHash        string    `gorm:"column:password_hash;not null"`
	ProfileImagePath    *string   `gorm:"column:profile_image_path;size:255"`
	CreatedAt           time.Time `gorm:"column:created_at;not null;default:current_timestamp"`
	UpdatedAt           time.Time `gorm:"column:updated_at;not null;default:current_timestamp"`

	// Relationships
	OrganizationUnit1 *models.OrganizationUnit `gorm:"foreignKey:OrganizationUnit1ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit2 *models.OrganizationUnit `gorm:"foreignKey:OrganizationUnit2ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit3 *models.OrganizationUnit `gorm:"foreignKey:OrganizationUnit3ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit4 *models.OrganizationUnit `gorm:"foreignKey:OrganizationUnit4ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	OrganizationUnit5 *models.OrganizationUnit `gorm:"foreignKey:OrganizationUnit5ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "users"
}
