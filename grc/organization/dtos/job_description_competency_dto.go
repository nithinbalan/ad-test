package dtos

type JobDescriptionMapCompetenciesDTO struct {
	CompetencyIDs []uint `json:"competency_ids" binding:"required,min=1,dive,required"`
}

type UpdateJDCompetencyDTO struct {
	CompetencyLevelID uint                                  `json:"competency_level_id" binding:"required"`
	Certifications    []JobDescriptionCompetencyCertFormDTO `json:"certifications"`
}

type JobDescriptionCompetencyCertFormDTO struct {
	Title       string `json:"title" binding:"required,max=255"`
	IsMandatory bool   `json:"is_mandatory"`
}
