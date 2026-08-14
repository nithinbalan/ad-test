package dtos

type IdentityReviewCampaignExportPDFDTO struct {
	CampaignName                  string
	ReferenceNumber               string
	Description                   *string
	StartDate                     string
	DueDate                       string
	GeneratedAt                   string
	ReviewTypeTitleEn             string
	ReminderDaysOffset            string
	AssetTitleEn                  string
	CampaignStatusTitleEn         string
	CampaignStatusTextColor       string
	CampaignStatusBorderColor     string
	CampaignStatusBackgroundColor string
	ReviewerName                  string
	CompletionPercentage          int
	CompletionProgressDash        float64
	Users                         []*IdentityReviewCampaignExportUserDTO
	Workflow                      *IdentityReviewCampaignWorkflowListDTO
}

type IdentityReviewCampaignExportUserDTO struct {
	ID                             uint
	Username                       string
	UserIdentifier                 *string
	LastLogin                      *string
	CurrentRole                    *string
	DepartmentTitleEn              string
	AccountTypeTitleEn             *string
	PrivilegeTitleEn               *string
	AccountDecisionID              *uint
	AccountDecisionTitleEn         *string
	AccountDecisionColorCode       *string
	AccountDecisionType            *string
	AccountDecisionTextColor       string
	AccountDecisionBorderColor     string
	AccountDecisionBackgroundColor string
}
