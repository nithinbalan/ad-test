package dtos

type MeetingByTypeDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Count   int    `json:"count"`
}

type MeetingByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type MeetingByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type MeetingByOrganizerDTO struct {
	UserID               uint   `json:"id"`
	UserFirstNameEn      string `json:"first_name_en"`
	UserFirstNameAr      string `json:"first_name_ar"`
	UserLastNameEn       string `json:"last_name_en"`
	UserLastNameAr       string `json:"last_name_ar"`
	UserTitleEn          string `json:"title_en"`
	UserTitleAr          string `json:"title_ar"`
	UserInitial          string `json:"initial"`
	UserColor            string `json:"color"`
	UserEmail            string `json:"email"`
	UserProfileImagePath string `json:"profile_image_path"`
	Count                int    `json:"count"`
}

type MeetingMinutesByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type MeetingActionPlanByStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type MeetingActionPlanByResponsibleUnitDTO struct {
	ResponsibleUnitID        uint   `json:"id"`
	ResponsibleUnitTitleEn   string `json:"title_en"`
	ResponsibleUnitTitleAr   string `json:"title_ar"`
	ResponsibleUnitColorCode string `json:"color_code"`
	Count                    uint   `json:"count"`
}

type MeetingActionPlanByResponsibleUserDTO struct {
	ID                   uint   `json:"id"`
	UserFirstNameEn      string `json:"first_name_en"`
	UserFirstNameAr      string `json:"first_name_ar"`
	UserLastNameEn       string `json:"last_name_en"`
	UserLastNameAr       string `json:"last_name_ar"`
	UserTitleEn          string `json:"title_en"`
	UserTitleAr          string `json:"title_ar"`
	UserInitial          string `json:"initial"`
	UserColor            string `json:"color"`
	UserEmail            string `json:"email"`
	UserProfileImagePath string `json:"profile_image_path"`
	Count                int    `json:"count"`
}

type ActionPlanByOpenVsClosedDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
	Count     int    `json:"count"`
}

type MeetingMOMCountDTO struct {
	WithMOMCount    int `json:"with_mom_count"`
	WithoutMOMCount int `json:"without_mom_count"`
}

type MeetingWithMOMDTO struct {
	MeetingID    uint   `json:"meeting_id"`
	MeetingTitle string `json:"meeting_title"`
	Minutes      string `json:"minutes"`
	Status       string `json:"status"`
}

type MeetingWithOpenActionPlanDTO struct {
	MeetingID           uint   `json:"meeting_id"`
	MeetingTitle        string `json:"meeting_title"`
	OpenActionPlanCount int    `json:"open_action_plan_count"`
}

type MeetingWithOverdueActionPlanDTO struct {
	MeetingID              uint   `json:"meeting_id"`
	MeetingTitle           string `json:"meeting_title"`
	OverdueActionPlanCount int    `json:"overdue_action_plan_count"`
}

type MeetingCountByYearDTO struct {
	TitleEn string `json:"title_en"` // Year in English
	TitleAr string `json:"title_ar"` // Year in Arabic
	Count   uint   `json:"count"`    // Count of meetings
}

type MeetingCountByUserDTO struct {
	UserID           uint   `json:"id"`
	UserFirstNameEn  string `json:"first_name_en"`
	UserFirstNameAr  string `json:"first_name_ar"`
	UserLastNameEn   string `json:"last_name_en"`
	UserLastNameAr   string `json:"last_name_ar"`
	UserTitleEn      string `json:"title_en"`
	UserTitleAr      string `json:"title_ar"`
	UserInitial      string `json:"initial"`
	UserColor        string `json:"color"`
	UserEmail        string `json:"email"`
	UserProfileImage string `json:"profile_image_path"`
	MeetingCount     int    `json:"count"`
}
