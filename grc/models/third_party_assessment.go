package models

import "time"

type ThirdPartyAssessment struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber              string     `json:"reference_number" gorm:"type:varchar(45);not null;unique" binding:"required"`
	Title                        string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	ThirdPartyRegisterID         uint       `json:"third_party_register_id" gorm:"not null" binding:"required"`
	QuestionSetID                uint       `json:"question_set_id" gorm:"not null" binding:"required"`
	Message                      *string    `json:"message" gorm:"type:nvarchar(MAX);default:null"`
	DueDate                      *time.Time `json:"due_date" gorm:"type:date;default:null"`
	ThirdPartyAssessmentStatusID uint       `jsom:"third_party_assessment_id"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                  *uint      `json:"updated_by_id,omitempty"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at,omitempty" gorm:"type:datetime"`

	// Relationships
	ThirdPartyRegister         *ThirdPartyRegister         `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	QuestionSet                *QuestionSet                `gorm:"foreignKey:QuestionSetID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThirdPartyAssessmentStatus *ThirdPartyAssessmentStatus `gorm:"foreignKey:ThirdPartyAssessmentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser              *User                       `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser              *User                       `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the default table name.
func (ThirdPartyAssessment) TableName() string {
	return "third_party_assessments"
}
