package dtos

type OrganizationStructureDTO struct {
	ID       uint                       `json:"id"`
	TitleEn  string                     `json:"title_en"`
	TitleAr  string                     `json:"title_ar"`
	Children []OrganizationStructureDTO `json:"children,omitempty"`
}
