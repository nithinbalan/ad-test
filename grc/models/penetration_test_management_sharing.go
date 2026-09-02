package models

import "time"

type PenetrationTestManagementSharing struct {
	ID                          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PenetrationTestID           uint       `json:"penetration_test_id" gorm:"not null;unique"`
	ResultsSharedWithManagement bool       `json:"results_shared_with_management" gorm:"not null;default:false"`
	ManagementNotifyDate        *time.Time `json:"management_notify_date,omitempty" gorm:"type:date;default:null"`
	ManagementResponseDate      *time.Time `json:"management_response_date,omitempty" gorm:"type:date;default:null"`
	RemediationNotes            *string    `json:"remediation_notes,omitempty" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID                 uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID                 *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                   *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	CreatedByUser *User `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestManagementSharing) TableName() string {
	return "penetration_test_management_sharings"
}
