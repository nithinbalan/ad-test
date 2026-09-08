package models

import "time"

type PenetrationTestRetestFindingDocument struct {
	ID                             uint      `json:"id" gorm:"primary_key;auto_increment"`
	PenetrationTestRetestFindingID uint      `json:"penetration_test_retest_finding_id"`
	URL                            string    `json:"url"`
	Name                           string    `json:"name" binding:"required"`
	Size                           int       `json:"size"`
	Token                          string    `json:"token"`
	OriginalName                   string    `json:"original_name"`
	Type                           string    `json:"type"`
	CreatedBYID                    uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                    uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                      time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                      time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	PenetrationTestRetestFinding *PenetrationTestRetestFinding `gorm:"foreignKey:PenetrationTestRetestFindingID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser                *User                         `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                *User                         `gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTestRetestFindingDocument) TableName() string {
	return "penetration_test_retest_finding_documents"
}
