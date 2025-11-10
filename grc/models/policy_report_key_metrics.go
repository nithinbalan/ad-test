package models

import "time"

type PolicyReportKeyMetrics struct {
	ID                             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID                 uint      `json:"policy_report_id" gorm:"not null;index"`
	ClausesDefinedCount            int       `json:"clauses_defined_count" gorm:"not null;default:0"`
	TotalAcknowledgementCount      int       `json:"total_acknowledgement_count" gorm:"not null;default:0"`
	AcknowledgementsCompletedCount int       `json:"acknowledgements_completed_count" gorm:"not null;default:0"`
	AcknowledgementsPendingCount   int       `json:"acknowledgements_pending_count" gorm:"not null;default:0"`
	AcknowledgementPercentage      float32   `json:"acknowledgement_percentage" gorm:"type:numeric(5,2);not null;default:0"`
	RisksReportedCount             int       `json:"risks_reported_count" gorm:"not null;default:0"`
	AssociatedStandardsCount       int       `json:"associated_standards_count" gorm:"not null;default:0"`
	AssociatedClausesCount         int       `json:"associated_clauses_count" gorm:"not null;default:0"`
	StrategicObjectivesLinkedCount int       `json:"strategic_objectives_linked_count" gorm:"not null;default:0"`
	ApprovalCycleCount             int       `json:"approval_cycle_count" gorm:"not null;default:0"`
	CreatedByID                    uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt                      time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportKeyMetrics) TableName() string {
	return "policy_report_key_metrics"
}
