package dtos

type UserTaskRawCountDTO struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type UserTaskStatusCountResponseDTO struct {
	Open       int `json:"open"`
	InProgress int `json:"in_progress"`
	Closed     int `json:"closed"`
	Total      int `json:"total"`
}
