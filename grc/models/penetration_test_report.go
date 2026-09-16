package models

import "time"

type PenetrationTestReport struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string     `json:"title" binding:"required,max=255"`
	StartDate   *time.Time `json:"start_date" gorm:"type:date"`
	EndDate     *time.Time `json:"end_date" gorm:"type:date"`
	CreatedBYID uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	CreatedByUser    *User                                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Sections         []*PenetrationTestReportSection         `gorm:"foreignKey:PenetrationTestReportID"`
	PenetrationTests []*PenetrationTestReportPenetrationTest `gorm:"foreignKey:PenetrationTestReportID"`
}

func (PenetrationTestReport) TableName() string {
	return "penetration_test_reports"
}
