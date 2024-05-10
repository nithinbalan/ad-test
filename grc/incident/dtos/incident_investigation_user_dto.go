package dtos

type IncidentInvestigationUserDTO struct {
	ID                      uint `json:"id"`
	IncidentInvestigationID int  `json:"incident_investigation_id"`
	UserID                  int  `json:"user_id"`
}
