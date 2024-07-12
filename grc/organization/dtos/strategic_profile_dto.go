package dtos

import "time"

type StrategicProfileFormDTO struct {
	TitleEn   string  `json:"title_en" binding:"required"`
	TitleAr   string  `json:"title_ar"`
	Address   string  `json:"address" binding:"required"`
	Location  *string `json:"location"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Industry  *string `json:"industry"`
	Phone     *string `json:"phone"`
	Linkedin  *string `json:"linkedin"`
	X         *string `json:"x"`
	Facebook  *string `json:"facebook"`
	Instagram *string `json:"instagram"`
	Youtube   *string `json:"youtube"`
	AboutEn   *string `json:"about_en"`
	AboutAr   *string `json:"about_ar"`
}

type StrategicProfileDetailsDTO struct {
	ID        int                            `json:"id"`
	TitleEn   string                         `json:"title_en"`
	TitleAr   string                         `json:"title_ar"`
	Address   string                         `json:"address"`
	Location  *string                        `json:"location"`
	Email     *string                        `json:"email"`
	Industry  *string                        `json:"industry"`
	Phone     *string                        `json:"phone"`
	LinkedIn  *string                        `json:"linkedin"`
	X         *string                        `json:"x"`
	Facebook  *string                        `json:"facebook"`
	Instagram *string                        `json:"instagram"`
	YouTube   *string                        `json:"youtube"`
	AboutEn   *string                        `json:"about_en"`
	AboutAr   *string                        `json:"about_ar"`
	VisionEn  *string                        `json:"vision_en"`
	VisionAr  *string                        `json:"vision_ar"`
	MissionEn *string                        `json:"mission_en"`
	MissionAr *string                        `json:"mission_ar"`
	CreatedAt *time.Time                     `json:"created_at"`
	UpdatedAt *time.Time                     `json:"updated_at"`
	Values    []StrategicProfileValueFormDTO `json:"values"`
}

type StrategicProfileVisionFormDTO struct {
	VisionEn *string `json:"vision_en" binding:"omitempty,max=500"`
	VisionAr *string `json:"vision_ar" binding:"omitempty,max=500"`
}

type StrategicProfileMissionFormDTO struct {
	MissionEn *string `json:"mission_en" binding:"omitempty,max=500"`
	MissionAr *string `json:"mission_ar" binding:"omitempty,max=500"`
}

type StrategicProfileValueFormDTO struct {
	ID            *uint   `json:"id"`
	TitleEn       string  `json:"title_en" binding:"max=255"`
	TitleAr       string  `json:"title_ar" binding:"max=255"`
	DescriptionEn *string `json:"description_en" binding:"omitempty,max=500"`
	DescriptionAr *string `json:"description_ar" binding:"omitempty,max=500"`
}

type StrategicProfileValueDetailsDTO struct {
	ID               int                        `json:"id"`
	TitleEn          string                     `json:"title_en" `
	TitleAr          string                     `json:"title_ar"`
	DescriptionEn    *string                    `json:"description_en"`
	DescriptionAr    *string                    `json:"description_ar"`
	StrategicProfile StrategicProfileDetailsDTO `json:"strategic_profile"`
}
