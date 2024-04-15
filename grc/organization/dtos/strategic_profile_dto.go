package dtos

import "time"

type StrategicProfileFormDTO struct {
	TitleEn   string  `json:"title_en" binding:"required"`
	TitleAr   string  `json:"title_ar" binding:"required"`
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
	Location  *string                        `json:"location,omitempty"`
	Email     *string                        `json:"email,omitempty"`
	Industry  *string                        `json:"industry,omitempty"`
	Phone     *string                        `json:"phone,omitempty"`
	LinkedIn  *string                        `json:"linkedin,omitempty"`
	X         *string                        `json:"x,omitempty"`
	Facebook  *string                        `json:"facebook,omitempty"`
	Instagram *string                        `json:"instagram,omitempty"`
	YouTube   *string                        `json:"youtube,omitempty"`
	AboutEn   *string                        `json:"about_en,omitempty"`
	AboutAr   *string                        `json:"about_ar,omitempty"`
	VisionEn  *string                        `json:"vision_en,omitempty"`
	VisionAr  *string                        `json:"vision_ar,omitempty"`
	MissionEn *string                        `json:"mission_en,omitempty"`
	MissionAr *string                        `json:"mission_ar,omitempty"`
	CreatedAt *time.Time                     `json:"created_at"`
	UpdatedAt *time.Time                     `json:"updated_at"`
	Values    []StrategicProfileValueFormDTO `json:"values,omitempty"`
}

type StrategicProfileVisionFormDTO struct {
	VisionEn *string `json:"vision_en"`
	VisionAr *string `json:"vision_ar"`
}

type StrategicProfileMisionFormDTO struct {
	MissionEn *string `json:"mission_en"`
	MissionAr *string `json:"mission_ar"`
}

type StrategicProfileValueFormDTO struct {
	ID            *uint   `json:"id"`
	TitleEn       string  `json:"title_en"`
	TitleAr       string  `json:"title_ar"`
	DescriptionEn *string `json:"description_en"`
	DescriptionAr *string `json:"description_ar"`
}
