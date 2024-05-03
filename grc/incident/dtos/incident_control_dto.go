package dtos

type IncidentControlDTO struct {
	ID                 uint   `json:"id"`
	IncidentRegisterID uint   `json:"incident_register_id"`
	ControlID          uint   `json:"control_id"`
}