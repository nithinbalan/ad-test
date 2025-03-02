package dtos

type DocumentByStatusDTO struct {
	DocumentStatusID        uint   `json:"id"`
	DocumentStatusTitleEn   string `json:"title_en"`
	DocumentStatusTitleAr   string `json:"title_ar"`
	DocumentStatusColorCode string `json:"color_code"`
	DocumentStatusType      string `json:"type"`
	Count                   int64  `json:"count"`
}

type DocumentByTypeDTO struct {
	DocumentTypeID      uint   `json:"id"`
	DocumentTypeTitleEn string `json:"title_en"`
	DocumentTypeTitleAr string `json:"title_ar"`
	Count               int64  `json:"count"`
}

type DocumentByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    int64  `json:"count"`
}

type DocumentByCategoryDTO struct {
	CategoryID      uint   `json:"id"`
	CategoryTitleEn string `json:"title_en"`
	CategoryTitleAr string `json:"title_ar"`
	Count           int64  `json:"count"`
}

type DocumentByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"id"`
	ResponsibleUserFirstNameEN      string `json:"first_name_en"`
	ResponsibleUserFirstNameAR      string `json:"first_name_ar"`
	ResponsibleUserLastNameEN       string `json:"last_name_en"`
	ResponsibleUserLastNameAR       string `json:"last_name_ar"`
	ResponsibleUserTitleEN          string `json:"title_en"`
	ResponsibleUserTitleAR          string `json:"title_ar"`
	ResponsibleUserInitial          string `json:"initial"`
	ResponsibleUserColor            string `json:"color"`
	ResponsibleUserEmail            string `json:"email"`
	ResponsibleUserProfileImagePath string `json:"profile_image"`
	Count                           int64  `json:"count"`
}

type DocumentByAccessDTO struct {
	AccessID          uint   `json:"id"`
	AccessTypeTitleEN string `json:"title_en"`
	AccessTypeTitleAR string `json:"title_ar"`
	AccessType        string `json:"type"`
	ColorCode         string `json:"color_code"`
	Count             int64  `json:"count"`
}

type DocumentByYearDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type DocumentByMonthDTO struct {
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int64  `json:"count"`
}

type DocumentChangeRequestByTypeDTO struct {
	DocumentChangeRequestTypeID      uint   `json:"id"`
	DocumentChangeRequestTypeTitleEn string `json:"title_en"`
	DocumentChangeRequestTypeTitleAr string `json:"title_ar"`
	DocumentChangeRequestType        string `json:"type"`
	Count                            int64  `json:"count"`
}

type DocumentChangeRequestByStatusDTO struct {
	DocumentChangeRequestStatusID      uint   `json:"id"`
	DocumentChangeRequestStatusTitleEn string `json:"title_en"`
	DocumentChangeRequestStatusTitleAr string `json:"title_ar"`
	DocumentChangeRequestStatusType    string `json:"type"`
	ColorCode                          string `json:"color_code"`
	Count                              int64  `json:"count"`
}

type DocumentsArchivedByQuarterDTO struct {
	Year    int    `json:"year"`
	Quarter int    `json:"quarter"`
	TitleEN string `json:"title_en"`
	TitleAR string `json:"title_ar"`
	Count   int64  `json:"count"`
}
