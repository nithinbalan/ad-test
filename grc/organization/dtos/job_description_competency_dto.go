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

type JobDescriptionCompetencyDetailDTO struct {
	ID                uint                             `json:"id"`
	CompetencyID      uint                             `json:"competency_id"`
	CompetencyTitleEn string                           `json:"competency_title_en"`
	CompetencyGroup   string                           `json:"competency_group"`
	CompetencyLevelID *uint                            `json:"competency_level_id,omitempty"`
	LevelTitleEn      *string                          `json:"level_title_en,omitempty"`
	Certifications    []JobDescriptionCertificationDTO `json:"certifications"`
}
type JobDescriptionCertificationDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	IsMandatory bool   `json:"is_mandatory"`
}
