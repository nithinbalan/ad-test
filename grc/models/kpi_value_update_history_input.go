package models

type KPIValueUpdateHistoryInput struct {
	ID                      uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	KPIValueUpdateHistoryID uint    `json:"kpi_value_update_history_id" gorm:"not null"`
	Variable                string  `json:"variable" gorm:"type:nvarchar(45);not null"`
	Input                   string  `json:"input" gorm:"type:nvarchar(255);not null"`
	Value                   float64 `json:"value" gorm:"type:decimal(18,2);not null"`
	CreatedBYID             *uint   `json:"created_by_id" gorm:"default:null"`

	KPIValueUpdateHistory *KPIValueUpdateHistory `gorm:"foreignKey:KPIValueUpdateHistoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (KPIValueUpdateHistoryInput) TableName() string {
	return "kpi_value_update_history_inputs"
}
