package dtos

type StrategicObjectiveListDTO struct {
	ID                     uint                                          `json:"id"`
	TitleEn                string                                        `json:"title_en"`
	TitleAr                string                                        `json:"title_ar"`
	Status                 string                                        `json:"status"`
	ComplianceRequirements []*ComplianceRequirementStrategicObjectiveDTO `json:"compliance_requirements"`
}

type StrategicObjectiveFormDTO struct {
	TitleEn string `json:"title_en" binding:"required"`
	TitleAr string `json:"title_ar" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

type StrategicObjectiveDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type ComplianceRequirementStrategicObjectiveDTO struct {
	ID                      uint                     `json:"id"`
	ComplianceRequirementID uint                     `json:"compliance_requirement_id"`
	StrategicObjectiveID    uint                     `json:"strategic_objective_id"`
	ComplianceRequirement   ComplianceRequirementDTO `json:"compliance_requirement"`
}

type ComplianceRequirementDTO struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
