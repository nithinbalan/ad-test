package repositories

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"math"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type DocumentDashboardRepositoryInterface interface {
	GetDocumentSummary() (*kpiDTOs.DocumentDashboardSummaryDTO, error)
	GetDocumentStatusByType() ([]kpiDTOs.DocumentStatusByTypeDTO, error)
	GetDocumentOnTimeReviewTrend() ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error)
}

type DocumentDashboardRepository struct {
	db *gorm.DB
}

type documentReviewTimelinessCounts struct {
	CurrentPeriodReviews       int64
	CurrentPeriodOnTimeReviews int64
}

type documentApprovalTimeResult struct {
	AverageApprovalTimeDays float64
}

type complianceRegisterLegislation struct {
	ID            uint
	LegislationNo string
}

type documentMonthlyReviewCounts struct {
	Month         time.Time
	TotalReviews  int64
	OnTimeReviews int64
}

var eccLegislationNumberPattern = regexp.MustCompile(`(?i)(^|[^a-z0-9])ECC([^a-z0-9]|$)`)

const documentOnTimeReviewTarget = 85.0

func NewDocumentDashboardRepository(db *gorm.DB) DocumentDashboardRepositoryInterface {
	return &DocumentDashboardRepository{db: db}
}

func (repo *DocumentDashboardRepository) GetDocumentSummary() (*kpiDTOs.DocumentDashboardSummaryDTO, error) {
	result := &kpiDTOs.DocumentDashboardSummaryDTO{}

	err := repo.db.Table("documents d").
		Select(`
			COUNT(d.id) AS total_documents,
			COALESCE(SUM(CASE WHEN ds.type = 'approved' THEN 1 ELSE 0 END), 0) AS approved_documents,
			COALESCE(SUM(CASE WHEN ds.type = 'published' THEN 1 ELSE 0 END), 0) AS published_documents
		`).
		Joins("LEFT JOIN document_statuses ds ON ds.id = d.document_status_id").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false).
		Scan(result).Error
	if err != nil {
		return nil, err
	}

	reviewCounts := &documentReviewTimelinessCounts{}
	err = repo.db.Table("document_reviews dr").
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
				THEN 1 ELSE 0 END), 0) AS current_period_on_time_reviews
		`).
		Joins("INNER JOIN documents d ON d.id = dr.document_id").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false).
		Scan(reviewCounts).Error
	if err != nil {
		return nil, err
	}

	if reviewCounts.CurrentPeriodReviews > 0 {
		result.ReviewedOnTimePercentage = roundToOneDecimal(
			float64(reviewCounts.CurrentPeriodOnTimeReviews) * 100 / float64(reviewCounts.CurrentPeriodReviews),
		)
	}

	approvalTime := &documentApprovalTimeResult{}
	err = repo.db.Table("document_versions dv").
		Select(`
			COALESCE(
				AVG(CAST(DATEDIFF(second, dv.submitted_at, final_approval.final_approved_at) AS FLOAT) / 86400.0),
				0
			) AS average_approval_time_days
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
		Where("d.is_internal = ? AND d.is_folder = ?", false, false).
		Scan(approvalTime).Error
	if err != nil {
		return nil, err
	}

	result.AverageApprovalTimeDays = roundToOneDecimal(approvalTime.AverageApprovalTimeDays)

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

	if len(eccComplianceRegisterIDs) > 0 && result.TotalDocuments > 0 {
		var mappedDocumentCount int64
		err = repo.db.Table("document_compliance_requirements dcr").
			Joins("INNER JOIN documents d ON d.id = dcr.document_id").
			Where("dcr.compliance_register_id IN ?", eccComplianceRegisterIDs).
			Where("dcr.compliance_clause_id IS NOT NULL").
			Where("d.is_internal = ? AND d.is_folder = ?", false, false).
			Distinct("dcr.document_id").
			Count(&mappedDocumentCount).Error
		if err != nil {
			return nil, err
		}

		result.ECCControlsDocumentedPercentage = roundToOneDecimal(
			float64(mappedDocumentCount) * 100 / float64(result.TotalDocuments),
		)
	}

	return result, nil
}

func roundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}

func (repo *DocumentDashboardRepository) GetDocumentStatusByType() ([]kpiDTOs.DocumentStatusByTypeDTO, error) {
	result := []kpiDTOs.DocumentStatusByTypeDTO{}

	err := repo.db.Table("document_types dt").
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
		`).
		Joins(`LEFT JOIN documents d
			ON d.document_type_id = dt.id
			AND d.is_internal = 0
			AND d.is_folder = 0`).
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

func (repo *DocumentDashboardRepository) GetDocumentOnTimeReviewTrend() ([]kpiDTOs.DocumentOnTimeReviewTrendDTO, error) {
	now := time.Now().UTC()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	startMonth := currentMonth.AddDate(0, -11, 0)
	nextMonth := currentMonth.AddDate(0, 1, 0)

	monthlyCounts := []documentMonthlyReviewCounts{}
	err := repo.db.Table("document_reviews dr").
		Select(`
			DATEFROMPARTS(YEAR(dr.review_date), MONTH(dr.review_date), 1) AS month,
			COUNT(dr.id) AS total_reviews,
			COALESCE(SUM(CASE
				WHEN dr.review_date <= d.next_review_date THEN 1 ELSE 0 END
			), 0) AS on_time_reviews
		`).
		Joins("INNER JOIN documents d ON d.id = dr.document_id").
		Where("dr.review_date >= ? AND dr.review_date < ?", startMonth, nextMonth).
		Where("dr.review_date <= CAST(GETDATE() AS date)").
		Where("d.next_review_date IS NOT NULL").
		Where("d.is_internal = ? AND d.is_folder = ?", false, false).
		Group("YEAR(dr.review_date), MONTH(dr.review_date)").
		Scan(&monthlyCounts).Error
	if err != nil {
		return nil, err
	}

	countsByMonth := make(map[string]documentMonthlyReviewCounts, len(monthlyCounts))
	for _, counts := range monthlyCounts {
		countsByMonth[counts.Month.Format("2006-01")] = counts
	}

	result := make([]kpiDTOs.DocumentOnTimeReviewTrendDTO, 0, 12)
	for month := startMonth; month.Before(nextMonth); month = month.AddDate(0, 1, 0) {
		monthKey := month.Format("2006-01")
		item := kpiDTOs.DocumentOnTimeReviewTrendDTO{
			Month:            monthKey,
			MonthLabel:       month.Format("Jan"),
			TargetPercentage: documentOnTimeReviewTarget,
		}

		if counts, ok := countsByMonth[monthKey]; ok && counts.TotalReviews > 0 {
			percentage := roundToOneDecimal(
				float64(counts.OnTimeReviews) * 100 / float64(counts.TotalReviews),
			)
			variance := roundToOneDecimal(percentage - documentOnTimeReviewTarget)
			item.OnTimePercentage = &percentage
			item.VarianceFromTarget = &variance
		}

		result = append(result, item)
	}

	return result, nil
}
