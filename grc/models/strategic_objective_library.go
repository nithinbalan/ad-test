package models

type StrategicObjectiveLibrary struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title    string `json:"title" gorm:"type:nvarchar(500);not null"`
	Category string `json:"category" gorm:"type:nvarchar(255);not null"`
}

func (StrategicObjectiveLibrary) TableName() string {
	return "strategic_objective_library"
}
