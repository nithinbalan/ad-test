package dtos

type RoleDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Status  string `json:"status"`
}

type RolePermissionFormDTO struct {
	ModuleID             uint `json:"module_id" binding:"required"`
	RolePermissionTypeID uint `json:"role_permission_type_id" binding:"required"`
	IsGranted            bool `json:"is_granted" binding:"required"`
}
