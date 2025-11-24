package services

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
	"jethur/utility"
	"sort"
	"time"
)

type KRIDashboardServiceInterface interface {
	GetKRIOverview(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIOverviewDTO, error)
	GetKRIReviewStatus(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIReviewStatusDTO, error)
	GetKRIMostFrequentBreaches() (*kpiDTOs.KriMostFrequentBreachedResponseDTO, error)
	GetKRIStatusCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIStatusCategoryDTO, error)
	GetKRISeverityBreakdown(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRISeverityBreakdownDTO, error)
	GetTopKRIRiskExposure(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriRiskExposureDTO, error)
	GetKRISeverityByOrganizationUnit(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriSeverityOrganizationUnitDTO, error)
	GetKRISeverityByCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriSeverityCategoryDTO, error)
}

type KRIDashboardService struct {
	repo repositories.KRIDashboardRepositoryInterface
}

func NewKRIDashboardService(repo repositories.KRIDashboardRepositoryInterface) KRIDashboardServiceInterface {
	return &KRIDashboardService{
		repo: repo,
	}
}

func (kds *KRIDashboardService) GetKRIOverview(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIOverviewDTO, error) {

	result, err := kds.repo.GetKRIOverview(params)
	if err != nil {
		return nil, err
	}

	// calculate high risk percent
	if result.TotalKRIs > 0 {
		result.HighRiskPercent = (result.HighRiskLinkedKRIs * 100) / result.TotalKRIs
	} else {
		result.HighRiskPercent = 0
	}

	return result, nil
}

func (kds *KRIDashboardService) GetKRIReviewStatus(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIReviewStatusDTO, error) {

	kriList, err := kds.repo.GetKRIReviewDates(params)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	todayMinus14 := today.AddDate(0, 0, -14)
	todayPlus7 := today.AddDate(0, 0, 7)

	result := &kpiDTOs.KRIReviewStatusDTO{}

	for _, item := range kriList {

		nextDate := item.NextReviewDate.Truncate(24 * time.Hour)

		// 14+ days overdue
		if nextDate.Before(todayMinus14) {
			result.FourteenDaysOverdueReviews++
			continue
		}

		// overdue but not older than 14 days
		if nextDate.Before(today) {
			result.OverdueReviews++
			continue
		}

		// upcoming (today to next 7 days)
		if nextDate.Before(todayPlus7) {
			result.UpcomingReviews++
		}
	}

	return result, nil
}

func (kds *KRIDashboardService) GetKRIStatusCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRIStatusCategoryDTO, error) {

	rawList, total, err := kds.repo.GetKRIStatusCategory(params)
	if err != nil {
		return nil, err
	}

	result := &kpiDTOs.KRIStatusCategoryDTO{
		Total:      total,
		Categories: []kpiDTOs.KRIStatusCategoryItemDTO{},
	}

	for _, item := range rawList {
		result.Categories = append(result.Categories, kpiDTOs.KRIStatusCategoryItemDTO{
			TitleEn: item.TitleEn,
			TitleAr: item.TitleAr,
			Count:   item.Count,
		})
	}

	return result, nil
}

func (kds *KRIDashboardService) GetKRIMostFrequentBreaches() (*kpiDTOs.KriMostFrequentBreachedResponseDTO, error) {

	rawList, err := kds.repo.GetKRIBreachedEvents()
	if err != nil {
		return nil, err
	}

	if len(rawList) == 0 {
		return nil, nil
	}

	breachMap := make(map[uint][]kpiDTOs.KriBreachRawDTO)

	for _, row := range rawList {

		if row.MedTh != nil && row.Value > *row.MedTh {
			breachMap[row.KriID] = append(breachMap[row.KriID], row)
		}
	}

	if len(breachMap) == 0 {
		return nil, nil
	}

	// Find most frequently breached KRI
	var topKRI uint
	var maxCount int
	for kriID, breaches := range breachMap {
		if len(breaches) > maxCount {
			maxCount = len(breaches)
			topKRI = kriID
		}
	}

	breaches := breachMap[topKRI]

	// Sort by latest
	sort.Slice(breaches, func(i, j int) bool {
		return breaches[i].ReviewedAt.After(breaches[j].ReviewedAt)
	})

	// Latest 3
	limit := 3
	if len(breaches) < 3 {
		limit = len(breaches)
	}

	finalBreaches := []kpiDTOs.KriBreachItemDTO{}
	for i := 0; i < limit; i++ {
		b := breaches[i]

		exceeded := 0
		if b.MedTh != nil {
			exceeded = b.Value - *b.MedTh
		}

		finalBreaches = append(finalBreaches, kpiDTOs.KriBreachItemDTO{
			Exceeded: exceeded,
			Date:     b.ReviewedAt.Format("02 Jan 2006"),
			Time:     b.ReviewedAt.Format("03:04 PM"),
		})
	}

	return &kpiDTOs.KriMostFrequentBreachedResponseDTO{
		KriID:    breaches[0].KriID,
		KriTitle: breaches[0].KriTitle,
		Breaches: finalBreaches,
	}, nil
}

func (kds *KRIDashboardService) GetKRISeverityBreakdown(params *kpiDTOs.KriDashboardQueryParamsDTO) (*kpiDTOs.KRISeverityBreakdownDTO, error) {

	rawList, total, err := kds.repo.GetKRISeverityBreakdown(params)
	if err != nil {
		return nil, err
	}

	result := &kpiDTOs.KRISeverityBreakdownDTO{
		Total:      total,
		Severities: []kpiDTOs.KRISeverityItemDTO{},
	}

	for _, item := range rawList {
		result.Severities = append(result.Severities, kpiDTOs.KRISeverityItemDTO{
			TitleEn:   item.TitleEn,
			TitleAr:   item.TitleAr,
			ColorCode: item.ColorCode,
			Type:      item.Type,
			Count:     item.Count,
		})
	}

	return result, nil
}

func (kds *KRIDashboardService) GetTopKRIRiskExposure(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriRiskExposureDTO, error) {

	rawList, err := kds.repo.GetKRIExposure(params)
	if err != nil {
		return nil, err
	}

	var items []kpiDTOs.KriRiskExposureDTO

	for _, r := range rawList {

		if r.Score == nil {
			continue
		}

		var threshold int
		if r.MediumThreshold != nil {
			threshold = *r.MediumThreshold
		} else if r.LowThreshold != nil {
			threshold = *r.LowThreshold
		} else {
			continue
		}

		if threshold == 0 {
			continue
		}

		exposure := float64(*r.Score) / float64(threshold)

		items = append(items, kpiDTOs.KriRiskExposureDTO{
			ID:        r.ID,
			Title:     r.Title,
			Score:     *r.Score,
			Threshold: threshold,
			Exposure:  exposure,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Exposure > items[j].Exposure
	})

	// if len(items) > 5 {
	// 	items = items[:5]
	// }

	return items, nil
}

func (kds *KRIDashboardService) GetKRISeverityByOrganizationUnit(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriSeverityOrganizationUnitDTO, error) {

	rawList, err := kds.repo.GetKRISeverityByOrganizationUnit(params)
	if err != nil {
		return nil, err
	}

	orgMap := make(map[uint]*kpiDTOs.KriSeverityOrganizationUnitDTO)

	for _, r := range rawList {

		if orgMap[r.OrganizationUnitID] == nil {
			orgMap[r.OrganizationUnitID] = &kpiDTOs.KriSeverityOrganizationUnitDTO{
				OrganizationUnitTitleEn: r.OrganizationUnitTitleEn,
				OrganizationUnitTitleAr: r.OrganizationUnitTitleAr,
				Severities:              []kpiDTOs.KRISeverityItemDTO{},
			}
		}

		orgMap[r.OrganizationUnitID].Severities = append(
			orgMap[r.OrganizationUnitID].Severities,
			kpiDTOs.KRISeverityItemDTO{
				Type:      utility.GetStringValue(r.RiskLevelType),
				TitleEn:   utility.GetStringValue(r.RiskLevelTitleEn),
				TitleAr:   utility.GetStringValue(r.RiskLevelTitleAr),
				ColorCode: utility.GetStringValue(r.RiskLevelColorCode),
				Count:     r.Count,
			},
		)
	}

	finalList := make([]kpiDTOs.KriSeverityOrganizationUnitDTO, 0, len(orgMap))
	for _, d := range orgMap {
		finalList = append(finalList, *d)
	}

	return finalList, nil
}

func (kds *KRIDashboardService) GetKRISeverityByCategory(params *kpiDTOs.KriDashboardQueryParamsDTO) ([]kpiDTOs.KriSeverityCategoryDTO, error) {

	rawList, err := kds.repo.GetKRISeverityByCategory(params)
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[uint]*kpiDTOs.KriSeverityCategoryDTO)

	for _, r := range rawList {

		if categoryMap[r.CategoryID] == nil {
			categoryMap[r.CategoryID] = &kpiDTOs.KriSeverityCategoryDTO{
				CategoryTitleEn: r.CategoryTitleEn,
				CategoryTitleAr: r.CategoryTitleAr,
				Severities:      []kpiDTOs.KRISeverityItemDTO{},
			}
		}

		categoryMap[r.CategoryID].Severities = append(
			categoryMap[r.CategoryID].Severities,
			kpiDTOs.KRISeverityItemDTO{
				Type:      utility.GetStringValue(r.RiskLevelType),
				TitleEn:   utility.GetStringValue(r.RiskLevelTitleEn),
				TitleAr:   utility.GetStringValue(r.RiskLevelTitleAr),
				ColorCode: utility.GetStringValue(r.RiskLevelColorCode),
				Count:     r.Count,
			},
		)
	}

	finalList := make([]kpiDTOs.KriSeverityCategoryDTO, 0, len(categoryMap))
	for _, d := range categoryMap {
		finalList = append(finalList, *d)
	}

	return finalList, nil
}
