package models

import (
	"time"
)

type BcpTestAndExerciseSystemRegister struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpTestAndExerciseID uint      `json:"bcp_test_and_exercise_id" binding:"required"`
	SystemRegisterID     uint      `json:"system_register_id" binding:"required"`
	CreatedBYID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpTestAndExercise *BcpTestAndExercise `gorm:"foreignKey:BcpTestAndExerciseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SystemRegister     *SystemRegister     `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedBYID"`
}

func (BcpTestAndExerciseSystemRegister) TableName() string {
	return "bcp_test_and_exercise_system_registers"
}
