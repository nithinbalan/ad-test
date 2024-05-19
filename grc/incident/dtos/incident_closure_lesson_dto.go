package dtos

type IncidentClosureLessonDTO struct {
	ID                uint   `json:"id"`
	IncidentClosureID uint   `json:"incident_closure_id"`
	Lessons           string `json:"lessons"`
}
