package models

type ComplianceAssessmentRating struct {
	ID        uint     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string   `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string   `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode string   `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Type      string   `json:"type" binding:"required,max=10,oneof=critical high medium low closed" gorm:"type:varchar(10);not null"`
	ScoreFrom *float64 `json:"score_from"`
	ScoreTo   *float64 `json:"score_to"`
	Status    string   `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ComplianceAssessmentRating to `compliance_assessment_ratings`
func (ComplianceAssessmentRating) TableName() string {
	return "compliance_assessment_ratings"
}
