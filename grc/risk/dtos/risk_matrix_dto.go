package dtos

type RiskMatrixDTO struct {
	ID              uint   `json:"id"`
	Title           string `json:"title"`
	CalculationType string `json:"calculation_type"`
	Type            string `json:"type"`
}
