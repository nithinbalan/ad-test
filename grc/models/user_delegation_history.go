package models

import "time"

type UserDelegationHistory struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          uint       `json:"user_id" binding:"required"`
	DelegateeUserID uint       `json:"delegatee_user_id" binding:"required"`
	Comments        *string    `json:"comments"`
	AssignedDate    time.Time  `json:"assigned_date" binding:"required"`
	ModuleCount     int        `json:"module_count"`
	TaskCount       int        `json:"task_count"`
	CreatedBYID     uint       `json:"created_by_id" binding:"required"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID     *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
}

func (UserDelegationHistory) TableName() string {
	return "user_delegation_history"
}
