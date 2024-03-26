package models

type ComplianceClause struct {
	ID                   uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint   `json:"compliance_register_id" binding:"required" gorm:"not null"`
	ClauseNo             string `json:"clause_no" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	Title                string `json:"title" binding:"required" gorm:"type:varchar(max);not null"`

	// Relationships
	ComplianceRegister ComplianceRegister `gorm:"foreignkey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by ComplianceClause to `compliance_clauses`
func (ComplianceClause) TableName() string {
	return "compliance_clauses"
}
