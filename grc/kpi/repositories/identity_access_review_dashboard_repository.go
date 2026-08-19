package repositories

import (
	identityAccessReviewDashboardDtos "jethur/grc/kpi/dtos"

	"gorm.io/gorm"
)

type IdentityAccessReviewDashboardRepositoryInterface interface {
	GetIdentityAccessReviewSummary() (*identityAccessReviewDashboardDtos.IdentityAccessReviewSummaryDTO, error)
	GetCampaignsByLifecycle() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCampaignByLifecycleDTO, error)
	GetCompletionBySystem() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCompletionBySystemDTO, error)
	GetActiveAndRecentCampaigns() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewRecentCampaignDTO, error)
}

type IdentityAccessReviewDashboardRepository struct {
	db *gorm.DB
}

func NewIdentityAccessReviewDashboardRepository(db *gorm.DB) IdentityAccessReviewDashboardRepositoryInterface {
	return &IdentityAccessReviewDashboardRepository{
		db: db,
	}
}

func (repo *IdentityAccessReviewDashboardRepository) GetIdentityAccessReviewSummary() (*identityAccessReviewDashboardDtos.IdentityAccessReviewSummaryDTO, error) {
	result := &identityAccessReviewDashboardDtos.IdentityAccessReviewSummaryDTO{}

	err := repo.db.Table("identity_review_campaigns AS irc").
		Select(`
			COALESCE(
				ROUND(
					CAST(SUM(CASE WHEN rcs.type = 'approved' THEN 1 ELSE 0 END) AS float)
					* 100.0 / NULLIF(COUNT(irc.id), 0),
					2
				),
				0
			) AS overall_completion_percentage,
			COALESCE(SUM(CASE
				WHEN irc.due_date < CAST(GETDATE() AS date)
					AND rcs.type <> 'approved' THEN 1
				ELSE 0
			END), 0) AS overdue_campaigns_count,
			COALESCE(SUM(CASE WHEN rcs.type = 'approved' THEN 1 ELSE 0 END), 0) AS approved_campaigns_count,
			COUNT(irc.id) AS total_campaigns_count,
			(
				SELECT COUNT(irap.id)
				FROM identity_review_action_plans AS irap
				JOIN action_plan_statuses AS aps ON aps.id = irap.action_plan_status_id
				WHERE aps.is_closed = 0
			) AS open_corrective_actions_count
		`).
		Joins("JOIN review_campaign_statuses AS rcs ON rcs.id = irc.campaign_status_id").
		Scan(result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *IdentityAccessReviewDashboardRepository) GetCampaignsByLifecycle() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCampaignByLifecycleDTO, error) {
	var result []*identityAccessReviewDashboardDtos.IdentityAccessReviewCampaignByLifecycleDTO

	err := repo.db.Table("review_campaign_statuses AS rcs").
		Select(`
			rcs.id AS campaign_status_id,
			rcs.title_en,
			rcs.title_ar,
			rcs.color_code,
			rcs.type,
			COUNT(irc.id) AS count
		`).
		Joins("LEFT JOIN identity_review_campaigns AS irc ON irc.campaign_status_id = rcs.id").
		Where("rcs.type <> ?", "closed").
		Group("rcs.id, rcs.title_en, rcs.title_ar, rcs.color_code, rcs.type").
		Order("rcs.id ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *IdentityAccessReviewDashboardRepository) GetCompletionBySystem() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewCompletionBySystemDTO, error) {
	var result []*identityAccessReviewDashboardDtos.IdentityAccessReviewCompletionBySystemDTO

	err := repo.db.Table("identity_review_campaign_systems AS ircs").
		Select(`
			COALESCE(ar.id, 0) AS asset_register_id,
			COALESCE(ar.title_en, '') AS asset_title_en,
			COALESCE(ar.title_ar, '') AS asset_title_ar,
			COALESCE(ROUND(AVG(CAST(irc.completion_percentage AS float)), 2), 0) AS completion_percentage
		`).
		Joins("JOIN identity_review_campaigns AS irc ON irc.id = ircs.identity_review_campaign_id").
		Joins("JOIN asset_register AS ar ON ar.id = ircs.asset_register_id").
		Where("ircs.asset_register_id IS NOT NULL").
		Group("ar.id, ar.title_en, ar.title_ar").
		Order("completion_percentage DESC, ar.title_en ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *IdentityAccessReviewDashboardRepository) GetActiveAndRecentCampaigns() ([]*identityAccessReviewDashboardDtos.IdentityAccessReviewRecentCampaignDTO, error) {
	var result []*identityAccessReviewDashboardDtos.IdentityAccessReviewRecentCampaignDTO

	err := repo.db.Table("identity_review_campaigns AS irc").
		Select(`
			irc.id,
			irc.reference_number,
			irc.campaign_name,
			COALESCE(ar.id, 0) AS asset_register_id,
			COALESCE(ar.title_en, '') AS asset_title_en,
			COALESCE(ar.title_ar, '') AS asset_title_ar,
			rt.id AS review_type_id,
			rt.title_en AS review_type_title_en,
			rt.title_ar AS review_type_title_ar,
			rt.type AS review_type,
			CONVERT(VARCHAR, irc.start_date, 23) AS start_date,
			CONVERT(VARCHAR, irc.due_date, 23) AS due_date,
			reviewer.id AS reviewer_id,
			reviewer.first_name_en AS reviewer_first_name_en,
			reviewer.first_name_ar AS reviewer_first_name_ar,
			reviewer.last_name_en AS reviewer_last_name_en,
			reviewer.last_name_ar AS reviewer_last_name_ar,
			reviewer.initial AS reviewer_initial,
			rcs.id AS campaign_status_id,
			rcs.title_en AS campaign_status_title_en,
			rcs.title_ar AS campaign_status_title_ar,
			rcs.color_code AS campaign_status_color_code,
			rcs.type AS campaign_status_type,
			irc.completion_percentage,
			irc.created_at
		`).
		Joins("JOIN review_types AS rt ON rt.id = irc.review_type_id").
		Joins("JOIN review_campaign_statuses AS rcs ON rcs.id = irc.campaign_status_id").
		Joins("JOIN users AS reviewer ON reviewer.id = irc.reviewer_id").
		Joins("LEFT JOIN identity_review_campaign_systems AS ircs ON ircs.identity_review_campaign_id = irc.id").
		Joins("LEFT JOIN asset_register AS ar ON ar.id = ircs.asset_register_id").
		Where("rcs.type <> ?", "closed").
		Order("irc.created_at DESC, irc.id DESC").
		Limit(10).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
