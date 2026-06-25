// awareness_plan_dashboard_repository.go

package repositories

import (
	awarenessDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"
	"math"
	"time"

	"gorm.io/gorm"
)

type AwarenessDashboardRepositoryInterface interface {
	GetAwarenessPlanSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPlanSummaryDTO, error)
	GetCampaignsByType(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsByTypeDTO, error)
	GetPriorityDistribution(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessPriorityDistributionDTO, error)
	GetCampaignsPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsPerMonthDTO, error)
	GetChannelMixPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessChannelMixPerMonthDTO, error)
	GetCoverageByAudience(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCoverageByAudienceDTO, error)

	GetAwarenessProgressSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessProgressSummaryDTO, error)
	GetCompletionByChannel(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByChannelDTO, error)
	GetCompletionPerCampaign(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionPerCampaignDTO, error)
	GetSurveyCompletionStatus(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessSurveyCompletionDTO, error)
	GetCompletionByDepartment(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByDepartmentDTO, error)
	GetPhishingResults(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPhishingResultsDTO, error)
}

type AwarenessDashboardRepository struct {
	db *gorm.DB
}

func NewAwarenessDashboardRepository(db *gorm.DB) *AwarenessDashboardRepository {
	return &AwarenessDashboardRepository{db: db}
}

func (repo *AwarenessDashboardRepository) applyDateFilter(query *gorm.DB, params *awarenessDTOs.AwarenessDashboardQueryParamsDTO, table string) *gorm.DB {
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(table+".created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	} else if params.Year > 0 {
		query = query.Where("YEAR("+table+".created_at) = ?", params.Year)
	}
	return query
}

func (repo *AwarenessDashboardRepository) GetAwarenessPlanSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPlanSummaryDTO, error) {
	summary := &awarenessDTOs.AwarenessPlanSummaryDTO{}

	campaignQuery := repo.db.Model(&models.AwarenessCampaign{})
	campaignQuery = repo.applyDateFilter(campaignQuery, params, "awareness_campaigns")
	if err := campaignQuery.Count(&summary.CampaignsPlanned).Error; err != nil {
		return nil, err
	}

	surveyQuery := repo.db.Model(&models.AwarenessCampaign{}).
		Joins("JOIN awareness_post_campaign_surveys ON awareness_post_campaign_surveys.linked_campaign_id = awareness_campaigns.id")
	surveyQuery = repo.applyDateFilter(surveyQuery, params, "awareness_campaigns")
	if err := surveyQuery.Count(&summary.CampaignsWithSurvey).Error; err != nil {
		return nil, err
	}

	if summary.CampaignsPlanned > 0 {
		summary.CampaignsWithSurveyPct = math.Round(float64(summary.CampaignsWithSurvey) / float64(summary.CampaignsPlanned) * 100)
	}

	if err := repo.db.Model(&models.AwarenessChannel{}).
		Where("status = ?", "active").
		Count(&summary.ActiveChannels).Error; err != nil {
		return nil, err
	}

	var plannedReach *int64
	reachQuery := repo.db.Model(&models.AwarenessCampaignDetail{}).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_campaign_details.awareness_campaign_id").
		Select("SUM(awareness_campaign_details.recipient_count)")
	reachQuery = repo.applyDateFilter(reachQuery, params, "awareness_campaigns")
	if err := reachQuery.Scan(&plannedReach).Error; err != nil {
		return nil, err
	}
	if plannedReach != nil {
		summary.PlannedReach = *plannedReach
	}

	return summary, nil
}

func (repo *AwarenessDashboardRepository) GetCampaignsByType(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsByTypeDTO, error) {
	var result []*awarenessDTOs.AwarenessCampaignsByTypeDTO

	query := repo.db.Table("awareness_campaigns").
		Select(`
			awareness_campaign_types.title_en AS type_title_en,
			awareness_campaign_types.title_ar AS type_title_ar,
			COUNT(awareness_campaigns.id) AS count
		`).
		Joins("JOIN awareness_campaign_types ON awareness_campaign_types.id = awareness_campaigns.awareness_campaign_type_id").
		Group("awareness_campaign_types.title_en, awareness_campaign_types.title_ar").
		Order("count DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetPriorityDistribution(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessPriorityDistributionDTO, error) {
	var result []*awarenessDTOs.AwarenessPriorityDistributionDTO

	var total int64
	totalQuery := repo.db.Model(&models.AwarenessCampaign{}).
		Where("awareness_priority_id IS NOT NULL")
	totalQuery = repo.applyDateFilter(totalQuery, params, "awareness_campaigns")
	if err := totalQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	query := repo.db.Table("awareness_campaigns").
		Select(`
			awareness_priorities.title_en AS priority_title_en,
			awareness_priorities.title_ar AS priority_title_ar,
			COUNT(awareness_campaigns.id) AS count,
			CAST(ROUND(COUNT(awareness_campaigns.id) * 100.0 / ?, 0) AS INT) AS percentage
		`, total).
		Joins("JOIN awareness_priorities ON awareness_priorities.id = awareness_campaigns.awareness_priority_id").
		Where("awareness_campaigns.awareness_priority_id IS NOT NULL").
		Group("awareness_priorities.title_en, awareness_priorities.title_ar").
		Order("count DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetCampaignsPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCampaignsPerMonthDTO, error) {
	var result []*awarenessDTOs.AwarenessCampaignsPerMonthDTO

	query := repo.db.Table("awareness_campaigns").
		Select("MONTH(awareness_campaigns.created_at) AS month, COUNT(awareness_campaigns.id) AS count").
		Group("MONTH(awareness_campaigns.created_at)").
		Order("month ASC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetChannelMixPerMonth(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessChannelMixPerMonthDTO, error) {
	var result []*awarenessDTOs.AwarenessChannelMixPerMonthDTO

	query := repo.db.Table("awareness_campaign_channels").
		Select(`
			MONTH(awareness_campaigns.created_at) AS month,
			awareness_channels.title_en AS channel_title_en,
			awareness_channels.title_ar AS channel_title_ar,
			COUNT(awareness_campaign_channels.id) AS count
		`).
		Joins("JOIN awareness_campaign_details ON awareness_campaign_details.id = awareness_campaign_channels.awareness_campaign_details_id").
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_campaign_details.awareness_campaign_id").
		Joins("JOIN awareness_channels ON awareness_channels.id = awareness_campaign_channels.awareness_channel_id").
		Group("MONTH(awareness_campaigns.created_at), awareness_channels.title_en, awareness_channels.title_ar").
		Order("month ASC, count DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetCoverageByAudience(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCoverageByAudienceDTO, error) {
	var result []*awarenessDTOs.AwarenessCoverageByAudienceDTO

	query := repo.db.Table("awareness_campaigns").
		Select("awareness_campaigns.target_audience, COUNT(awareness_campaigns.id) AS count").
		Group("awareness_campaigns.target_audience").
		Order("count DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetAwarenessProgressSummary(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessProgressSummaryDTO, error) {
	summary := &awarenessDTOs.AwarenessProgressSummaryDTO{}

	totalQuery := repo.db.Model(&models.AwarenessCampaign{})
	totalQuery = repo.applyDateFilter(totalQuery, params, "awareness_campaigns")
	if err := totalQuery.Count(&summary.TotalCampaigns).Error; err != nil {
		return nil, err
	}

	completedQuery := repo.db.Model(&models.AwarenessCampaign{}).
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Where("awareness_campaign_statuses.type = ?", "completed")
	completedQuery = repo.applyDateFilter(completedQuery, params, "awareness_campaigns")
	if err := completedQuery.Count(&summary.CompletedCampaigns).Error; err != nil {
		return nil, err
	}

	if summary.TotalCampaigns > 0 {
		summary.OverallCompletionPct = math.Round(float64(summary.CompletedCampaigns) / float64(summary.TotalCampaigns) * 100)
	}

	var totalWithSurvey, completedSurveys int64
	surveyTotalQuery := repo.db.Model(&models.AwarenessPostCampaignSurvey{}).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_post_campaign_surveys.linked_campaign_id")
	surveyTotalQuery = repo.applyDateFilter(surveyTotalQuery, params, "awareness_campaigns")
	if err := surveyTotalQuery.Count(&totalWithSurvey).Error; err != nil {
		return nil, err
	}

	surveyCompletedQuery := repo.db.Model(&models.AwarenessPostCampaignSurvey{}).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_post_campaign_surveys.linked_campaign_id").
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Where("awareness_campaign_statuses.type = ?", "completed")
	surveyCompletedQuery = repo.applyDateFilter(surveyCompletedQuery, params, "awareness_campaigns")
	if err := surveyCompletedQuery.Count(&completedSurveys).Error; err != nil {
		return nil, err
	}

	if totalWithSurvey > 0 {
		summary.SurveyCompletionPct = math.Round(float64(completedSurveys) / float64(totalWithSurvey) * 100)
	}

	var onTime int64
	now := time.Now()
	onTimeQuery := repo.db.Model(&models.AwarenessCampaign{}).
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Where("awareness_campaign_statuses.type = ?", "completed").
		Where("awareness_campaigns.end_date IS NOT NULL").
		Where("awareness_campaigns.end_date >= ?", now)
	onTimeQuery = repo.applyDateFilter(onTimeQuery, params, "awareness_campaigns")
	if err := onTimeQuery.Count(&onTime).Error; err != nil {
		return nil, err
	}

	if summary.CompletedCampaigns > 0 {
		summary.DeliveredOnTimePct = math.Round(float64(onTime) / float64(summary.CompletedCampaigns) * 100)
	}

	return summary, nil
}

func (repo *AwarenessDashboardRepository) GetCompletionByChannel(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByChannelDTO, error) {
	var result []*awarenessDTOs.AwarenessCompletionByChannelDTO

	query := repo.db.Table("awareness_campaign_channels").
		Select(`
			awareness_channels.title_en AS channel_title_en,
			awareness_channels.title_ar AS channel_title_ar,
			SUM(CASE WHEN awareness_campaign_statuses.type = 'completed' THEN 1 ELSE 0 END) AS completed,
			SUM(CASE WHEN awareness_campaign_statuses.type = 'submitted' THEN 1 ELSE 0 END) AS in_progress,
			SUM(CASE WHEN awareness_campaign_statuses.type NOT IN ('completed', 'submitted') THEN 1 ELSE 0 END) AS not_started
		`).
		Joins("JOIN awareness_campaign_details ON awareness_campaign_details.id = awareness_campaign_channels.awareness_campaign_details_id").
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_campaign_details.awareness_campaign_id").
		Joins("JOIN awareness_channels ON awareness_channels.id = awareness_campaign_channels.awareness_channel_id").
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Group("awareness_channels.title_en, awareness_channels.title_ar").
		Order("completed DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetCompletionPerCampaign(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionPerCampaignDTO, error) {
	var result []*awarenessDTOs.AwarenessCompletionPerCampaignDTO

	query := repo.db.Table("awareness_campaigns").
		Select(`
			awareness_campaigns.id AS campaign_id,
			awareness_campaigns.reference_number,
			awareness_campaigns.campaign_title,
			awareness_campaign_statuses.title_en AS status_title_en,
			awareness_campaign_statuses.title_ar AS status_title_ar,
			CASE
				WHEN awareness_campaign_statuses.type = 'completed' THEN 100
				WHEN awareness_campaign_statuses.type = 'submitted' THEN 50
				ELSE 0
			END AS completion_pct
		`).
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Order("awareness_campaigns.created_at ASC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetSurveyCompletionStatus(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessSurveyCompletionDTO, error) {
	summary := &awarenessDTOs.AwarenessSurveyCompletionDTO{}

	totalQuery := repo.db.Model(&models.AwarenessPostCampaignSurvey{}).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_post_campaign_surveys.linked_campaign_id")
	totalQuery = repo.applyDateFilter(totalQuery, params, "awareness_campaigns")
	if err := totalQuery.Count(&summary.TotalSurveys).Error; err != nil {
		return nil, err
	}

	completedQuery := repo.db.Model(&models.AwarenessPostCampaignSurvey{}).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_post_campaign_surveys.linked_campaign_id").
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Where("awareness_campaign_statuses.type = ?", "completed")
	completedQuery = repo.applyDateFilter(completedQuery, params, "awareness_campaigns")
	if err := completedQuery.Count(&summary.CompletedSurveys).Error; err != nil {
		return nil, err
	}

	if summary.TotalSurveys > 0 {
		summary.CompletionPct = math.Round(float64(summary.CompletedSurveys) / float64(summary.TotalSurveys) * 100)
	}

	return summary, nil
}

func (repo *AwarenessDashboardRepository) GetCompletionByDepartment(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) ([]*awarenessDTOs.AwarenessCompletionByDepartmentDTO, error) {
	var result []*awarenessDTOs.AwarenessCompletionByDepartmentDTO

	query := repo.db.Table("awareness_campaigns").
		Select(`
			organization_units.title_en AS department_title_en,
			organization_units.title_ar AS department_title_ar,
			COUNT(awareness_campaigns.id) AS campaign_count,
			CAST(ROUND(
				SUM(CASE WHEN awareness_campaign_statuses.type = 'completed' THEN 1.0 ELSE 0.0 END)
				/ COUNT(awareness_campaigns.id) * 100, 0
			) AS INT) AS completion_pct
		`).
		Joins("JOIN organization_units ON organization_units.id = awareness_campaigns.organization_unit_id").
		Joins("JOIN awareness_campaign_statuses ON awareness_campaign_statuses.id = awareness_campaigns.awareness_campaign_status_id").
		Group("organization_units.title_en, organization_units.title_ar").
		Order("completion_pct DESC")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *AwarenessDashboardRepository) GetPhishingResults(params *awarenessDTOs.AwarenessDashboardQueryParamsDTO) (*awarenessDTOs.AwarenessPhishingResultsDTO, error) {
	summary := &awarenessDTOs.AwarenessPhishingResultsDTO{}

	query := repo.db.Table("awareness_campaign_details").
		Select(`
			AVG(CAST(click_rate AS FLOAT)) AS avg_click_rate,
			AVG(CAST(data_submission_rate AS FLOAT)) AS avg_data_submission_rate,
			AVG(CAST(report_rate AS FLOAT)) AS avg_report_rate
		`).
		Joins("JOIN awareness_campaigns ON awareness_campaigns.id = awareness_campaign_details.awareness_campaign_id").
		Joins("JOIN awareness_campaign_types ON awareness_campaign_types.id = awareness_campaigns.awareness_campaign_type_id").
		Where("awareness_campaign_types.type = ?", "phishing-campaign").
		Where("awareness_campaign_details.click_rate IS NOT NULL")

	query = repo.applyDateFilter(query, params, "awareness_campaigns")

	if err := query.Scan(&summary).Error; err != nil {
		return nil, err
	}

	summary.AvgClickRate = math.Round(summary.AvgClickRate)
	summary.AvgDataSubmissionRate = math.Round(summary.AvgDataSubmissionRate)
	summary.AvgReportRate = math.Round(summary.AvgReportRate)

	return summary, nil
}
