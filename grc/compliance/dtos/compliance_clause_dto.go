package dtos

type ComplianceClauseDTO struct {
	ID       int    `json:"id"`
	ClauseNo string `json:"clause_no"`
	Title    string `json:"title"`
}
