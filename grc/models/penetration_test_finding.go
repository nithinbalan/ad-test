package models

import "time"

type PenetrationTestFinding struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber       string     `json:"reference_number" gorm:"type:varchar(45);not null;unique"`
	PenetrationTestID     uint       `json:"penetration_test_id" gorm:"not null" binding:"required"`
	Title                 string     `json:"title" gorm:"type:nvarchar(255);not null" binding:"required"`
	SeverityLevelID       uint       `json:"severity_level_id" gorm:"not null" binding:"required"`
	Description           string     `json:"description" gorm:"type:nvarchar(max);not null" binding:"required"`
	Recommendation        string     `json:"recommendation" gorm:"type:nvarchar(max);not null" binding:"required"`
	OwnerID               uint       `json:"owner_id" gorm:"not null" binding:"required"`
	TargetRemediationDate *time.Time `json:"target_remediation_date" gorm:"type:date;default:null"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID           *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	PenetrationTest *PenetrationTest `json:"penetration_test,omitempty" gorm:"foreignKey:PenetrationTestID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel   *SeverityLevel   `json:"severity_level,omitempty" gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Owner           *User            `json:"owner,omitempty" gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser   *User            `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestFinding) TableName() string {
	return "penetration_test_findings"
}
