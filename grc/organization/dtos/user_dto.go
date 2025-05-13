package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	settingDTOs "jethur/grc/settings/dtos"
	"time"
)

type UserListDTO struct {
	ID                       uint    `json:"id"`
	FirstNameEn              string  `json:"first_name_en"`
	LastNameEn               string  `json:"last_name_en"`
	FirstNameAr              string  `json:"first_name_ar"`
	LastNameAr               string  `json:"last_name_ar"`
	OrganizationUnit1TitleEn string  `json:"organization_unit_1_title_en" gorm:"column:organization_unit_1_title_en"`
	OrganizationUnit2TitleEn string  `json:"organization_unit_2_title_en,omitempty" gorm:"column:organization_unit_2_title_en"`
	OrganizationUnit3TitleEn string  `json:"organization_unit_3_title_en,omitempty" gorm:"column:organization_unit_3_title_en"`
	OrganizationUnit4TitleEn string  `json:"organization_unit_4_title_en,omitempty" gorm:"column:organization_unit_4_title_en"`
	OrganizationUnit5TitleEn string  `json:"organization_unit_5_title_en,omitempty" gorm:"column:organization_unit_5_title_en"`
	OrganizationUnit1TitleAr string  `json:"organization_unit_1_title_ar" gorm:"column:organization_unit_1_title_ar"`
	OrganizationUnit2TitleAr string  `json:"organization_unit_2_title_ar,omitempty" gorm:"column:organization_unit_2_title_ar"`
	OrganizationUnit3TitleAr string  `json:"organization_unit_3_title_ar,omitempty" gorm:"column:organization_unit_3_title_ar"`
	OrganizationUnit4TitleAr string  `json:"organization_unit_4_title_ar,omitempty" gorm:"column:organization_unit_4_title_ar"`
	OrganizationUnit5TitleAr string  `json:"organization_unit_5_title_ar,omitempty" gorm:"column:organization_unit_5_title_ar"`
	DesignationTitleEn       string  `json:"designation_title_en"`
	DesignationTitleAr       string  `json:"designation_title_ar"`
	RolesEn                  string  `json:"roles_en"`
	RolesAr                  string  `json:"roles_ar"`
	Email                    string  `json:"email"`
	ProfileImagePath         *string `json:"profile_image_path"`
	Initial                  string  `json:"initial"`
	Color                    string  `json:"color"`
	MobileNumber             string  `json:"mobile_number"`
	UserStatusTitleEn        string  `json:"user_status_title_en"`
	UserStatusTitleAr        string  `json:"user_status_title_ar"`
	UserStatusColorCode      string  `json:"user_status_color_code"`
	UserStatusIsActive       bool    `json:"user_status_is_active"`
}

type UserListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	UserIDs     []int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
	RoleIDs     []int
	DateFrom    string
	DateTo      string
}

type UserListResponseDTO struct {
	Items      []*UserListDTO               `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type UserFormDTO struct {
	OrganizationUnit1ID uint    `json:"organization_unit_1_id" binding:"required"`
	OrganizationUnit2ID *uint   `json:"organization_unit_2_id" binding:"omitempty"`
	OrganizationUnit3ID *uint   `json:"organization_unit_3_id" binding:"omitempty"`
	OrganizationUnit4ID *uint   `json:"organization_unit_4_id" binding:"omitempty"`
	OrganizationUnit5ID *uint   `json:"organization_unit_5_id" binding:"omitempty"`
	DesignationID       uint    `json:"designation_id" binding:"required"`
	FirstNameEn         string  `json:"first_name_en" binding:"required,max=255"`
	LastNameEn          string  `json:"last_name_en" binding:"required,max=255"`
	FirstNameAr         string  `json:"first_name_ar" binding:"required,max=255"`
	LastNameAr          string  `json:"last_name_ar" binding:"required,max=255"`
	Email               string  `json:"email" binding:"required,email"`
	MobileNumber        *string `json:"mobile_number" binding:"omitempty,max=45"`
	ProfileImagePath    *string `json:"profile_image_path" binding:"omitempty"`
	RoleIDs             []uint  `json:"role_ids" binding:"required"`
	UserStatusID        uint    `json:"user_status_id" binding:"required"`
	Password            string  `json:"password" binding:"omitempty"`
	IsSentEmail         bool    `json:"is_sent_email"`
}

type UserUpdateDTO struct {
	OrganizationUnit1ID *uint      `json:"organization_unit_1_id,omitempty"`
	OrganizationUnit2ID *uint      `json:"organization_unit_2_id,omitempty"`
	OrganizationUnit3ID *uint      `json:"organization_unit_3_id,omitempty"`
	OrganizationUnit4ID *uint      `json:"organization_unit_4_id,omitempty"`
	OrganizationUnit5ID *uint      `json:"organization_unit_5_id,omitempty"`
	DesignationID       *uint      `json:"designation_id,omitempty"`
	FirstNameEn         *string    `json:"first_name_en,omitempty"`
	LastNameEn          *string    `json:"last_name_en,omitempty"`
	FirstNameAr         *string    `json:"first_name_ar,omitempty"`
	LastNameAr          *string    `json:"last_name_ar,omitempty"`
	Email               *string    `json:"email,omitempty"`
	MobileNumber        *string    `json:"mobile_number,omitempty"`
	ProfileImagePath    *string    `json:"profile_image_path,omitempty"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
}

type SendNewPasswordEmailDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	RedirectURL string `json:"redirect_url"`
}

type OrganizationUnit1DTO struct {
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}
type UserMinimalDTO struct {
	ID                      int                             `json:"id"`
	FirstNameEN             string                          `json:"first_name_en"`
	LastNameEN              string                          `json:"last_name_en"`
	FirstNameAR             string                          `json:"first_name_ar"`
	LastNameAR              string                          `json:"last_name_ar"`
	ProfileImagePath        *string                         `json:"profile_image_path"`
	Email                   string                          `json:"email"`
	Initial                 string                          `json:"initial"`
	Color                   string                          `json:"color"`
	MobileNumber            *string                         `json:"mobile_number"`
	PrimaryOrganizationUnit settingDTOs.OrganizationUnitDTO `json:"primary_organization_unit"`
}

type UserRoleDTO struct {
	ID     uint                `json:"id"`
	UserID uint                `json:"user_id"`
	RoleID uint                `json:"role_id"`
	Role   settingDTOs.RoleDTO `json:"role"`
}

type UserDetailsDTO struct {
	ID               uint       `json:"id"`
	FirstNameEn      string     `json:"first_name_en"`
	LastNameEn       string     `json:"last_name_en"`
	FirstNameAr      string     `json:"first_name_ar"`
	LastNameAr       string     `json:"last_name_ar"`
	Email            string     `json:"email"`
	MobileNumber     *string    `json:"mobile_number"`
	ProfileImagePath *string    `json:"profile_image_path"`
	Initial          string     `json:"initial"`
	Color            string     `json:"color"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relationship DTOs
	OrganizationUnit1 settingDTOs.OrganizationUnitDTO  `json:"organization_unit_1"`
	OrganizationUnit2 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_2"`
	OrganizationUnit3 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_3"`
	OrganizationUnit4 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_4"`
	OrganizationUnit5 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_5"`
	UserStatus        *UserStatusDTO                   `json:"user_status"`
	Designation       *DesignationDTO                  `json:"designation"`
	Roles             *[]settingDTOs.RoleDTO           `json:"roles"`
	CreatedByUser     *UserMinimalDTO                  `json:"created_by_user"`
	UpdatedByUser     *UserMinimalDTO                  `json:"updated_by_user"`
}

type UserLoginHistoryListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	LoginAtFrom string
	LoginAtTo   string
}

type UserLoginHistoryListDTO struct {
	ID          uint      `json:"id"`
	LoginAt     time.Time `json:"login_at"`
	LoginStatus string    `json:"login_status"`
	IPAddress   string    `json:"ip_address"`
	Location    *string   `json:"location"`
	DeviceInfo  *string   `json:"device_info"`
	BrowserInfo *string   `json:"browser_info"`
}

type UserLoginHistoryListResponseDTO struct {
	Items      []*UserLoginHistoryListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type CurrentUserDetailsDTO struct {
	ID                  uint           `json:"id"`
	FirstNameEn         string         `json:"first_name_en"`
	LastNameEn          string         `json:"last_name_en"`
	FirstNameAr         string         `json:"first_name_ar"`
	LastNameAr          string         `json:"last_name_ar"`
	Email               string         `json:"email"`
	MobileNumber        *string        `json:"mobile_number"`
	ProfileImagePath    *string        `json:"profile_image_path"`
	Initial             string         `json:"initial"`
	Color               string         `json:"color"`
	LastLoginAt         *time.Time     `json:"last_login_at"`
	UserSetting         UserSettingDTO `json:"user_setting"`
	PrimaryUnitID       *uint          `json:"primary_unit_id"`
	OrganizationUnit1ID *uint          `json:"organization_unit_1_id"`
	OrganizationUnit2ID *uint          `json:"organization_unit_2_id"`
	OrganizationUnit3ID *uint          `json:"organization_unit_3_id"`
	OrganizationUnit4ID *uint          `json:"organization_unit_4_id"`
	OrganizationUnit5ID *uint          `json:"organization_unit_5_id"`

	// Relationship DTOs
	OrganizationUnit1 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_1"`
	OrganizationUnit2 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_2"`
	OrganizationUnit3 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_3"`
	OrganizationUnit4 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_4"`
	OrganizationUnit5 *settingDTOs.OrganizationUnitDTO `json:"organization_unit_5"`
	UserStatus        *UserStatusDTO                   `json:"user_status"`
	Designation       *DesignationDTO                  `json:"designation"`
	Roles             *[]settingDTOs.RoleDTO           `json:"roles"`
}

type UserNotificationListDTO struct {
	ID                          uint      `json:"id"`
	UserID                      uint      `json:"user_id"`
	UserFirstNameEn             string    `json:"user_first_name_en"`
	UserFirstNameAr             string    `json:"user_first_name_ar"`
	UserLastNameEn              string    `json:"user_last_name_en"`
	UserLastNameAr              string    `json:"user_last_name_ar"`
	UserInitial                 string    `json:"user_initial"`
	UserColor                   string    `json:"user_color"`
	UserEmail                   string    `json:"user_email"`
	UserProfileImagePath        string    `json:"user_profile_image_path"`
	ModuleID                    uint      `json:"module_id"`
	ModuleIdentifier            string    `json:"module_identifier"`
	ModuleTitleEn               string    `json:"module_title_en"`
	ModuleTitleAr               string    `json:"module_title_ar"`
	NotificationType            string    `json:"notification_type"`
	ReferenceNumber             string    `json:"reference_number"`
	HeadingEn                   string    `json:"heading_en"`
	HeadingAr                   string    `json:"heading_ar"`
	TitleEn                     string    `json:"title_en"`
	Entity                      string    `json:"entity"`
	Action                      string    `json:"action"`
	TitleAr                     string    `json:"title_ar"`
	URLParam                    string    `json:"url_param"`
	QueryParam                  string    `json:"query_param"`
	IsNew                       bool      `json:"is_new"`
	IsRead                      bool      `json:"is_read"`
	IsSystemNotification        bool      `json:"is_system_notification"`
	CreatedAt                   time.Time `json:"created_at"`
	CreatedUserID               uint      `json:"created_user_id"`
	CreatedUserFirstNameEn      string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string    `json:"created_user_last_name_ar"`
	CreatedUserInitial          string    `json:"created_user_initial"`
	CreatedUserColor            string    `json:"created_user_color"`
	CreatedUserEmail            string    `json:"created_user_email"`
	CreatedUserProfileImagePath string    `json:"created_user_profile_image_path"`
}

type UserNotificationListQueryParamsDTO struct {
	Sort        string `json:"sort" form:"sort"`
	SortBy      string `json:"sort_by" form:"sort_by"`
	Keyword     string `json:"keyword" form:"keyword"`
	ModuleIDs   []uint `json:"module_ids" form:"module_ids"`
	Page        int    `json:"page" form:"page"`
	PageSize    int    `json:"page_size" form:"page_size"`
	QuickFilter []string
}

type UserNotificationListResponseDTO struct {
	Items      []*UserNotificationListDTO   `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MarkNotificationAsReadFormDTO struct {
	Type                string `json:"type" binding:"required,oneof=all individual"`
	UserNotificationIDs []uint `json:"user_notification_ids" binding:"required_if=Type individual,dive,gt=0"`
}

type UserPrimaryDetailsDTO struct {
	ID               int     `json:"id"`
	FirstNameEN      string  `json:"first_name_en"`
	LastNameEN       string  `json:"last_name_en"`
	Email            string  `json:"email"`
	Designation      string  `json:"designation"`
	PrimaryUnit      string  `json:"primary_unit"`
	MobileNumber     *string `json:"mobile_number"`
	ProfileImagePath *string `json:"profile_image_path"`
}

type UpdateOrganizationUnitDTO struct {
	UserID             *uint `json:"user_id"`
	OrganizationUnitID uint  `json:"organization_unit_id" binding:"required"` // New Organization Unit ID
}
type ADUserDTO struct {
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            string  `json:"email"`
	OrganizationUnit *string `json:"organization_unit,omitempty"`
	Designation      *string `json:"designation,omitempty"`
}
