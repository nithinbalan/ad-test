package models

type RiskMatrix struct {
	ID              uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title           string `json:"title" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	CalculationType string `json:"calculation_type" binding:"required,max=20,oneof=multiplication addition" gorm:"type:varchar(20);not null"`
	Type            string `json:"type" binding:"required,max=10,oneof=erm isms" gorm:"type:varchar(10);not null"`

	//Relationship
	RiskMatrixImpact     []*RiskMatrixImpact     `gorm:"foreignKey:RiskMatrixID"`
	RiskMatrixLikelihood []*RiskMatrixLikelihood `gorm:"foreignKey:RiskMatrixID"`
	RiskMatrixRiskRating []*RiskMatrixRiskRating `gorm:"foreignKey:RiskMatrixID"`
}

// TableName overrides the table name used by RiskMatrix to `risk_matrix`
func (RiskMatrix) TableName() string {
	return "risk_matrix"
}
