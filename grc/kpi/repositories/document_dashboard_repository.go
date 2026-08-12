package repositories

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"math"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type DocumentDashboardRepositoryInterface interface {
	GetDocumentSummary(params *kpiDTOs.DocumentDashboardQueryParamsDTO) (*kpiDTOs.DocumentDashboardSummaryDTO, error)
	GetDocumentStatusByType(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentStatusByTypeDTO, error)
	GetDocumentsByConfidentialityClassification(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentConfidentialityClassificationDTO, error)
	GetDocumentsByResponsibleUnit(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentResponsibleUnitDTO, error)
	GetDocumentReviewCountTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentReviewCountTrendDTO, error)
}

type DocumentDashboardRepository struct {
	db *gorm.DB
}

type documentReviewTimelinessCounts struct {
	CurrentPeriodReviews       int64
	CurrentPeriodOnTimeReviews int64
	CurrentThirtyDayReviews    int64
	CurrentThirtyDayOnTime     int64
	PreviousThirtyDayReviews   int64
	PreviousThirtyDayOnTime    int64
}

type documentThirtyDayComparisonCounts struct {
	CurrentTotal      int64
	PreviousTotal     int64
	CurrentApproved   int64
	PreviousApproved  int64
	CurrentPublished  int64
	PreviousPublished int64
}

type documentApprovalTimeResult struct {
	AverageApprovalTimeDays         float64
	CurrentAverageApprovalTimeDays  *float64
	PreviousAverageApprovalTimeDays *float64
}

type complianceRegisterLegislation struct {
	ID            uint
	LegislationNo string
}

type documentMonthlyReviewCounts struct {
	Month         time.Time
	DocumentCount int64
}

var eccLegislationNumberPattern = regexp.MustCompile(`(?i)(^|[^a-z0-9])ECC([^a-z0-9]|$)`)

func NewDocumentDashboardRepository(db *gorm.DB) DocumentDashboardRepositoryInterface {
	return &DocumentDashboardRepository{db: db}
}

func (repo *DocumentDashboardRepository) GetDocumentSummary(params *kpiDTOs.DocumentDashboardQueryParamsDTO) (*kpiDTOs.DocumentDashboardSummaryDTO, error) {
	result := &kpiDTOs.DocumentDashboardSummaryDTO{}
	documentCategoryIDs := params.Filters["document_category_id"]

	documentCountsQuery := repo.db.Table("documents d").
		Select(`
			COUNT(d.id) AS total_documents,
			COALESCE(SUM(CASE WHEN ds.type IN ('approved', 'published') THEN 1 ELSE 0 END), 0) AS approved_documents,
			COALESCE(SUM(CASE WHEN ds.type = 'published' THEN 1 ELSE 0 END), 0) AS published_documents
		`).
		Joins("LEFT JOIN document_statuses ds ON ds.id = d.document_status_id").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false)
	if len(documentCategoryIDs) > 0 {
		documentCountsQuery = documentCountsQuery.Where("d.document_category_id IN ?", documentCategoryIDs)
	}
	err := documentCountsQuery.Scan(result).Error
	if err != nil {
		return nil, err
	}

	comparisonCounts := &documentThirtyDayComparisonCounts{}
	comparisonQuery := repo.db.Table("documents d").
		Select(`
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, 1, CAST(GETDATE() AS date))
				THEN 1 ELSE 0 END), 0) AS current_total,
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, -29, CAST(GETDATE() AS date))
				THEN 1 ELSE 0 END), 0) AS previous_total,
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, 1, CAST(GETDATE() AS date))
					AND ds.type IN ('approved', 'published')
				THEN 1 ELSE 0 END), 0) AS current_approved,
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, -29, CAST(GETDATE() AS date))
					AND ds.type IN ('approved', 'published')
				THEN 1 ELSE 0 END), 0) AS previous_approved,
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, 1, CAST(GETDATE() AS date))
					AND ds.type = 'published'
				THEN 1 ELSE 0 END), 0) AS current_published,
			COALESCE(SUM(CASE
				WHEN d.created_at >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND d.created_at < DATEADD(day, -29, CAST(GETDATE() AS date))
					AND ds.type = 'published'
				THEN 1 ELSE 0 END), 0) AS previous_published
		`).
		Joins("LEFT JOIN document_statuses ds ON ds.id = d.document_status_id").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false)
	if len(documentCategoryIDs) > 0 {
		comparisonQuery = comparisonQuery.Where("d.document_category_id IN ?", documentCategoryIDs)
	}
	err = comparisonQuery.Scan(comparisonCounts).Error
	if err != nil {
		return nil, err
	}

	result.TotalDocumentsChangePercentage = calculatePercentageChange(
		comparisonCounts.CurrentTotal,
		comparisonCounts.PreviousTotal,
	)
	result.ApprovedDocumentsChangePercentage = calculatePercentageChange(
		comparisonCounts.CurrentApproved,
		comparisonCounts.PreviousApproved,
	)
	result.PublishedDocumentsChangePercentage = calculatePercentageChange(
		comparisonCounts.CurrentPublished,
		comparisonCounts.PreviousPublished,
	)

	reviewCounts := &documentReviewTimelinessCounts{}
	reviewCountsQuery := repo.db.Table("document_reviews dr").
		Select(`
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(month, -12, CAST(GETDATE() AS date))
					AND dr.review_date <= CAST(GETDATE() AS date)
					AND d.next_review_date IS NOT NULL
				THEN 1 ELSE 0 END), 0) AS current_period_reviews,
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(month, -12, CAST(GETDATE() AS date))
					AND dr.review_date <= CAST(GETDATE() AS date)
					AND d.next_review_date IS NOT NULL
					AND dr.review_date <= d.next_review_date
				THEN 1 ELSE 0 END), 0) AS current_period_on_time_reviews,
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND dr.review_date < DATEADD(day, 1, CAST(GETDATE() AS date))
					AND d.next_review_date IS NOT NULL
				THEN 1 ELSE 0 END), 0) AS current_thirty_day_reviews,
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND dr.review_date < DATEADD(day, 1, CAST(GETDATE() AS date))
					AND d.next_review_date IS NOT NULL
					AND dr.review_date <= d.next_review_date
				THEN 1 ELSE 0 END), 0) AS current_thirty_day_on_time,
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND dr.review_date < DATEADD(day, -29, CAST(GETDATE() AS date))
					AND d.next_review_date IS NOT NULL
				THEN 1 ELSE 0 END), 0) AS previous_thirty_day_reviews,
			COALESCE(SUM(CASE
				WHEN dr.review_date >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND dr.review_date < DATEADD(day, -29, CAST(GETDATE() AS date))
					AND d.next_review_date IS NOT NULL
					AND dr.review_date <= d.next_review_date
				THEN 1 ELSE 0 END), 0) AS previous_thirty_day_on_time
		`).
		Joins("INNER JOIN documents d ON d.id = dr.document_id").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false)
	if len(documentCategoryIDs) > 0 {
		reviewCountsQuery = reviewCountsQuery.Where("d.document_category_id IN ?", documentCategoryIDs)
	}
	err = reviewCountsQuery.Scan(reviewCounts).Error
	if err != nil {
		return nil, err
	}

	if reviewCounts.CurrentPeriodReviews > 0 {
		result.ReviewedOnTimePercentage = roundToOneDecimal(
			float64(reviewCounts.CurrentPeriodOnTimeReviews) * 100 / float64(reviewCounts.CurrentPeriodReviews),
		)
	}
	if reviewCounts.CurrentThirtyDayReviews > 0 && reviewCounts.PreviousThirtyDayReviews > 0 {
		currentOnTimePercentage := float64(reviewCounts.CurrentThirtyDayOnTime) * 100 / float64(reviewCounts.CurrentThirtyDayReviews)
		previousOnTimePercentage := float64(reviewCounts.PreviousThirtyDayOnTime) * 100 / float64(reviewCounts.PreviousThirtyDayReviews)
		result.ReviewedOnTimeChangePercentage = calculateMetricPercentageChange(
			currentOnTimePercentage,
			previousOnTimePercentage,
		)
	}

	approvalTime := &documentApprovalTimeResult{}
	approvalTimeQuery := repo.db.Table("document_versions dv").
		Select(`
			COALESCE(
				AVG(CAST(DATEDIFF(second, dv.submitted_at, final_approval.final_approved_at) AS FLOAT) / 86400.0),
				0
			) AS average_approval_time_days,
			AVG(CASE
				WHEN final_approval.final_approved_at >= DATEADD(day, -29, CAST(GETDATE() AS date))
					AND final_approval.final_approved_at < DATEADD(day, 1, CAST(GETDATE() AS date))
				THEN CAST(DATEDIFF(second, dv.submitted_at, final_approval.final_approved_at) AS FLOAT) / 86400.0
			END) AS current_average_approval_time_days,
			AVG(CASE
				WHEN final_approval.final_approved_at >= DATEADD(day, -59, CAST(GETDATE() AS date))
					AND final_approval.final_approved_at < DATEADD(day, -29, CAST(GETDATE() AS date))
				THEN CAST(DATEDIFF(second, dv.submitted_at, final_approval.final_approved_at) AS FLOAT) / 86400.0
			END) AS previous_average_approval_time_days
		`).
		Joins("INNER JOIN documents d ON d.id = dv.document_id").
		Joins(`INNER JOIN (
			SELECT
				dw.document_version_id,
				MAX(dwr.created_at) AS final_approved_at
			FROM document_workflows dw
			INNER JOIN document_workflow_statuses current_status
				ON current_status.id = dw.document_workflow_status_id
				AND current_status.type = 'approve'
			INNER JOIN document_workflow_reviews dwr
				ON dwr.document_workflow_id = dw.id
			INNER JOIN document_workflow_statuses review_status
				ON review_status.id = dwr.document_workflow_status_id
				AND review_status.type = 'approve'
			WHERE dw.level = (
				SELECT MAX(last_workflow.level)
				FROM document_workflows last_workflow
				WHERE last_workflow.document_version_id = dw.document_version_id
			)
			GROUP BY dw.document_version_id
		) final_approval ON final_approval.document_version_id = dv.id`).
		Where("dv.is_latest_version = ?", true).
		Where("dv.submitted_at IS NOT NULL").
		Where("final_approval.final_approved_at >= dv.submitted_at").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false)
	if len(documentCategoryIDs) > 0 {
		approvalTimeQuery = approvalTimeQuery.Where("d.document_category_id IN ?", documentCategoryIDs)
	}
	err = approvalTimeQuery.Scan(approvalTime).Error
	if err != nil {
		return nil, err
	}

	result.AverageApprovalTimeDays = roundToOneDecimal(approvalTime.AverageApprovalTimeDays)
	if approvalTime.CurrentAverageApprovalTimeDays != nil && approvalTime.PreviousAverageApprovalTimeDays != nil {
		result.AverageApprovalTimeChangePercentage = calculateMetricPercentageChange(
			*approvalTime.CurrentAverageApprovalTimeDays,
			*approvalTime.PreviousAverageApprovalTimeDays,
		)
	}

	complianceRegisters := []complianceRegisterLegislation{}
	err = repo.db.Table("compliance_register").
		Select("id, legislation_no").
		Where("is_latest = ? AND legislation_no IS NOT NULL", true).
		Scan(&complianceRegisters).Error
	if err != nil {
		return nil, err
	}

	eccComplianceRegisterIDs := make([]uint, 0)
	for _, complianceRegister := range complianceRegisters {
		if eccLegislationNumberPattern.MatchString(complianceRegister.LegislationNo) {
			eccComplianceRegisterIDs = append(eccComplianceRegisterIDs, complianceRegister.ID)
		}
	}

	if len(eccComplianceRegisterIDs) > 0 {
		var totalECCControlCount int64
		err = repo.db.Table("compliance_controls cc").
			Where("cc.compliance_register_id IN ?", eccComplianceRegisterIDs).
			Distinct("cc.control_id").
			Count(&totalECCControlCount).Error
		if err != nil {
			return nil, err
		}

		var documentedECCControlCount int64
		documentedControlsQuery := repo.db.Table("compliance_controls cc").
			Joins("INNER JOIN document_controls dc ON dc.control_id = cc.control_id").
			Joins("INNER JOIN documents d ON d.id = dc.document_id").
			Where("cc.compliance_register_id IN ?", eccComplianceRegisterIDs).
			Where("d.is_internal = ? AND d.is_folder = ?", false, false)
		if len(documentCategoryIDs) > 0 {
			documentedControlsQuery = documentedControlsQuery.Where("d.document_category_id IN ?", documentCategoryIDs)
		}
		err = documentedControlsQuery.
			Distinct("cc.control_id").
			Count(&documentedECCControlCount).Error
		if err != nil {
			return nil, err
		}

		if totalECCControlCount > 0 {
			result.ECCControlsDocumentedPercentage = roundToOneDecimal(
				float64(documentedECCControlCount) * 100 / float64(totalECCControlCount),
			)
		}
	}

	return result, nil
}

func roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func calculatePercentageChange(currentCount, previousCount int64) *float64 {
	if previousCount == 0 {
		return nil
	}

	percentage := roundToOneDecimal(
		float64(currentCount-previousCount) * 100 / float64(previousCount),
	)
	return &percentage
}

func calculateMetricPercentageChange(currentValue, previousValue float64) *float64 {
	if previousValue == 0 {
		return nil
	}

	percentage := roundToOneDecimal(
		(currentValue - previousValue) * 100 / previousValue,
	)
	return &percentage
}

func (repo *DocumentDashboardRepository) GetDocumentStatusByType(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentStatusByTypeDTO, error) {
	result := []kpiDTOs.DocumentStatusByTypeDTO{}
	documentCategoryIDs := params.Filters["document_category_id"]

	query := repo.db.Table("document_types dt").
		Select(`
			dt.id AS document_type_id,
			dt.type,
			dt.title_en,
			dt.title_ar,
			COALESCE(SUM(CASE
				WHEN ds.type IN ('approved', 'published')
					AND NOT (
						ds.type = 'published'
						AND d.next_review_date IS NOT NULL
						AND d.next_review_date < CAST(GETDATE() AS date)
					)
				THEN 1 ELSE 0 END), 0) AS approved,
			COALESCE(SUM(CASE
				WHEN ds.type = 'in-review' THEN 1 ELSE 0 END), 0) AS in_review,
			COALESCE(SUM(CASE
				WHEN ds.type IN ('draft', 'rejected', 'reverted') THEN 1 ELSE 0 END), 0) AS draft,
			COALESCE(SUM(CASE
				WHEN ds.type = 'published'
					AND d.next_review_date IS NOT NULL
					AND d.next_review_date < CAST(GETDATE() AS date)
				THEN 1 ELSE 0 END), 0) AS overdue,
			COALESCE(SUM(CASE
				WHEN ds.type IN ('approved', 'published', 'in-review', 'draft', 'rejected', 'reverted')
				THEN 1 ELSE 0 END), 0) AS total
		`)
	documentJoin := `LEFT JOIN documents d
		ON d.document_type_id = dt.id
		AND d.is_internal = 0
		AND d.is_folder = 0`
	if len(documentCategoryIDs) > 0 {
		documentJoin += " AND d.document_category_id IN ?"
		query = query.Joins(documentJoin, documentCategoryIDs)
	} else {
		query = query.Joins(documentJoin)
	}
	err := query.
		Joins("LEFT JOIN document_statuses ds ON ds.id = d.document_status_id").
		Where("dt.status = ?", "active").
		Group("dt.id, dt.type, dt.title_en, dt.title_ar").
		Order("dt.id ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *DocumentDashboardRepository) GetDocumentsByConfidentialityClassification(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentConfidentialityClassificationDTO, error) {
	result := []kpiDTOs.DocumentConfidentialityClassificationDTO{}
	documentCategoryIDs := params.Filters["document_category_id"]

	query := repo.db.Table("confidentiality_classifications cc").
		Select(`
			cc.id AS confidentiality_classification_id,
			cc.title_en,
			cc.title_ar,
			COUNT(d.id) AS count
		`)
	documentJoin := `LEFT JOIN documents d
		ON d.confidentiality_classification_id = cc.id
		AND d.is_internal = 0
		AND d.is_folder = 0`
	if len(documentCategoryIDs) > 0 {
		documentJoin += " AND d.document_category_id IN ?"
		query = query.Joins(documentJoin, documentCategoryIDs)
	} else {
		query = query.Joins(documentJoin)
	}
	err := query.
		Where("cc.status = ?", "active").
		Group("cc.id, cc.title_en, cc.title_ar").
		Order("COUNT(d.id) DESC, cc.title_en ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var totalDocumentCount int64
	for _, item := range result {
		totalDocumentCount += item.Count
	}
	if totalDocumentCount > 0 {
		for index := range result {
			result[index].Percentage = roundToOneDecimal(
				float64(result[index].Count) * 100 / float64(totalDocumentCount),
			)
		}
	}

	return result, nil
}

func (repo *DocumentDashboardRepository) GetDocumentsByResponsibleUnit(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentResponsibleUnitDTO, error) {
	result := []kpiDTOs.DocumentResponsibleUnitDTO{}
	documentCategoryIDs := params.Filters["document_category_id"]

	query := repo.db.Table("organization_units ou").
		Select(`
			ou.id AS responsible_unit_id,
			ou.title_en,
			ou.title_ar,
			COUNT(d.id) AS count
		`)
	documentJoin := `LEFT JOIN documents d
		ON d.responsible_unit_id = ou.id
		AND d.is_internal = 0
		AND d.is_folder = 0`
	if len(documentCategoryIDs) > 0 {
		documentJoin += " AND d.document_category_id IN ?"
		query = query.Joins(documentJoin, documentCategoryIDs)
	} else {
		query = query.Joins(documentJoin)
	}
	err := query.
		Group("ou.id, ou.title_en, ou.title_ar").
		Order("COUNT(d.id) DESC, ou.title_en ASC").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var totalDocumentCount int64
	for _, item := range result {
		totalDocumentCount += item.Count
	}
	if totalDocumentCount > 0 {
		for index := range result {
			result[index].Percentage = roundToOneDecimal(
				float64(result[index].Count) * 100 / float64(totalDocumentCount),
			)
		}
	}

	return result, nil
}

func (repo *DocumentDashboardRepository) GetDocumentReviewCountTrend(params *kpiDTOs.DocumentDashboardQueryParamsDTO) ([]kpiDTOs.DocumentReviewCountTrendDTO, error) {
	documentCategoryIDs := params.Filters["document_category_id"]
	now := time.Now().UTC()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startMonth := currentMonth.AddDate(0, -11, 0)
	nextMonth := currentMonth.AddDate(0, 1, 0)

	monthlyCounts := []documentMonthlyReviewCounts{}
	query := repo.db.Table("document_reviews dr").
		Select(`
			DATEFROMPARTS(YEAR(dr.review_date), MONTH(dr.review_date), 1) AS month,
			COUNT(DISTINCT dr.document_id) AS document_count
		`).
		Joins("INNER JOIN documents d ON d.id = dr.document_id").
		Where("dr.review_date >= ? AND dr.review_date < ?", startMonth, nextMonth).
		Where("dr.review_date <= CAST(GETDATE() AS date)").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false)
	if len(documentCategoryIDs) > 0 {
		query = query.Where("d.document_category_id IN ?", documentCategoryIDs)
	}
	err := query.
		Group("YEAR(dr.review_date), MONTH(dr.review_date)").
		Scan(&monthlyCounts).Error
	if err != nil {
		return nil, err
	}

	countsByMonth := make(map[string]documentMonthlyReviewCounts, len(monthlyCounts))
	for _, counts := range monthlyCounts {
		countsByMonth[counts.Month.Format("2006-01")] = counts
	}

	result := make([]kpiDTOs.DocumentReviewCountTrendDTO, 0, 12)
	for month := startMonth; month.Before(nextMonth); month = month.AddDate(0, 1, 0) {
		monthKey := month.Format("2006-01")
		item := kpiDTOs.DocumentReviewCountTrendDTO{
			Month:      monthKey,
			MonthLabel: month.Format("Jan"),
		}

		if counts, ok := countsByMonth[monthKey]; ok {
			item.Count = counts.DocumentCount
		}

		result = append(result, item)
	}

	return result, nil
}
