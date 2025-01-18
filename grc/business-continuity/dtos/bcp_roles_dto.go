package dtos

type BcpRolesListDTO struct {
	ID                         uint    `json:"id"`
	BusinessContinuityPlanID   uint    `json:"business_continuity_plan_id"`
	ResponsibleUserID          uint    `json:"responsible_user_id"`
	ResponsibleUnitID          uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn     string  `json:"responsible_unit_title_en"`
	Responsibility             string  `json:"responsibility"`
	BackupUserID               *uint   `json:"backup_user_id"`
	CreatedByID                uint    `json:"created_by_id"`
	ResponsibleUserFirstNameEn string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn  string  `json:"responsible_user_last_name_en"`
	ResponsibleUserEmail       string  `json:"responsible_user_email"`
	BackupUserFirstNameEn      string  `json:"backup_user_first_name_en"`
	BackupUserLastNameEn       string  `json:"backup_user_last_name_en"`
	BackupUserEmail            string  `json:"backup_user_email"`
	CreatedAt                  string  `json:"created_at"`
	UpdatedAt                  *string `json:"updated_at"`
	CreatedUserFirstNameEn     string  `json:"created_user_first_name_en"`
	CreatedUserLastNameEn      string  `json:"created_user_last_name_en"`
}

type BcpRolesFormDTO struct {
	ResponsibleUserID uint   `json:"responsible_user_id" binding:"required"`
	ResponsibleUnitID uint   `json:"responsible_unit_id" binding:"required"`
	Responsibility    string `json:"responsibility" binding:"required,max=5000"`
	BackupUserID      *uint  `json:"backup_user_id"`
}

type BcpRolesListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}
