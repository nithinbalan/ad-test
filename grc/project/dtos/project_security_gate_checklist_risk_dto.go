package dtos

import (
	"time"

	organizationDTOs "jethur/grc/organization/dtos"
	riskDTOs "jethur/grc/risk/dtos"
)

type ProjectSecurityGateChecklistRiskFormDTO struct {
	RiskCode          string `json:"risk_code" binding:"required,max=500"`
	RegistrationDate  string `json:"registration_date" binding:"required,datetime=2006-01-02"`
	Title             string `json:"title" binding:"required,max=255"`
	Justification     string `json:"justification" binding:"required"`
	DivisionID        uint   `json:"division_id" binding:"required"`
	ResponsibleUnitID uint   `json:"responsible_unit_id" binding:"required"`
	ResponsibleUserID uint   `json:"responsible_user_id" binding:"required"`
	RiskCategoryID    uint   `json:"risk_category_id" binding:"required"`
}

type ProjectSecurityGateChecklistRiskDTO struct {
	ID               uint                                  `json:"id"`
	ReferenceNumber  *string                               `json:"reference_number"`
	RiskCode         *string                               `json:"risk_code"`
	Title            string                                `json:"title"`
	RegistrationDate *time.Time                            `json:"registration_date"`
	Justification    *string                               `json:"justification"`
	RiskStatus       *riskDTOs.RiskStatusDTO               `json:"risk_status"`
	RiskCategory     *riskDTOs.RiskCategoryDTO             `json:"risk_category"`
	Division         *organizationDTOs.OrganizationUnitDTO `json:"division"`
	ResponsibleUnit  *organizationDTOs.OrganizationUnitDTO `json:"responsible_unit"`
	ResponsibleUser  *organizationDTOs.UserMinimalDTO      `json:"responsible_user"`
	CreatedAt        time.Time                             `json:"created_at"`
}
