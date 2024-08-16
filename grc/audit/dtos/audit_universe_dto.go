package dtos

type AuditUniverseFormDTO struct {
	ProcessIDs               []*uint `json:"process_ids" binding:"omitempty"`
	ControlIDs               []*uint `json:"control_ids" binding:"omitempty"`
	DocumentIDs              []*uint `json:"document_ids" binding:"omitempty"`
	IncidentRegisterIDs      []*uint `json:"incident_register_ids" binding:"omitempty"`
	ComplianceRequirementIDs []*uint `json:"compliance_requirement_ids" binding:"omitempty"`
	RiskRegisterIDs          []*uint `json:"risk_register_ids" binding:"omitempty"`
	UserIDs                  []*uint `json:"user_ids" binding:"omitempty"`
	Comment                  string  `json:"comment" binding:"required"`
	AuditorRiskRatingID      uint    `json:"auditor_risk_rating_id" binding:"required"`
}
