package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskResidualControlAnalysisServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskControlAnalysisDetailDTO, error)
	Create(riskID uint, currentUserID uint) error
	Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetControlSummary(id uint) (map[string]interface{}, error)
}

// Service and implementations
type RiskResidualControlAnalysisService struct {
	repo            riskRepositories.RiskResidualControlAnalysisRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
	controlSvc      organizationServices.ControlServiceInterface
	controlTestSvc  organizationServices.ControlTestServiceInterface
	riskTimelineSvc RiskTimelineServiceInterface
}

func NewRiskResidualControlAnalysisService(
	repo riskRepositories.RiskResidualControlAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
) RiskResidualControlAnalysisServiceInterface {
	return &RiskResidualControlAnalysisService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
		controlSvc:      controlSvc,
		controlTestSvc:  controlTestSvc,
		riskTimelineSvc: riskTimelineSvc,
	}
}

func (rcas *RiskResidualControlAnalysisService) GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error) {
	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rcas.repo.GetAll(riskID, params)
}

func (rcas *RiskResidualControlAnalysisService) GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error) {
	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskControlAnalysis, err := rcas.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskControlAnalysis, nil
}
func (rcas RiskResidualControlAnalysisService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskControlAnalysisDetailDTO, error) {
	var riskResidualControlAnalysisDTO riskDTOs.RiskControlAnalysisDetailDTO

	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskResidualControlAnalysis, err := rcas.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskResidualControlAnalysis, &riskResidualControlAnalysisDTO)
	if err != nil {
		return nil, err
	}

	return &riskResidualControlAnalysisDTO, nil
}

func (rcas *RiskResidualControlAnalysisService) Create(riskID uint, currentUserID uint) error {
	tx := rcas.repo.BeginTransaction()

	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	riskControls, err := rcas.repo.GetAllControlsByRiskID(riskID)
	if err != nil {
		return err
	}

	for _, riskControl := range riskControls {
		control, err := rcas.controlSvc.GetByID(riskControl.ControlID)
		if err != nil {
			return err
		}

		riskControlAnalysis := &models.RiskControlAnalysis{
			RiskControlID:                 riskControl.ID,
			Type:                          "residual",
			ControlEfficiencyStatusID:     *control.ControlEfficiencyStatusID,
			ControlImplementationStatusID: *control.ControlImplementationStatusID,
			CreatedByID:                   currentUserID,
			CreatedAt:                     time.Now(),
		}

		_, err = rcas.repo.Create(riskControlAnalysis)
		if err != nil {
			return err
		}
	}
	//Update timeline
	//todo
	timelineType, err := rcas.riskTimelineSvc.GetByType("residual_control_analysis")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := rcas.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := rcas.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (rcas *RiskResidualControlAnalysisService) Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error {
	tx := rcas.repo.BeginTransaction()

	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskControlAnalysis, err := rcas.GetByID(riskID, id)
	if err != nil {
		return err
	}

	testId, err := rcas.controlTestSvc.Create(existingRiskControlAnalysis.RiskControl.ControlID, currentUserID, &riskID, controlTestFormDTO, "risk", tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingRiskControlAnalysis.ControlEfficiencyStatusID = controlTestFormDTO.ControlEfficiencyStatusID
	existingRiskControlAnalysis.ControlImplementationStatusID = controlTestFormDTO.ControlImplementationStatusID
	existingRiskControlAnalysis.ControlTestID = testId
	existingRiskControlAnalysis.Score = controlTestFormDTO.Score
	existingRiskControlAnalysis.CreatedByID = currentUserID //todo (tested_by or currentUserID)
	existingRiskControlAnalysis.CreatedAt = time.Now()

	err = rcas.repo.Update(existingRiskControlAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rcas *RiskResidualControlAnalysisService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := rcas.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		riskControlAnalysis, err := rcas.GetByID(riskID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = rcas.repo.Delete(riskControlAnalysis.ID, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
	}

	//Delete-> update timeline performedAt as null
	//todo
	timelineType, err := rcas.riskTimelineSvc.GetByType("residual_control_analysis")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	timelineStatus, err := rcas.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    nil,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      nil,
	}
	if err := rcas.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(successfulIds) > 0 {
		tx.Commit()
		go func() {
			err := rcas.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskControlAnalysis{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIds, nil
}
func (rcas *RiskResidualControlAnalysisService) GetControlSummary(riskID uint) (map[string]interface{}, error) {
	riskControlAnalyses, err := rcas.repo.GetControlSummary(riskID)
	if err != nil {
		return nil, err
	}

	if len(riskControlAnalyses) == 0 {
		return map[string]interface{}{
			"total_control_score":       0,
			"implementation_status":     []map[string]interface{}{},
			"efficiency_percentage":     0,
			"implementation_percentage": 0,
			"total_efficiency":          []map[string]interface{}{},
			"total_control_count":       0,
		}, nil
	}

	implementationStatusCount := make(map[string]int)
	efficiencyStatusCount := make(map[string]int)
	implementationDetails := make(map[string]models.ControlImplementationStatus)
	efficiencyDetails := make(map[string]models.ControlEfficiencyStatus)

	for _, rca := range riskControlAnalyses {
		implStatus := rca.ControlImplementationStatus.Type
		implementationStatusCount[implStatus]++
		implementationDetails[implStatus] = *rca.ControlImplementationStatus

		effStatus := rca.ControlEfficiencyStatus.Type
		efficiencyStatusCount[effStatus]++
		efficiencyDetails[effStatus] = *rca.ControlEfficiencyStatus
	}

	totalControls := len(riskControlAnalyses)

	implementationStatus := make([]map[string]interface{}, 0)
	for status, count := range implementationStatusCount {
		implDetail := implementationDetails[status]
		implementationStatus = append(implementationStatus, map[string]interface{}{
			"id":         implDetail.ID,
			"type":       implDetail.Type,
			"title_en":   implDetail.TitleEn,
			"title_ar":   implDetail.TitleAr,
			"count":      count,
			"color_code": implDetail.ColorCode,
		})
	}

	efficiencyStatus := make([]map[string]interface{}, 0)
	for status, count := range efficiencyStatusCount {
		effDetail := efficiencyDetails[status]
		efficiencyStatus = append(efficiencyStatus, map[string]interface{}{
			"id":         effDetail.ID,
			"type":       effDetail.Type,
			"title_en":   effDetail.TitleEn,
			"title_ar":   effDetail.TitleAr,
			"count":      count,
			"color_code": effDetail.ColorCode,
		})
	}

	summary := map[string]interface{}{
		"total_control_score":       (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"implementation_status":     implementationStatus,
		"efficiency_percentage":     (float64(efficiencyStatusCount["efficient"]) / float64(totalControls)) * 100,
		"implementation_percentage": (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"total_efficiency":          efficiencyStatus,
		"total_control_count":       totalControls,
	}

	return summary, nil
}
