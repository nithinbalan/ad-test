package dtos

type IncidentComplianceRequirementDTO struct {
	ID                      *uint `json:"id"`
	IncidentRegisterID      uint  `json:"incident_register_id"`
	ComplianceRequirementID uint  `json:"compliance_requirement_id"`
}
