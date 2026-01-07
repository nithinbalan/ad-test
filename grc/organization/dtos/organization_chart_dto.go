package dtos

type OrganizationChartUserListDTO struct {
	ID               uint    `json:"id"`
	FirstNameEn      string  `json:"first_name_en"`
	LastNameEn       string  `json:"last_name_en"`
	FirstNameAr      string  `json:"first_name_ar"`
	LastNameAr       string  `json:"last_name_ar"`
	Email            string  `json:"email"`
	MobileNumber     *string `json:"mobile_number"`
	ProfileImagePath *string `json:"profile_image_path"`
	Initial          string  `json:"initial"`
	Color            string  `json:"color"`

	Designation string `json:"designation"`

	PrimaryOrganizationID *uint  `json:"primary_organization_unit_id" gorm:"column:primary_organization_unit_id"`
	PrimaryOrganization   string `json:"primary_organization_unit" gorm:"column:primary_organization_unit"`
	
	ReportingUserID             *uint  `json:"reporting_user_id"`
	ReportingUserFirstNameEn    string `json:"reporting_user_first_name_en"`
	ReportingUserFirstNameAr    string `json:"reporting_user_first_name_ar"`

	RolesEn string `json:"roles_en"`
	RolesAr string `json:"roles_ar"`
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

	Designation           DesignationDTO `json:"designation"`
	PrimaryOrganizationID *uint          `json:"primary_organization_unit_id" mapstructure:"-"`
	PrimaryOrganization   string         `json:"primary_organization_unit" mapstructure:"-"`

	SystemRoles   []OrganizationChartRoleDTO         `json:"system_roles"`
	ReportingUser *OrganizationChartReportingUserDTO `json:"reporting_user"`
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

type OrganizationChartUserListQueryParamsDTO struct {
	Filters map[string][]string
}

type ChangeReportingUserFormDTO struct {
	ReportingUserID *uint `json:"reporting_user_id"`
}
