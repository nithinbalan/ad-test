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
	Competency        CompetencyDTO                    `json:"competency"`
	CompetencyLevelID *uint                            `json:"competency_level_id"`
	CompetencyLevel   CompetencyLevelDTO               `json:"competency_level"`
	Certifications    []JobDescriptionCertificationDTO `json:"certifications"`
}
type JobDescriptionCertificationDTO struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	IsMandatory bool   `json:"is_mandatory"`
}
type CompetencyLevelDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en" `
	TitleAr string `json:"title_ar" `
	Status  string `json:"status" `
}
type CompetencyDTO struct {
	ID                uint                      `json:"id"`
	TitleEn           string                    `json:"title_en" `
	TitleAr           string                    `json:"title_ar" `
	Status            string                    `json:"status" `
	CompetencyGroupID uint                      `json:"competency_group_id"`
	CompetencyGroup   CompetencyGroupDetailsDTO `json:"competency_group"`
}
type CompetencyGroupDetailsDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}
