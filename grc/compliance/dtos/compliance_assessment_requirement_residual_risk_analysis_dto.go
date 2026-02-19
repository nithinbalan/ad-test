package dtos

import organizationDTOs "jethur/grc/organization/dtos"

type ComplianceAssessmentRequirementResidualRiskFormDTO struct {
	ComplianceAssessmentImpactTypeID uint    `json:"compliance_assessment_impact_type_id" binding:"required"`
	Score                            *int    `json:"score" binding:"required,min=0,max=5"`
	Comment                          *string `json:"comment"`
}

type ComplianceAssessmentRequirementResidualRiskDTO struct {
	ID                             uint                              `json:"id"`
	ComplianceAssessmentImpactType ComplianceAssessmentImpactTypeDTO `json:"compliance_assessment_impact_type"`
	Score                          int                               `json:"score"`
	Comment                        *string                           `json:"comment"`
}

type ComplianceAssessmentResidualControlFormDTO struct {
	ComplianceAssessmentImpactTypeID uint  `json:"compliance_assessment_impact_type_id" binding:"required"`
	ControlID                        *uint `json:"control_id"`
}

type ComplianceAssessmentResidualControlListDTO struct {
	ID                               uint    `json:"id"`
	ComplianceAssessmentImpactTypeID uint    `json:"compliance_assessment_impact_type_id"`
	ImpactTypeTitleEn                string  `json:"impact_type_title_en"`
	ImpactTypeTitleAr                string  `json:"impact_type_title_ar"`
	InherentRiskScore                *string `json:"inherent_risk_score"`
	ControlID                        uint    `json:"control_id"`
	ControlReferenceNumber           string  `json:"control_reference_number"`
	ControlTitleEn                   string  `json:"control_title_en"`
	ControlTitleAr                   string  `json:"control_title_ar"`
	ControlDescription               *string `json:"control_description"`
	LinkControl                      bool    `json:"link_control"`
	IsNotApplicable                  bool    `json:"is_not_applicable"`
}

type ComplianceAssessmentResidualControlListQueryParamsDTO struct {
	Sort        string   `json:"sort"`
	SortBy      string   `json:"sort_by"`
	Keyword     string   `json:"keyword"`
	Status      string   `json:"status"`
	QuickFilter []string `json:"quick_filter"`
}

type ComplianceAssessmentResidualControlListResponseDTO struct {
	Items []ComplianceAssessmentResidualControlListDTO `json:"items"`
}

type ComplianceAssessmentResidualControlDTO struct {
	ID                             uint                               `json:"id"`
	ComplianceAssessmentImpactType ComplianceAssessmentImpactTypeDTO  `json:"compliance_assessment_impact_type"`
	Control                        *organizationDTOs.ControlDetailDTO `json:"control"`
	LinkControl                    bool                               `json:"link_control"`
	IsNotApplicable                bool                               `json:"is_not_applicable"`
}
