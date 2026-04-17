package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type ControlListDTO struct {
	ID                                   uint       `json:"id"`
	ControlTypeTitleEn                   string     `json:"control_type_title_en"`
	ControlTypeTitleAr                   string     `json:"control_type_title_ar"`
	ControlCategoryTitleEn               string     `json:"control_category_title_en"`
	ControlCategoryTitleAr               string     `json:"control_category_title_ar"`
	OwnerUnitTitleEn                     *string    `json:"owner_unit_title_en"`
	OwnerUnitTitleAr                     *string    `json:"owner_unit_title_ar"`
	OwnerID                              *uint      `json:"owner_id"`
	OwnerFirstNameEn                     *string    `json:"owner_first_name_en"`
	OwnerFirstNameAr                     *string    `json:"owner_first_name_ar"`
	OwnerLastNameEn                      *string    `json:"owner_last_name_en"`
	OwnerLastNameAr                      *string    `json:"owner_last_name_ar"`
	OwnerInitial                         string     `json:"owner_initial"`
	OwnerColor                           string     `json:"owner_color"`
	OwnerEmail                           string     `json:"owner_email"`
	OwnerProfileImagePath                string     `json:"owner_profile_image_path"`
	ControlImplementationStatusTitleEn   string     `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string     `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string     `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string     `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string     `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string     `json:"control_efficiency_status_color_code"`
	ReferenceNumber                      string     `json:"reference_number"`
	Code                                 *string    `json:"code"`
	TitleEn                              string     `json:"title_en"`
	TitleAr                              string     `json:"title_ar"`
	Description                          *string    `json:"description"`
	Objective                            *string    `json:"objective"`
	ReviewFrequencyTitleEn               string     `json:"review_frequency_title_en"`
	ReviewFrequencyTitleAr               string     `json:"review_frequency_title_ar"`
	LastReviewDate                       *time.Time `json:"last_review_date"`
	NextReviewDate                       *time.Time `json:"next_review_date"`
	LastReviewNote                       *string    `json:"last_review_note"`
	TestResponsibleUserID                uint       `json:"test_responsible_user_id"`
	TestResponsibleUserFirstNameEn       string     `json:"test_responsible_user_first_name_en"`
	TestResponsibleUserFirstNameAr       string     `json:"test_responsible_user_first_name_ar"`
	TestResponsibleUserLastNameEn        string     `json:"test_responsible_user_last_name_en"`
	TestResponsibleUserLastNameAr        string     `json:"test_responsible_user_last_name_ar"`
	TestResponsibleUserInitial           string     `json:"test_responsible_user_initial"`
	TestResponsibleUserColor             string     `json:"test_responsible_user_color"`
	TestResponsibleUserEmail             string     `json:"test_responsible_user_email"`
	TestResponsibleUserProfileImagePath  string     `json:"test_responsible_user_profile_image_path"`
	TestFrequencyTitleEn                 *string    `json:"test_frequency_title_en"`
	TestFrequencyTitleAr                 *string    `json:"test_frequency_title_ar"`
	LastTestDate                         *time.Time `json:"last_test_date"`
	NextTestDate                         *time.Time `json:"next_test_date"`
	ImplementedAt                        *time.Time `json:"implemented_at"`
	IsApplicable                         bool       `json:"is_applicable"`
	Justification                        *string    `json:"justification"`
	ImplementationGuideline              *string    `json:"implementation_guideline"`
	ApplicabilityChangeDate              *string    `json:"applicability_change_date"`
	LifecycleStatusTitleEn               string     `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr               string     `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode             string     `json:"lifecycle_status_color_code"`
	CreatedUserID                        uint       `json:"created_user_id"`
	CreatedUserFirstNameEn               string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr               string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn                string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr                string     `json:"created_user_last_name_ar"`
	CreatedUserInitial                   string     `json:"created_user_initial"`
	CreatedUserColor                     string     `json:"created_user_color"`
	CreatedUserEmail                     string     `json:"created_user_email"`
	CreatedUserProfileImagePath          string     `json:"created_user_profile_image_path"`
	UpdatedUserID                        *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   string     `json:"updated_user_initial"`
	UpdatedUserColor                     string     `json:"updated_user_color"`
	UpdatedUserEmail                     string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath          string     `json:"updated_user_profile_image_path"`
	CreatedAt                            time.Time  `json:"created_at"`
	UpdatedAt                            *time.Time `json:"updated_at"`
	RiskRatingTitleEn                    *string    `json:"risk_rating_title_en"`
	RiskRatingTitleAr                    *string    `json:"risk_rating_title_ar"`
	RiskRatingColorCode                  *string    `json:"risk_rating_color_code"`
	RiskRatingValue                      *uint      `json:"risk_rating_value"`
	ControlAutomationID                  *uint      `json:"control_automation_id"`
	ControlAutomationTitleEn             *string    `json:"control_automation_title_en"`
	ControlAutomationTitleAr             *string    `json:"control_automation_title_ar"`
	ControlAutomationType                *string    `json:"control_automation_type"`
	ControlAutomationScore               *float64   `json:"control_automation_score"`
	ControlAutomationStatus              *string    `json:"control_automation_status"`
	ControlConsistencyID                 *uint      `json:"control_consistency_id"`
	ControlConsistencyTitleEn            *string    `json:"control_consistency_title_en"`
	ControlConsistencyTitleAr            *string    `json:"control_consistency_title_ar"`
	ControlConsistencyType               *string    `json:"control_consistency_type"`
	ControlConsistencyScore              *float64   `json:"control_consistency_score"`
	ControlConsistencyStatus             *string    `json:"control_consistency_status"`
	ControlConsistencyDescriptionEn      *string    `json:"control_consistency_description_en"`
	ControlConsistencyDescriptionAr      *string    `json:"control_consistency_description_ar"`
	ControlEffectivenessScore            *float64   `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID         *uint      `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn    *string    `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    *string    `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingMinScore   *float64   `json:"control_effectiveness_rating_min_score"`
	ControlEffectivenessRatingMaxScore   *float64   `json:"control_effectiveness_rating_max_score"`
	ControlEffectivenessRatingType       *string    `json:"control_effectiveness_rating_type"`
	RiskReduction                        *float64   `json:"risk_reduction"`
}

type ControlMinimalListDTO struct {
	ID              uint    `json:"id"`
	ReferenceNumber string  `json:"reference_number"`
	Code            *string `json:"code"`
	TitleEn         string  `json:"title_en"`
	TitleAr         string  `json:"title_ar"`
	Description     *string `json:"description"`
}

type ControlListQueryParamsDTO struct {
	Sort                  string
	SortBy                string
	Keyword               string
	Status                string
	Page                  int
	PageSize              int
	Filters               map[string][]string
	QuickFilter           []string
	LastReviewDateFrom    string
	LastReviewDateTo      string
	NextReviewDateFrom    string
	NextReviewDateTo      string
	LastTestDateFrom      string
	LastTestDateTo        string
	NextTestDateFrom      string
	NextTestDateTo        string
	ImplementedAtFrom     string
	ImplementedAtTo       string
	Export                bool
	DateFrom              string
	DateTo                string
	ComplianceRegisterIDs []int
}

type ControlListResponseDTO struct {
	Items      []*ControlListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlFormDTO struct {
	ControlTypeID                 uint    `json:"control_type_id" binding:"required"`
	ControlCategoryID             *uint   `json:"control_category_id" `
	OwnerID                       uint    `json:"owner_id" binding:"required"`
	ControlImplementationStatusID *uint   `json:"control_implementation_status_id" binding:"omitempty"`
	ControlEfficiencyStatusID     *uint   `json:"control_efficiency_status_id" binding:"omitempty"`
	Code                          *string `json:"code" binding:"omitempty"`
	TitleEn                       string  `json:"title_en" binding:"required"`
	TitleAr                       string  `json:"title_ar" binding:"required"`
	Description                   *string `json:"description" binding:"omitempty"`
	Objective                     *string `json:"objective" binding:"omitempty"`
	EfficiencyScore               *uint   `json:"efficiency_score" binding:"omitempty"`
	ReviewFrequencyID             *uint   `json:"review_frequency_id" binding:"omitempty"`
	LastReviewDate                *string `json:"last_review_date" binding:"omitempty,datetime=2006-01-02"`
	NextReviewDate                *string `json:"next_review_date" binding:"omitempty,datetime=2006-01-02"`
	LastReviewNote                *string `json:"last_review_note" binding:"omitempty"`
	TestResponsibleUserID         *uint   `json:"test_responsible_user_id"`
	TestFrequencyID               *uint   `json:"test_frequency_id,omitempty"`
	LastTestDate                  *string `json:"last_test_date" binding:"omitempty,datetime=2006-01-02"`
	NextTestDate                  *string `json:"next_test_date" binding:"omitempty,datetime=2006-01-02"`
	ImplementedAt                 *string `json:"implemented_at" binding:"omitempty,datetime=2006-01-02"`
	IsApplicable                  bool    `json:"is_applicable"`
	Justification                 *string `json:"justification" binding:"omitempty"`
	ApplicabilityChangeDate       *string `json:"applicability_change_date" binding:"omitempty,datetime=2006-01-02"`
	ImplementationGuideline       *string `json:"implementation_guideline" binding:"omitempty"`
	ControlAutomationID           *uint   `json:"control_automation_id" binding:"omitempty"`
	ControlConsistencyID          *uint   `json:"control_consistency_id" binding:"omitempty"`
}

type ControlParsedDates struct {
	LastReviewDate          *time.Time
	NextReviewDate          *time.Time
	LastTestDate            *time.Time
	NextTestDate            *time.Time
	ImplementedAt           *time.Time
	ApplicabilityChangeDate *time.Time
}

type ControlDetailDTO struct {
	ID                            uint                             `json:"id"`
	ReferenceNumber               string                           `json:"reference_number"`
	Code                          string                           `json:"code"`
	TitleEn                       string                           `json:"title_en"`
	TitleAr                       string                           `json:"title_ar"`
	Description                   *string                          `json:"description"`
	Objective                     *string                          `json:"objective"`
	ControlType                   ControlTypeDTO                   `json:"control_type"`
	ControlCategory               *ControlCategoryDTO              `json:"control_category"`
	ResponsibleUnit               settingsDTOs.OrganizationUnitDTO `json:"owner_unit"`
	ResponsibleUserID             *uint                            `json:"owner_id"`
	Responsibleuser               *UserMinimalDTO                  `json:"owner"`
	ControlImplementationStatus   *ControlImplementationStatusDTO  `json:"control_implementation_status"`
	ControlEfficiencyStatus       *ControlEfficiencyStatusDTO      `json:"control_efficiency_status"`
	EfficiencyScore               *uint                            `json:"efficiency_score"`
	ReviewFrequency               *masterDTOs.ReviewFrequencyDTO   `json:"review_frequency"`
	LastReviewDateString          *string                          `json:"last_review_date"`
	NextReviewDateString          *string                          `json:"next_review_date"`
	LastReviewNote                *string                          `json:"last_review_note"`
	TestResponsibleUser           UserMinimalDTO                   `json:"test_responsible_user"`
	TestFrequency                 *masterDTOs.ReviewFrequencyDTO   `json:"test_frequency"`
	LastTestDateString            *string                          `json:"last_test_date"`
	NextTestDateString            *string                          `json:"next_test_date"`
	ImplementedAtString           *string                          `json:"implemented_at"`
	CreatedAt                     time.Time                        `json:"created_at"`
	UpdatedAt                     *time.Time                       `json:"updated_at"`
	CreatedByUser                 UserMinimalDTO                   `json:"created_by"`
	UpdatedByUser                 *UserMinimalDTO                  `json:"updated_by"`
	LifecycleStatus               masterDTOs.LifecycleStatusDTO    `json:"status"`
	ControlReviewStatus           ControlReviewStatusDTO           `json:"review_status"`
	LatestControlTest             ControlTestDetailsDTO            `json:"latest_control_test"`
	LatestControlReview           *ControlReviewDTO                `json:"control_review"`
	MarkedAsAuditableItem         bool                             `json:"marked_as_auditable_item"`
	Risks                         []*ControlRiskDTO                `json:"risks"`
	RiskMatrixRating              *RiskMatrixRiskRatingDTO         `json:"risk_matrix_risk_rating"`
	RiskMatrixRatingValue         *RiskMatrixRiskRatingValueDTO    `json:"risk_matrix_rating_value"`
	IsApplicable                  bool                             `json:"is_applicable"`
	Justification                 *string                          `json:"justification"`
	ImplementationGuideline       *string                          `json:"implementation_guideline"`
	ApplicabilityChangeDateString *string                          `json:"applicability_change_date"`
	ComplianceStatusID            uint                             `json:"compliance_status_id"`
	ComplianceStatus              *ComplianceStatusDTO             `json:"compliance_status"`
	ControlAutomationID           *uint                            `json:"control_automation_id"`
	ControlConsistencyID          *uint                            `json:"control_consistency_id"`
	ControlAutomation             *ControlAutomationDetailDTO      `json:"control_automation"`
	ControlConsistency            *ControlConsistencyDetailDTO     `json:"control_consistency"`
	AssociatedDocuments           []ControlAssociatedDocumentDTO   `json:"associated_documents"`
	ControlEffectivenessScore     *float64                         `json:"control_effectiveness_score"`
	ControlEffectivenessRating    *ControlEffectivenessRatingDTO   `json:"control_effectiveness_rating"`
	RiskReduction                 *float64                         `json:"risk_reduction"`
}
type ComplianceStatusDTO struct {
	ID        int    `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	Type      string `json:"type"`
	ColorCode string `json:"color_code"`
	Status    string `json:"status"`
	ScoreFrom *int   `json:"score_from"`
	ScoreTo   *int   `json:"score_to"`
}

type ControlRiskDTO struct {
	ID             *uint          `json:"id"`
	ControlID      uint           `json:"control_id"`
	RiskRegisterID uint           `json:"risk_register_id"`
	Risk           RiskMinimalDTO `json:"risk"`
}
type ControlReviewDTO struct {
	ID             uint                          `json:"id"`
	Comment        string                        `json:"comment"`
	ReviewedAt     time.Time                     `json:"reviewed_at"`
	UpdatedAt      *time.Time                    `json:"updated_at"`
	ReviewedByUser UserMinimalDTO                `json:"reviewed_by"`
	UpdatedByUser  *UserMinimalDTO               `json:"updated_by"`
	Documents      []*commonDTOs.DocumentFormDTO `json:"documents"`
}

type ControlMailDTO struct {
	FirstNameEn string `json:"first_name_en"`
	LastNameEn  string `json:"last_name_en"`
	Subject     string `json:"subject"`
	MessageEn   string `json:"message_en"`
	RedirectURL string `json:"redirect_url"`
}

type ControlAutomationDetailDTO struct {
	ID            uint            `json:"id"`
	TitleEn       string          `json:"title_en"`
	TitleAr       string          `json:"title_ar"`
	Type          string          `json:"type"`
	Score         float64         `json:"score"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedByID   uint            `json:"created_by_id"`
	CreatedByUser *UserMinimalDTO `json:"created_by_user"`
}

type ControlConsistencyDetailDTO struct {
	ID            uint            `json:"id"`
	TitleEn       string          `json:"title_en"`
	TitleAr       string          `json:"title_ar"`
	Type          string          `json:"type"`
	Score         float64         `json:"score"`
	Status        string          `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedByID   uint            `json:"created_by_id"`
	CreatedByUser *UserMinimalDTO `json:"created_by_user"`
}

type ControlApplicableFormDTO struct {
	Justification           *string `json:"justification" binding:"omitempty"`
	ApplicabilityChangeDate *string `json:"applicability_change_date" binding:"omitempty,datetime=2006-01-02"`
}

type ControlAssociatedDocumentDTO struct {
	ID         uint         `json:"id"`
	ControlID  uint         `json:"control_id"`
	DocumentID uint         `json:"document_id"`
	Document   *DocumentDTO `json:"document"`
}

type DocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
