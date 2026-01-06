package dtos

type OrganizationChartUserDTO struct {
	ID              uint   `json:"id"`
	FirstNameEn     string `json:"first_name_en"`
	LastNameEn      string `json:"last_name_en"`
	FirstNameAr     string `json:"first_name_ar"`
	LastNameAr      string `json:"last_name_ar"`
	Email           string `json:"email"`
	Designation     string `json:"designation"`
	ReportingUserID *uint  `json:"reporting_user_id"`
}


type OrganizationChartUserDetailsDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameAr       string  `json:"last_name_ar"`
	Email            string  `json:"email"`
	MobileNumber     *string `json:"mobile_number"`
	ProfileImagePath *string `json:"profile_image_path"`
	IsHead           bool    `json:"is_head"`

	Designation           string  `json:"designation" mapstructure:"-"`
	PrimaryOrganizationID *uint   `json:"primary_organization_unit_id" mapstructure:"-"`
	PrimaryOrganization   string  `json:"primary_organization_unit" mapstructure:"-"`

	SystemRoles []OrganizationChartRoleDTO         `json:"system_roles"`
	ReportingTo *OrganizationChartReportingUserDTO `json:"reporting_to"`
}


type OrganizationChartRoleDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type OrganizationChartReportingUserDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	ProfileImagePath *string `json:"profile_image_path"`
	Designation      string  `json:"designation"`
}
