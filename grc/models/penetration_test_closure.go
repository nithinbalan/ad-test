package models

import "time"

type PenetrationTestClosure struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestID uint       `json:"penetration_test_id" gorm:"not null;unique"`
	ClosureOwnerID    uint       `json:"closure_owner_id" gorm:"not null" binding:"required"`
	ClosureDate       time.Time  `json:"closure_date" gorm:"type:date;not null" binding:"required"`
	ClosureRemarks    string     `json:"closure_remarks" gorm:"type:nvarchar(max);not null" binding:"required"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	ClosureOwner  *User `json:"closure_owner,omitempty" gorm:"foreignKey:ClosureOwnerID"`
	CreatedByUser *User `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestClosure) TableName() string {
	return "penetration_test_closures"
}
