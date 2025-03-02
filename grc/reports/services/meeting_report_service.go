package services

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
	settingService "jethur/grc/settings/services"
)

type MeetingReportServiceInterface interface {
	GetMeetingByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByTypeDTO, error)
	GetMeetingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByStatusDTO, error)
	GetMeetingByOrganizer(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByOrganizerDTO, error)
	GetMeetingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByResponsibleUnitDTO, error)
	GetMeetingMinutesByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingMinutesByStatusDTO, error)

	GetMeetingActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByStatusDTO, error)
	GetMeetingActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUserDTO, error)
	GetMeetingActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUnitDTO, error)
	GetActionPlanByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ActionPlanByOpenVsClosedDTO, error)

	GetMeetingMOMCount(params *reportDTOs.ReportListQueryParamsDTO) (*reportDTOs.MeetingMOMCountDTO, error)
	GetMeetingsWithMOM(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithMOMDTO, error)
	GetMeetingsWithOpenActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOpenActionPlanDTO, error)
	GetMeetingsWithOverdueActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOverdueActionPlanDTO, error)
	GetMeetingCountByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByYearDTO, error)
	GetMeetingCountByUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByUserDTO, error)
}

type MeetingReportService struct {
	repo                 repositories.MeetingReportRepositoryInterface
	organizationLevelSvc settingService.OrganizationLevelServiceInterface
}

func NewMeetingReportService(repo repositories.MeetingReportRepositoryInterface,
	organizationLevelSvc settingService.OrganizationLevelServiceInterface) MeetingReportServiceInterface {
	return &MeetingReportService{
		repo:                 repo,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (mrs *MeetingReportService) GetMeetingByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByTypeDTO, error) {
	return mrs.repo.GetMeetingByType(params)
}

func (mrs *MeetingReportService) GetMeetingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByStatusDTO, error) {
	return mrs.repo.GetMeetingByStatus(params)
}

func (mrs *MeetingReportService) GetMeetingByOrganizer(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByOrganizerDTO, error) {
	return mrs.repo.GetMeetingByOrganizer(params)
}

func (mrs *MeetingReportService) GetMeetingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByResponsibleUnitDTO, error) {
	primaryLevel, err := mrs.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}
	var primaryOrgUnitField string

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}
	return mrs.repo.GetMeetingByResponsibleUnit(params, primaryOrgUnitField)
}

func (mrs *MeetingReportService) GetMeetingMinutesByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingMinutesByStatusDTO, error) {
	return mrs.repo.GetMeetingMinutesByStatus(params)
}

func (mrs *MeetingReportService) GetMeetingActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByStatusDTO, error) {
	return mrs.repo.GetMeetingActionPlanByStatus(params)
}

func (mrs *MeetingReportService) GetMeetingActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUserDTO, error) {
	return mrs.repo.GetMeetingActionPlanByResponsibleUser(params)
}

func (mrs *MeetingReportService) GetMeetingActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUnitDTO, error) {
	primaryLevel, err := mrs.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}
	var primaryOrgUnitField string

	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}
	return mrs.repo.GetMeetingActionPlanByResponsibleUnit(params, primaryOrgUnitField)
}

func (mrs *MeetingReportService) GetActionPlanByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ActionPlanByOpenVsClosedDTO, error) {
	return mrs.repo.GetActionPlanByOpenVsClosed(params)
}

func (mrs *MeetingReportService) GetMeetingMOMCount(params *reportDTOs.ReportListQueryParamsDTO) (*reportDTOs.MeetingMOMCountDTO, error) {
	return mrs.repo.GetMeetingMOMCount(params)
}

func (mrs *MeetingReportService) GetMeetingsWithMOM(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithMOMDTO, error) {
	return mrs.repo.GetMeetingsWithMOM(params)
}

func (mrs *MeetingReportService) GetMeetingsWithOpenActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOpenActionPlanDTO, error) {
	return mrs.repo.GetMeetingsWithOpenActionPlans(params)
}

func (mrs *MeetingReportService) GetMeetingsWithOverdueActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOverdueActionPlanDTO, error) {
	return mrs.repo.GetMeetingsWithOverdueActionPlans(params)
}

func (mrs *MeetingReportService) GetMeetingCountByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByYearDTO, error) {
	return mrs.repo.GetMeetingCountByYear(params)
}

func (mrs *MeetingReportService) GetMeetingCountByUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByUserDTO, error) {
	return mrs.repo.GetMeetingCountByUser(params)
}
