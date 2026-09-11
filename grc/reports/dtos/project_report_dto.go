package dtos

type ProjectByTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type ProjectByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type ProjectByPhaseDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type ProjectByClassificationDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type ProjectByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type ProjectByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type ProjectByYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type ProjectActionPlanByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type ProjectActionPlanByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEn      string `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"last_name_en"`
	ResponsibleUserLastNameAr       string `json:"last_name_ar"`
	ResponsibleUserTitleEn          string `json:"title_en"`
	ResponsibleUserTitleAr          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type ProjectActionPlanByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"id"`
	ResponsibleUnitTitleEn string `json:"title_en"`
	ResponsibleUnitTitleAr string `json:"title_ar"`
	ColorCode              string `json:"color_code"`
	Count                  uint   `json:"count"`
}
