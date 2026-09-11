package models

import "time"

type ThreatIntelligencePlannedAdvisory struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" binding:"required,max=100" gorm:"type:varchar(100)"`
	CaseNumber               string     `json:"case_number" binding:"required" gorm:"type:nvarchar(255);not null"`
	TitleEn                  string     `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr                  string     `json:"title_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	ReceiptDate              time.Time  `json:"receipt_date" binding:"required" gorm:"type:datetime;not null"`
	CaseLink                 string     `json:"case_link" binding:"required" gorm:"type:nvarchar(max);not null"`
	AdvisorySourceID         uint       `json:"advisory_source_id" binding:"required" gorm:"not null"`
	CaseClassificationID     uint       `json:"case_classification_id" binding:"required" gorm:"not null"`
	AdvisoryTypeID           uint       `json:"advisory_type_id" binding:"required" gorm:"not null"`
	ImportanceID             uint       `json:"importance_id" binding:"required" gorm:"not null"`
	SecurityClassificationID uint       `json:"security_classification_id" binding:"required" gorm:"not null"`
	ThreatType               *string    `json:"threat_type" gorm:"type:nvarchar(max);default:null"`
	ThreatSeverityID         *uint      `json:"threat_severity_id" gorm:"default:null"`
	LikelihoodID             *uint      `json:"likelihood_id" gorm:"default:null"`
	RiskLevelID              *uint      `json:"risk_level_id" gorm:"default:null"`
	RiskLevelValueID         *uint      `json:"risk_level_value_id" gorm:"default:null"`
	Description              *string    `json:"description" gorm:"type:nvarchar(max);default:null"`
	AffectedAssets           *string    `json:"affected_assets" gorm:"type:nvarchar(max);default:null"`
	Tags                     *string    `json:"tags" gorm:"type:nvarchar(max);default:null"`
	IOCs                     *string    `json:"iocs" gorm:"column:iocs;type:nvarchar(max);default:null"`
	AdvisoryRecipientID      uint       `json:"advisory_recipient_id" binding:"required" gorm:"not null"`
	OrganizationUnitID       uint       `json:"organization_unit_id" binding:"required" gorm:"not null"`
	ExecutedByID             *uint      `json:"executed_by_id" gorm:"default:null"`
	AdvisoryStatusID         uint       `json:"advisory_status_id" binding:"required" gorm:"not null"`
	InternalStatusID         *uint      `json:"internal_status_id" gorm:"default:null"`
	NotifyDate               *time.Time `json:"notify_date" gorm:"type:datetime;default:null"`
	DeadlineDate             *time.Time `json:"deadline_date" gorm:"type:datetime;default:null"`
	ResolveDate              *time.Time `json:"resolve_date" gorm:"type:datetime;default:null"`
	ReplyDate                *time.Time `json:"reply_date" gorm:"type:datetime;default:null"`
	ProcessingTime           *string    `json:"processing_time" gorm:"type:nvarchar(255);default:null"`
	SAEDTicketNo             *string    `json:"saed_ticket_no" gorm:"type:nvarchar(255);default:null"`
	Attachment               *string    `json:"attachment" gorm:"type:nvarchar(max);default:null"`
	ReportLink               *string    `json:"report_link" gorm:"type:nvarchar(max);default:null"`
	Notes                    *string    `json:"notes" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AdvisorySource         *PlannedAdvisorySource                       `gorm:"foreignKey:AdvisorySourceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CaseClassification     *PlannedAdvisoryClassification               `gorm:"foreignKey:CaseClassificationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AdvisoryType           *PlannedAdvisoryType                         `gorm:"foreignKey:AdvisoryTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Importance             *PlannedAdvisoryImportance                   `gorm:"foreignKey:ImportanceID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SecurityClassification *PlannedAdvisorySecurityClassification       `gorm:"foreignKey:SecurityClassificationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AdvisoryRecipient      *User                                        `gorm:"foreignKey:AdvisoryRecipientID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OrganizationUnit       *OrganizationUnit                            `gorm:"foreignKey:OrganizationUnitID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ExecutedBy             *User                                        `gorm:"foreignKey:ExecutedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AdvisoryStatus         *PlannedAdvisoryStatus                       `gorm:"foreignKey:AdvisoryStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	InternalStatus         *PlannedAdvisoryInternalStatus               `gorm:"foreignKey:InternalStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                                        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ThreatSeverity         *ThreatIntelligenceSeverity                  `gorm:"foreignKey:ThreatSeverityID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskLevel              *ThreatIntelligenceRiskLevel                 `gorm:"foreignKey:RiskLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Likelihood             *ThreatIntelligenceLikelihood                `gorm:"foreignKey:LikelihoodID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskLevelValue         *ThreatIntelligenceRiskLevelValue            `gorm:"foreignKey:RiskLevelValueID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Documents              []*ThreatIntelligencePlannedAdvisoryDocument `gorm:"foreignKey:ThreatIntelligencePlannedAdvisoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by ThreatIntelligencePlannedAdvisory to `threat_intelligence_planned_advisories`
func (ThreatIntelligencePlannedAdvisory) TableName() string {
	return "threat_intelligence_planned_advisories"
}
