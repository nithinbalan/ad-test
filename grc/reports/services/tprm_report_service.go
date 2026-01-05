package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type TPRMReportServiceInterface interface {
	GetThirdPartyByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByTypeDTO, error)
	GetThirdPartyByCountry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCountryDTO, error)
	GetThirdPartyByDepartment(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByDepartmentDTO, error)
	GetThirdPartyByCriticality(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCriticalityDTO, error)
	GetThirdPartyByLatestAssessmentStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error)
	GetThirdPartyByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByRiskRatingDTO, error)
	GetThirdPartyServiceExpiry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyServiceExpiryDTO, error)

	GetAssessmentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error)
	GetAssessmentsByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByRiskRatingDTO, error)
	GetAssessmentsByThirdPartyCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByCategoryDTO, error)
	GetAssessmentsByTemplate(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByTemplateDTO, error)
	GetOverdueAssessments(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPOverdueAssessmentsDTO, error)
}

type TPRMReportService struct {
	repo repositories.TPRMReportRepositoryInterface
}

func NewTPRMReportService(repo repositories.TPRMReportRepositoryInterface) TPRMReportServiceInterface {
	return &TPRMReportService{
		repo: repo,
	}
}

func (trs *TPRMReportService) GetThirdPartyByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByTypeDTO, error) {
	return trs.repo.GetThirdPartyByType(params)
}

func (trs *TPRMReportService) GetThirdPartyByCountry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCountryDTO, error) {
	return trs.repo.GetThirdPartyByCountry(params)
}

func (trs *TPRMReportService) GetThirdPartyByDepartment(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByDepartmentDTO, error) {
	return trs.repo.GetThirdPartyByDepartment(params)
}

func (trs *TPRMReportService) GetThirdPartyByCriticality(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByCriticalityDTO, error) {
	return trs.repo.GetThirdPartyByCriticality(params)
}

func (trs *TPRMReportService) GetThirdPartyByLatestAssessmentStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error) {
	return trs.repo.GetThirdPartyByLatestAssessmentStatus(params)
}

func (trs *TPRMReportService) GetThirdPartyByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyByRiskRatingDTO, error) {
	return trs.repo.GetThirdPartyByRiskRating(params)
}

func (trs *TPRMReportService) GetThirdPartyServiceExpiry(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ThirdPartyServiceExpiryDTO, error) {

	repoData, err := trs.repo.GetThirdPartyServiceExpiry(params)
	if err != nil {
		return nil, err
	}

	expiryTitles := map[string]struct {
		En string
		Ar string
	}{
		"expired":     {En: "Expired", Ar: "منتهي"},
		"expiring_30": {En: "Expiring in 30 Days", Ar: "تنتهي خلال 30 يوماً"},
		"expiring_60": {En: "Expiring in 60 Days", Ar: "تنتهي خلال 60 يوماً"},
	}

	for _, item := range repoData {
		if t, ok := expiryTitles[item.Category]; ok {
			item.TitleEn = t.En
			item.TitleAr = t.Ar
		}
	}

	return repoData, nil
}

func (trs *TPRMReportService) GetAssessmentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByStatusDTO, error) {
	return trs.repo.GetAssessmentsByStatus(params)
}

func (trs *TPRMReportService) GetAssessmentsByRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByRiskRatingDTO, error) {
	return trs.repo.GetAssessmentsByRiskRating(params)
}

func (trs *TPRMReportService) GetAssessmentsByThirdPartyCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByCategoryDTO, error) {
	return trs.repo.GetAssessmentsByThirdPartyCategory(params)
}

func (trs *TPRMReportService) GetAssessmentsByTemplate(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPAssessmentsByTemplateDTO, error) {
	return trs.repo.GetAssessmentsByTemplate(params)
}

func (trs *TPRMReportService) GetOverdueAssessments(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.TPOverdueAssessmentsDTO, error) {

    repoData, err := trs.repo.GetOverdueAssessments(params)
    if err != nil {
        return nil, err
    }

    titles := map[string]struct {
        En string
        Ar string
    }{
        "overdue":     {En: "Overdue", Ar: "متأخر"},
        "upcoming_7":  {En: "Due in 7 Days", Ar: "مستحق خلال 7 أيام"},
        "upcoming_14": {En: "Due in 14 Days", Ar: "مستحق خلال 14 يومًا"},
    }

    for _, item := range repoData {
        if t, ok := titles[item.Category]; ok {
            item.TitleEn = t.En
            item.TitleAr = t.Ar
        }
    }

    return repoData, nil
}
