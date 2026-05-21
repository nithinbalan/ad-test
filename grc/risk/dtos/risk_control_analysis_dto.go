package dtos

import organizationDTOs "jethur/grc/organization/dtos"

type RiskControlAnalysisListDTO struct {
	ID                                   uint                                            `json:"id"`
	ControlID                            uint                                            `json:"control_id"`
	ControlTestID                        *uint                                           `json:"control_test_id"`
	ControlReferenceNumber               string                                          `json:"control_reference_number"`
	ControlCode                          *string                                         `json:"control_code"`
	ControlTitleEn                       string                                          `json:"control_title_en"`
	ControlTitleAr                       string                                          `json:"control_title_ar"`
	Score                                string                                          `json:"score"`
	ControlDescription                   *string                                         `json:"control_description"`
	ControlTypeTitleEn                   string                                          `json:"control_type_title_en"`
	ControlTypeTitleAr                   string                                          `json:"control_type_title_ar"`
	ControlType                          *string                                         `json:"control_type"`
	ControlCategoryTitleEn               string                                          `json:"control_category_title_en"`
	ControlCategoryTitleAr               string                                          `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string                                         `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string                                         `json:"owner_unit_title_ar"`
	OwnerID                              *uint                                           `json:"owner_id"`
	OwnerFirstNameEn                     *string                                         `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string                                         `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string                                         `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string                                         `json:"owner_last_name_ar"`
	OwnerInitial                         string                                          `json:"owner_initial"`
	OwnerColor                           string                                          `json:"owner_color"`
	OwnerEmail                           string                                          `json:"owner_email"`
	OwnerProfileImagePath                string                                          `json:"owner_profile_image_path"`
	Type                                 string                                          `json:"type" binding:"required"`
	ControlImplementationStatusTitleEn   *string                                         `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   *string                                         `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode *string                                         `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string                                         `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string                                         `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string                                         `json:"control_efficiency_status_color_code"`
	ControlEffectivenessScore            *float64                                        `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID         *uint                                           `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn    *string                                         `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    *string                                         `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingMinScore   *float64                                        `json:"control_effectiveness_rating_min_score"`
	ControlEffectivenessRatingMaxScore   *float64                                        `json:"control_effectiveness_rating_max_score"`
	ControlEffectivenessRatingType       *string                                         `json:"control_effectiveness_rating_type"`
	ControlEffectivenessRatingColorCode  *string                                         `json:"control_effectiveness_rating_color_code"`
	RiskReduction                        *float64                                        `json:"risk_reduction"`
	ControlAutomationID                  *uint                                           `json:"control_automation_id"`
	ControlAutomationTitleEn             *string                                         `json:"control_automation_title_en"`
	ControlAutomationTitleAr             *string                                         `json:"control_automation_title_ar"`
	ControlAutomationType                *string                                         `json:"control_automation_type"`
	ControlConsistencyID                 *uint                                           `json:"control_consistency_id"`
	ControlConsistencyTitleEn            *string                                         `json:"control_consistency_title_en"`
	ControlConsistencyTitleAr            *string                                         `json:"control_consistency_title_ar"`
	ControlConsistencyType               *string                                         `json:"control_consistency_type"`
	AssociatedDocuments                  []organizationDTOs.ControlAssociatedDocumentDTO `json:"associated_documents" gorm:"-"`
	LastTestDate                         *string                                         `json:"last_test_date"`
	ControlStatusID                      *uint                                           `json:"control_status_id"`
	ControlStatusTitleEn                 *string                                         `json:"control_status_title_en"`
	ControlStatusTitleAr                 *string                                         `json:"control_status_title_ar"`
	ControlStatusColorCode               *string                                         `json:"control_status_color_code"`
	ControlStatusType                    *string                                         `json:"control_status_type"`
	ControlTestTitle                     *string                                         `json:"control_test_title"`
	ControlTestDescription               *string                                         `json:"control_test_description"`
	ControlTestComments                  *string                                         `json:"control_test_comments"`
	ControlTestScore                     *uint                                           `json:"control_test_score"`
	ControlTestEffectivenessScore        *float64                                        `json:"control_test_effectiveness_score"`
	ControlTestRiskReduction             *float64                                        `json:"control_test_risk_reduction"`
	ControlTestTestedAt                  *string                                         `json:"control_test_tested_at"`
	ControlTestTestedByFirstNameEn       *string                                         `json:"control_test_tested_by_first_name_en"`
	ControlTestTestedByLastNameEn        *string                                         `json:"control_test_tested_by_last_name_en"`
	ControlTestTestedByInitial           *string                                         `json:"control_test_tested_by_initial"`
	ControlTestTestedByColor             *string                                         `json:"control_test_tested_by_color"`
	ControlTestDocuments                 []organizationDTOs.ControlTestDocumentDTO       `json:"control_test_documents" gorm:"-"`
	CreatedAt                            string                                          `json:"created_at"`
	CreatedUserID                        uint                                            `json:"created_user_id"`
	CreatedUserFirstNameEn               string                                          `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string                                          `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string                                          `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string                                          `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string                                          `json:"created_user_initial"`
	CreatedUserColor                     string                                          `json:"created_user_color"`
	CreatedUserEmail                     string                                          `json:"created_user_email"`
	CreatedUserProfileImagePath          string                                          `json:"created_user_profile_image_path"`
}

type RiskControlAnalysisListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
}

type RiskControlAnalysisFormDTO struct {
	ControlEfficiencyStatusID     uint `json:"control_efficiency_status_id" binding:"required"`
	ControlImplementationStatusID uint `json:"control_implementation_status_id" binding:"required"`
}

type RiskControlAnalysisDetailDTO struct {
	ID                          uint                                            `json:"id"`
	ControlTestID               *uint                                           `json:"control_test_id"`
	ControlEfficiencyStatus     organizationDTOs.ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	ControlImplementationStatus organizationDTOs.ControlImplementationStatusDTO `json:"control_implementation_status" `
	Score                       uint                                            `json:"score"`
	ControlEffectivenessScore   *float64                                        `json:"control_effectiveness_score"`
	ControlEffectivenessRating  *organizationDTOs.ControlEffectivenessRatingDTO `json:"control_effectiveness_rating"`
	RiskReduction               *float64                                        `json:"risk_reduction"`
}

type RiskResidualControlAnalysisDetailDTO struct {
	ID                          uint                                            `json:"id"`
	ControlTestID               *uint                                           `json:"control_test_id"`
	ControlEfficiencyStatus     organizationDTOs.ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	ControlImplementationStatus organizationDTOs.ControlImplementationStatusDTO `json:"control_implementation_status" `
	Score                       uint                                            `json:"score"`
	ControlEffectivenessScore   *float64                                        `json:"control_effectiveness_score"`
	ControlEffectivenessRating  *organizationDTOs.ControlEffectivenessRatingDTO `json:"control_effectiveness_rating"`
	RiskReduction               *float64                                        `json:"risk_reduction"`
}
