package models

// UserRole represents the 'user_roles' table that links users and roles in a many-to-many relationship
type UserRole struct {
	ID     uint `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	RoleID uint `gorm:"not null" json:"role_id"`

	// You can also include navigation properties to User and Role if needed
	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}

// TableName overrides the table name used by UserRole to `user_roles`
func (UserRole) TableName() string {
	return "user_roles"
}
