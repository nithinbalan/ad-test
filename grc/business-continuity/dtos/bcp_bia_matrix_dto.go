package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	riskDTOs "jethur/grc/risk/dtos"
)

type BcpBiaMatrixDetailDTO struct {
	ID                 uint                    `json:"id"`
	ReferenceNumber    string                  `json:"reference_number"`
	BusinessFunctionID uint                    `json:"business_function_id"`
	Rto                *uint                   `json:"rto"`
	RtoUnitId          *uint                   `json:"rto_unit_id"`
	RtoUnit            *masterDTOs.TimeUnitDTO `json:"rto_unit"`
	Mtpd               *uint                   `json:"mtpd"`
	MtpdUnitId         *uint                   `json:"mtpd_unit_id"`
	MtpdUnit           *masterDTOs.TimeUnitDTO `json:"mtpd_unit"`
	Rpo                *uint                   `json:"rpo"`
	RpoUnitId          *uint                   `json:"rpo_unit_id"`
	RpoUnit            *masterDTOs.TimeUnitDTO `json:"rpo_unit"`
	Processes          []*BcpProcessDetailsDTO `json:"processes"`
}

type BcpProcessDetailsDTO struct {
	ID          uint                    `json:"id"`
	Process     ProcessDetailDTO        `json:"process"`
	Rto         *uint                   `json:"rto"`
	RtoUnitId   *uint                   `json:"rto_uint_id"`
	RtoUnit     *masterDTOs.TimeUnitDTO `json:"rto_unit"`
	Mtpd        *uint                   `json:"mtpd"`
	MtpdUnitId  *uint                   `json:"mtpd_uint_id"`
	MtpdUnit    *masterDTOs.TimeUnitDTO `json:"mtpd_unit"`
	Rpo         *uint                   `json:"rpo"`
	RpoUnitId   *uint                   `json:"rpo_uint_id"`
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
	Rto        *uint             `json:"rto"`
	RtoUnitId  *uint             `json:"rto_unit_id"`
	Mtpd       *uint             `json:"mtpd"`
	MtpdUnitId *uint             `json:"mtpd_unit_id"`
	Rpo        *uint             `json:"rpo"`
	RpoUnitId  *uint             `json:"rpo_unit_id"`
	Processes  []*ProcessFormDTO `json:"processes"`
}

type ProcessFormDTO struct {
	ID          uint  `json:"id"`
	ProcessID   uint  `json:"process_id"`
	RiskLevelID *uint `json:"risk_level_id"`
	Rto         *uint `json:"rto"`
	RtoUnitId   *uint `json:"rto_unit_id"`
	Mtpd        *uint `json:"mtpd"`
	MtpdUnitId  *uint `json:"mtpd_unit_id"`
	Rpo         *uint `json:"rpo"`
	RpoUnitId   *uint `json:"rpo_unit_id"`
}
