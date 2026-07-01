package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ProcessEmployeeAnalysisListDTO struct {
	ID                          uint                                 `json:"id"`
	ProcessID                   uint                                 `json:"process_id"`
	ProcessTitleEn              string                               `json:"process_title_en"`
	ProcessTitleAr              string                               `json:"process_title_ar"`
	DesignationID               uint                                 `json:"designation_id"`
	DesignationTitleEn          string                               `json:"designation_title_en"`
	DesignationTitleAr          string                               `json:"designation_title_ar"`
	NormalCount                 *int                                 `json:"normal_count"`
	CriticalCount               *int                                 `json:"critical_count"`
	CreatedBYID                 *uint                                `json:"created_by_id"`
	CreatedUserFirstNameEn      string                               `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                               `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                               `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                               `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                               `json:"created_user_initial"`
	CreatedUserColor            string                               `json:"created_user_color"`
	CreatedUserEmail            string                               `json:"created_user_email"`
	CreatedUserProfileImagePath string                               `json:"created_user_profile_image_path"`
	UpdatedBYID                 *uint                                `json:"updated_by_id"`
	UpdatedUserFirstNameEn      string                               `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string                               `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string                               `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string                               `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string                               `json:"updated_user_initial"`
	UpdatedUserColor            string                               `json:"updated_user_color"`
	UpdatedUserEmail            string                               `json:"updated_user_email"`
	UpdatedUserProfileImagePath string                               `json:"updated_user_profile_image_path"`
	CreatedAt                   string                               `json:"created_at"`
	UpdatedAt                   string                               `json:"updated_at"`
	Employees                   []ProcessEmployeeAnalysisEmployeeDTO `json:"employees" gorm:"-"`
}

type ProcessEmployeeAnalysisListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Export   bool
}

type ProcessEmployeeAnalysisListResponseDTO struct {
	Items      []*ProcessEmployeeAnalysisListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO      `json:"pagination"`
}

type ProcessEmployeeAnalysisDetailDTO struct {
	ID            uint                                 `json:"id"`
	ProcessID     uint                                 `json:"process_id"`
	Process       *ProcessMinimalDTO                   `json:"process"`
	DesignationID uint                                 `json:"designation_id"`
	Designation   *DesignationDTO                      `json:"designation"`
	NormalCount   *int                                 `json:"normal_count"`
	CriticalCount *int                                 `json:"critical_count"`
	Employees     []ProcessEmployeeAnalysisEmployeeDTO `json:"employees"`
	CreatedByUser *UserMinimalDTO                      `json:"created_by"`
	UpdatedByUser *UserMinimalDTO                      `json:"updated_by"`
	CreatedAt     time.Time                            `json:"created_at"`
	UpdatedAt     *time.Time                           `json:"updated_at"`
}

type ProcessEmployeeAnalysisFormDTO struct {
	DesignationID uint                       `json:"designation_id" binding:"required"`
	NormalCount   *int                       `json:"normal_count"`
	CriticalCount *int                       `json:"critical_count"`
	Users         []EmployeeWithToolsFormDTO `json:"users"`
}

type EmployeeWithToolsFormDTO struct {
	UserID uint                  `json:"user_id" binding:"required"`
	Tools  []EmployeeToolFormDTO `json:"equipment_tools"`
}

type EmployeeToolFormDTO struct {
	EquipmentToolID uint `json:"equipment_tool_id" binding:"required"`
	IsHas           bool `json:"is_has" binding:"omitempty"`
	IsNeeded        bool `json:"is_needed" binding:"omitempty"`
}

type ProcessEmployeeAnalysisEmployeeDTO struct {
	ID                        uint                                     `json:"id"`
	ProcessEmployeeAnalysisID uint                                     `json:"process_employee_analysis_id"`
	UserID                    uint                                     `json:"user_id"`
	User                      UserMinimalDTO                           `json:"user"`
	Tools                     []ProcessEmployeeAnalysisEmployeeToolDTO `json:"equipment_tools"`
}

type ProcessEmployeeAnalysisEmployeeToolDTO struct {
	ID            uint                    `json:"id"`
	EquipmentTool ProcessEquipmentToolDTO `json:"equipment_tool"`
	IsHas         bool                    `json:"is_has"`
	IsNeeded      bool                    `json:"is_needed"`
}
