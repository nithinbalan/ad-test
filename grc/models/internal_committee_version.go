package models

import (
	"time"
)

type InternalCommitteeVersion struct {
	ID                        uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	InternalCommitteeID       uint       `json:"internal_committee_id" binding:"required" gorm:"not null"`
	CommitteeNameEn           string     `json:"committee_name_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	CommitteeNameAr           string     `json:"committee_name_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	CommitteeTypeID           uint       `json:"committee_type_id" binding:"required" gorm:"not null"`
	CommitteeTermID           uint       `json:"committee_term_id" binding:"required" gorm:"not null"`
	TermStartDate             *time.Time `json:"term_start_date" gorm:"type:date"`
	TermEndDate               *time.Time `json:"term_end_date" gorm:"type:date"`
	CommitteeClassificationID uint       `json:"committee_classification_id" binding:"required" gorm:"not null"`
	Justification             *string    `json:"justification" gorm:"type:nvarchar(max)"`
	RolesAndResponsibilities  *string    `json:"roles_and_responsibilities" gorm:"type:nvarchar(max)"`
	CommitteeAuthorities      *string    `json:"committee_authorities" gorm:"type:nvarchar(max)"`
	CommitteeStatusID         uint       `json:"committee_status_id" binding:"required" gorm:"not null"`
	IsDisbandRequested        bool       `json:"is_disband_requested" gorm:"type:bit;default:0"`
	Version                   *float64   `json:"version" gorm:"type:decimal(4,2)"`
	IsLatestVersion           *bool      `json:"is_latest_version" gorm:"type:bit"`
	CreatedBYID               uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID               *uint      `json:"updated_by_id"`
	CreatedAt                 time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                 *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	SubmittedByID             *uint      `json:"submitted_by_id"`
	ApprovedByID              *uint      `json:"approved_by_id"`
	SubmittedAt               *time.Time `json:"submitted_at" gorm:"type:datetime;default:null"`
	ApprovedAt                *time.Time `json:"approved_at" gorm:"type:datetime;default:null"`
	NextReviewLevel           *int       `json:"next_review_level"`
	WorkflowStatusID          *uint      `json:"workflow_status_id"`
	CanceledByID              *uint      `json:"canceled_by_id"`
	CanceledAt                *time.Time `json:"canceled_at" gorm:"type:datetime;default:null"`
	DisbandByID               *uint      `json:"disband_by_id"`
	DisbandAt                 *time.Time `json:"disband_at" gorm:"type:datetime;default:null"`
	CompletedByID             *uint      `json:"completed_by_id"`
	CompletedAt               *time.Time `json:"completed_at" gorm:"type:datetime;default:null"`
	IsDetailsUpdated          bool       `json:"is_details_updated" gorm:"type:bit;default:0"`
	// Relationships
	InternalCommittee                          *InternalCommittee                           `gorm:"foreignKey:InternalCommitteeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CommitteeType                              *CommitteeType                               `gorm:"foreignKey:CommitteeTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CommitteeTerm                              *CommitteeTerm                               `gorm:"foreignKey:CommitteeTermID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CommitteeClassification                    *CommitteeClassification                     `gorm:"foreignKey:CommitteeClassificationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CommitteeStatus                            *InternalCommitteeStatus                     `gorm:"foreignKey:CommitteeStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	WorkflowStatus                             *WorkflowStatus                              `gorm:"foreignKey:WorkflowStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	InternalCommitteeVersionNomineeDepartments []*InternalCommitteeVersionNomineeDepartment `gorm:"foreignKey:InternalCommitteeVersionID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`

	CreatedByUser   *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser   *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubmittedByUser *User `gorm:"foreignKey:SubmittedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ApprovedByUser  *User `gorm:"foreignKey:ApprovedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CanceledByUser  *User `gorm:"foreignKey:CanceledByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	DisbandByUser   *User `gorm:"foreignKey:DisbandByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CompletedByUser *User `gorm:"foreignKey:CompletedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by InternalCommitteeVersion
func (InternalCommitteeVersion) TableName() string {
	return "internal_committee_versions"
}
