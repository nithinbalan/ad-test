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

	"gorm.io/gorm"
)

type RiskInherentControlAnalysisServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error)
	Create(riskID uint, currentUserID uint) error
	Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetControlSummary(id uint) (map[string]interface{}, error)
}

// Service and implementations
type RiskInherentControlAnalysisService struct {
	repo            riskRepositories.RiskInherentControlAnalysisRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
	controlSvc      organizationServices.ControlServiceInterface
	controlTestSvc  organizationServices.ControlTestServiceInterface
}

func NewRiskInherentControlAnalysisService(
	repo riskRepositories.RiskInherentControlAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
) RiskInherentControlAnalysisServiceInterface {
	return &RiskInherentControlAnalysisService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
		controlSvc:      controlSvc,
		controlTestSvc:  controlTestSvc,
	}
}

func (rcas *RiskInherentControlAnalysisService) GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error) {
	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return rcas.repo.GetAll(riskID, params)
}

func (rcas *RiskInherentControlAnalysisService) GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error) {
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

func (rcas *RiskInherentControlAnalysisService) Create(riskID uint, currentUserID uint) error {
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
			Type:                          "inherent",
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

	return nil
}

func (rcas *RiskInherentControlAnalysisService) Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error {
	tx := rcas.repo.BeginTransaction()

	if _, err := rcas.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskControlAnalysis, err := rcas.GetByID(riskID, id)
	if err != nil {
		return err
	}

	existingRiskControlAnalysis.ControlEfficiencyStatusID = controlTestFormDTO.ControlEfficiencyStatusID
	existingRiskControlAnalysis.ControlImplementationStatusID = controlTestFormDTO.ControlImplementationStatusID

	err = rcas.repo.Update(existingRiskControlAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = rcas.controlTestSvc.Create(existingRiskControlAnalysis.RiskControl.ControlID, currentUserID, &riskID, controlTestFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rcas *RiskInherentControlAnalysisService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
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
func (rcas *RiskInherentControlAnalysisService) GetControlSummary(riskID uint) (map[string]interface{}, error) {
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
