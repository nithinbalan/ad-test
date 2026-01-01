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
	ThirdPartyAssessmentStatusID uint       `json:"third_party_assessment_id"`
	UID                          *string    `json:"uid"`
	Link                         *string    `json:"link"`
	IsSend                       bool       `json:"is_send"`
	SentBYID                     *uint      `json:"sent_by_id" gorm:"default:null"`
	SentAt                       *time.Time `json:"sent_at" gorm:"default:null"`
	ImportedBYID                 *uint      `json:"imported_by_id" gorm:"default:null"`
	ImportedAt                   *time.Time `json:"imported_at" gorm:"type:datetime;default:null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null;column:created_by_id"`
	UpdatedBYID                  *uint      `json:"updated_by_id,omitempty" gorm:"default:null;column:updated_by_id"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    *time.Time `json:"updated_at,omitempty" gorm:"type:datetime"`
	NextReviewLevel              *int       `json:"next_review_level" gorm:"default:null"`
	SubmittedByID                *uint      `json:"submitted_by_id" gorm:"default:null"`
	SubmittedAt                  *time.Time `json:"submitted_at" gorm:"default:null"`
	RiskLevelID                  *uint      `json:"risk_level_id"`
	Score                        *float64   `json:"score"`

	// Relationships
	ThirdPartyRegister         *ThirdPartyRegister         `gorm:"foreignKey:ThirdPartyRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	QuestionSet                *QuestionSet                `gorm:"foreignKey:QuestionSetID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ThirdPartyAssessmentStatus *ThirdPartyAssessmentStatus `gorm:"foreignKey:ThirdPartyAssessmentStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Sections                   []*TPRMAssessmentSection    `gorm:"foreignKey:ThirdPartyAssessmentID"`
	// Questions                  []*TPRMAssessmentSectionQuestion `gorm:"foreignKey:ThirdPartyAssessmentID"`
	RiskLevel       *RiskLevel `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User      `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser   *User      `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SubmittedByUser *User      `gorm:"foreignKey:SubmittedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ImportedByUser  *User      `gorm:"foreignKey:ImportedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SentByUser      *User      `gorm:"foreignKey:SentBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	ThirdPartyAssessmentWorkflows []*ThirdPartyAssessmentWorkflow `gorm:"foreignKey:ThirdPartyAssessmentID"`
}

// TableName overrides the default table name.
func (ThirdPartyAssessment) TableName() string {
	return "third_party_assessments"
}
