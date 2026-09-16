package dtos

type PenetrationTestByTypeDTO struct {
	PenetrationTestTypeID      uint   `json:"id"`
	PenetrationTestTypeTitleEn string `json:"title_en"`
	PenetrationTestTypeTitleAr string `json:"title_ar"`
	Count                      uint   `json:"count"`
}

type PenetrationTestByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type PenetrationTestByClassificationDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type PenetrationTestByPlanDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type PenetrationTestByExecutorDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type PenetrationTestByResponsibleUserDTO struct {
	ResponsibleUserID               uint    `json:"id"`
	ResponsibleUserFirstNameEn      string  `json:"first_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"first_name_ar"`
	ResponsibleUserLastNameEn       string  `json:"last_name_en"`
	ResponsibleUserLastNameAr       string  `json:"last_name_ar"`
	ResponsibleUserTitleEn          string  `json:"title_en"`
	ResponsibleUserTitleAr          string  `json:"title_ar"`
	ResponsibleUserInitial          *string `json:"initial"`
	ResponsibleUserColor            *string `json:"color"`
	ResponsibleUserEmail            string  `json:"email"`
	ResponsibleUserProfileImagePath *string `json:"profile_image"`
	Count                           int64   `json:"count"`
}

type PenetrationTestByCreatedUserDTO struct {
	CreatedUserID               uint    `json:"id"`
	CreatedUserFirstNameEn      string  `json:"first_name_en"`
	CreatedUserFirstNameAr      string  `json:"first_name_ar"`
	CreatedUserLastNameEn       string  `json:"last_name_en"`
	CreatedUserLastNameAr       string  `json:"last_name_ar"`
	CreatedUserTitleEn          string  `json:"title_en"`
	CreatedUserTitleAr          string  `json:"title_ar"`
	CreatedUserInitial          *string `json:"initial"`
	CreatedUserColor            *string `json:"color"`
	CreatedUserEmail            string  `json:"email"`
	CreatedUserProfileImagePath *string `json:"profile_image"`
	Count                       int64   `json:"count"`
}

type PenetrationTestByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint    `json:"id"`
	ResponsibleUnitTitleEn   string  `json:"title_en"`
	ResponsibleUnitTitleAr   string  `json:"title_ar"`
	ResponsibleUnitColorCode *string `json:"color_code"`
	Count                    uint    `json:"count"`
}

type PenetrationTestByYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type PenetrationTestByPeriodDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   uint   `json:"count"`
}

type PenetrationTestFindingByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}

type PenetrationTestFindingByResponsibleUnitDTO struct {
	ResponsibleUnitID      uint   `json:"id"`
	ResponsibleUnitTitleEn string `json:"title_en"`
	ResponsibleUnitTitleAr string `json:"title_ar"`
	ColorCode              string `json:"color_code"`
	Count                  uint   `json:"count"`
}

type PenetrationTestFindingByResponsibleUserDTO struct {
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

type PenetrationTestFindingBySeverityLevelDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     uint   `json:"count"`
}
