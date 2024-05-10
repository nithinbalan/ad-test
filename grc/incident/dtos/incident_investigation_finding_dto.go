package dtos

type IncidentInvestigationFindingDTO struct {
	ID                      uint   `json:"id"`
	IncidentInvestigationID uint   `json:"incident_investigation_id"`
	Title                   string `json:"title"`
}
