package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	riskDTOs "jethur/grc/risk/dtos"
)

type BcpKeyProcessListDTO struct {
	ID                           uint    `json:"id"`
	ProcessID                    uint    `json:"process_id"`
	ProcessReferenceNumber       string  `json:"process_reference_number"`
	ProcessTitleEn               string  `json:"process_title_en"`
	ProcessTitleAr               string  `json:"process_title_ar"`
	ProcessDescription           string  `json:"process_description"`
	ProcessOwnerUnitTitleEn      *string `json:"process_owner_unit_title_en"`
	ProcessOwnerUnitTitleAr      *string `json:"process_owner_unit_title_ar"`
	ProcessOwnerID               *uint   `json:"process_owner_id"`
	ProcessOwnerFirstNameEn      *string `json:"process_owner_first_name_en"`
	ProcessOwnerFirstNameAr      *string `json:"process_owner_first_name_ar"`
	ProcessOwnerLastNameEn       *string `json:"process_owner_last_name_en"`
	ProcessOwnerLastNameAr       *string `json:"process_owner_last_name_ar"`
	ProcessOwnerInitial          string  `json:"process_owner_initial"`
	ProcessOwnerColor            string  `json:"process_owner_color"`
	ProcessOwnerEmail            string  `json:"process_owner_email"`
	ProcessOwnerProfileImagePath string  `json:"process_owner_profile_image_path"`
	LifecycleStatusTitleEn       string  `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr       string  `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode     string  `json:"lifecycle_status_color_code"`
	RiskRatingTitleEn            *string `json:"risk_rating_title_en"`
	RiskRatingTitleAr            *string `json:"risk_rating_title_ar"`
	RiskRatingColorCode          *string `json:"risk_rating_color_code"`
	RiskRatingValue              *uint   `json:"risk_rating_value"`
}

type BcpKeyProcessListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}

// type BcpKeyProcessListResponseDTO struct {
// 	Items      []*BcpKeyProcessListDTO      `json:"items"`
// 	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
// }

type BcpBiaMatrixDetailDTO struct {
	ID                 uint                    `json:"id"`
	ReferenceNumber    string                  `json:"reference_number"`
	BusinessFunctionID uint                    `json:"business_function_id"`
	Rto                *int                    `json:"rto"`
	RtoUnitId          *uint                   `json:"rto_unit_id"`
	RtoUnit            *masterDTOs.TimeUnitDTO `json:"rto_unit"`
	Mtpd               *int                    `json:"mtpd"`
	MtpdUnitId         *uint                   `json:"mtpd_unit_id"`
	MtpdUnit           *masterDTOs.TimeUnitDTO `json:"mtpd_unit"`
	Rpo                *int                    `json:"rpo"`
	RpoUnitId          *uint                   `json:"rpo_unit_id"`
	RpoUnit            *masterDTOs.TimeUnitDTO `json:"rpo_unit"`
	Processes          []*BcpProcessDetailsDTO `json:"processes"`
}

type BcpProcessDetailsDTO struct {
	ID          uint                    `json:"id"`
	Process     ProcessDetailDTO        `json:"process"`
	Rto         *int                    `json:"rto"`
	RtoUnitId   *uint                   `json:"rto_unit_id"`
	RtoUnit     *masterDTOs.TimeUnitDTO `json:"rto_unit"`
	Mtpd        *int                    `json:"mtpd"`
	MtpdUnitId  *uint                   `json:"mtpd_unit_id"`
	MtpdUnit    *masterDTOs.TimeUnitDTO `json:"mtpd_unit"`
	Rpo         *int                    `json:"rpo"`
	RpoUnitId   *uint                   `json:"rpo_unit_id"`
	RpoUnit     *masterDTOs.TimeUnitDTO `json:"rpo_unit"`
	RiskLevelId *uint                   `json:"risk_level_id"`
	RiskLevel   riskDTOs.RiskLevelDTO   `json:"risk_level"`
}
type ProcessDetailDTO struct {
	ID              uint   `json:"id"`
	ReferenceNumber string `json:"reference_number"`
	TitleEn         string `json:"title_en"`
	TitleAr         string `json:"title_ar"`
}

type BcpProcessFormDTO struct {
	Rto        *int              `json:"rto"`
	RtoUnitId  *uint             `json:"rto_unit_id"`
	Mtpd       *int              `json:"mtpd"`
	MtpdUnitId *uint             `json:"mtpd_unit_id"`
	Rpo        *int              `json:"rpo"`
	RpoUnitId  *uint             `json:"rpo_unit_id"`
	Processes  []*ProcessFormDTO `json:"processes"`
}

type ProcessFormDTO struct {
	ID          uint  `json:"id"`
	ProcessID   uint  `json:"process_id"`
	RiskLevelID *uint `json:"risk_level_id"`
	Rto         *int  `json:"rto"`
	RtoUnitId   *uint `json:"rto_unit_id"`
	Mtpd        *int  `json:"mtpd"`
	MtpdUnitId  *uint `json:"mtpd_unit_id"`
	Rpo         *int  `json:"rpo"`
	RpoUnitId   *uint `json:"rpo_unit_id"`
}
