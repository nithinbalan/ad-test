package models

import "time"

type PolicyReportClause struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID uint      `json:"policy_report_id" gorm:"not null;index"`
	ClauseNumber   *string   `json:"clause_number" gorm:"type:varchar(50);index"`
	ClauseTitleEn  string    `json:"clause_title_en" gorm:"type:nvarchar(255);not null"`
	ClauseTitleAr  string    `json:"clause_title_ar" gorm:"type:nvarchar(255)"`
	Description    *string   `json:"description" gorm:"type:nvarchar(max)"`
	Controls       *string   `json:"controls" gorm:"type:nvarchar(max)"`
	CreatedByID    uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportClause) TableName() string {
	return "policy_report_clauses"
}
