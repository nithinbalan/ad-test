package dtos

type TimezoneDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Offset      string `json:"offset"`
	DisplayName string `json:"display_name"`
}
