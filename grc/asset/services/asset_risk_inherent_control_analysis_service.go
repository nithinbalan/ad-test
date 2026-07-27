package services

import (
	"errors"
	"fmt"
	"time"

	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRiskInherentControlAnalysisServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskControlAnalysisDetailDTO, error)
	Create(riskID uint, currentUserID uint) error
	Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetControlSummary(id uint) (map[string]interface{}, error)
}

// Service and implementations
type AssetRiskInherentControlAnalysisService struct {
	repo                 repositories.AssetRiskInherentControlAnalysisRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface
	controlSvc           organizationServices.ControlServiceInterface
	controlTestSvc       organizationServices.ControlTestServiceInterface
}

func NewAssetRiskInherentControlAnalysisService(
	repo repositories.AssetRiskInherentControlAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlTestSvc organizationServices.ControlTestServiceInterface,
) AssetRiskInherentControlAnalysisServiceInterface {
	return &AssetRiskInherentControlAnalysisService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		assetRiskRegisterSvc: assetRiskRegisterSvc,
		controlSvc:           controlSvc,
		controlTestSvc:       controlTestSvc,
	}
}

func (svc *AssetRiskInherentControlAnalysisService) GetAll(riskID uint, params *riskDTOs.RiskControlAnalysisListQueryParamsDTO) ([]*riskDTOs.RiskControlAnalysisListDTO, error) {
	if _, err := svc.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	controlAnalyses, err := svc.repo.GetAll(riskID, params)
	if err != nil {
		return nil, err
	}

	// Populate associated documents for each control
	for _, analysis := range controlAnalyses {
		if analysis.ControlID > 0 {
			control, err := svc.controlSvc.GetFullDetailsByID(analysis.ControlID)
			if err == nil && control != nil && len(control.AssociatedDocuments) > 0 {
				analysis.AssociatedDocuments = control.AssociatedDocuments
			} else {
				analysis.AssociatedDocuments = []organizationDTOs.ControlAssociatedDocumentDTO{}
			}
		} else {
			analysis.AssociatedDocuments = []organizationDTOs.ControlAssociatedDocumentDTO{}
		}
	}

	return controlAnalyses, nil
}

func (svc AssetRiskInherentControlAnalysisService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskControlAnalysisDetailDTO, error) {
	var riskControlAnalysisDTO riskDTOs.RiskControlAnalysisDetailDTO

	if _, err := svc.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskControlAnalysis, err := svc.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskControlAnalysis, &riskControlAnalysisDTO)
	if err != nil {
		return nil, err
	}

	return &riskControlAnalysisDTO, nil
}

func (svc *AssetRiskInherentControlAnalysisService) GetByID(riskID uint, id uint) (*models.RiskControlAnalysis, error) {
	if _, err := svc.assetRiskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskControlAnalysis, err := svc.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskControlAnalysis, nil
}

func (svc *AssetRiskInherentControlAnalysisService) Create(riskID uint, currentUserID uint) (err error) {
	tx := svc.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = svc.assetRiskRegisterSvc.GetByID(riskID, tx); err != nil {
		return err
	}

	riskControls, err := svc.repo.GetAllControlsByRiskID(riskID)
	if err != nil {
		return err
	}
	var controlAnalysis []*models.RiskControlAnalysis
	var ids []uint
	for _, riskControl := range riskControls {
		control, err := svc.controlSvc.GetByID(riskControl.ControlID)
		if err != nil {
			return err
		}

		riskControlAnalysis := &models.RiskControlAnalysis{
			RiskControlID:                 riskControl.ID,
			Type:                          "inherent",
			ControlEfficiencyStatusID:     control.ControlEfficiencyStatusID,
			ControlImplementationStatusID: *control.ControlImplementationStatusID,
			Score:                         control.EfficiencyScore,
			CreatedByID:                   currentUserID,
			CreatedAt:                     time.Now(),
		}

		_, err = svc.repo.Create(riskControlAnalysis)
		if err != nil {
			return err
		}
		controlAnalysis = append(controlAnalysis, riskControlAnalysis)
		ids = append(ids, riskControlAnalysis.ID)
	}

	timelineType, err := svc.assetRiskRegisterSvc.GetTimelineTypeByType("control_analysis", tx)
	if err != nil {
		return err
	}
	timelineStatus, err := svc.assetRiskRegisterSvc.GetTimelineStatusByType("completed", tx)
	if err != nil {
		return err
	}
	riskTimeline, err := svc.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return err
	}
	currentTime := time.Now()
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = &currentUserID
	riskTimeline.PerformedAt = &currentTime
	if err = svc.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return err
	}

	go func() {
		logErr := svc.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.RiskControlAnalysis{}.TableName(),
			ids,
			controlAnalysis,
			currentUserID,
			nil,
		)

		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskControlAnalysis{}.TableName(),
				"recordID":  ids,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (svc *AssetRiskInherentControlAnalysisService) Update(riskID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) (err error) {
	tx := svc.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if _, err = svc.assetRiskRegisterSvc.GetByID(riskID, tx); err != nil {
		return err
	}

	existingRiskControlAnalysis, err := svc.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldControlAnalysis, err := svc.GetFullDetailsByID(riskID, id)
	if err != nil {
		return err
	}

	testId, err := svc.controlTestSvc.Create(existingRiskControlAnalysis.RiskControl.ControlID, currentUserID, &riskID, controlTestFormDTO, "asset_risk", tx)
	if err != nil {
		return err
	}

	// Get the created control test to retrieve calculated effectiveness values (within transaction)
	controlTest, err := svc.controlTestSvc.GetByID(testId, tx)
	if err != nil {
		return err
	}

	existingRiskControlAnalysis.ControlEfficiencyStatusID = controlTestFormDTO.ControlEfficiencyStatusID
	existingRiskControlAnalysis.ControlImplementationStatusID = controlTestFormDTO.ControlImplementationStatusID
	existingRiskControlAnalysis.ControlTestID = testId
	existingRiskControlAnalysis.Score = controlTestFormDTO.Score
	existingRiskControlAnalysis.ControlEffectivenessScore = controlTest.ControlEffectivenessScore
	existingRiskControlAnalysis.ControlEffectivenessRatingID = controlTest.ControlEffectivenessRatingID
	existingRiskControlAnalysis.RiskReduction = controlTest.RiskReduction
	existingRiskControlAnalysis.CreatedByID = currentUserID
	existingRiskControlAnalysis.CreatedAt = time.Now()

	if err = svc.repo.Update(existingRiskControlAnalysis, tx); err != nil {
		return err
	}

	go func() {
		existingControlAnalysis, logErr := svc.GetFullDetailsByID(riskID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch risk control analysis", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"control_efficiency_status":     oldControlAnalysis.ControlEfficiencyStatus.TitleEn,
			"control_implementation_status": oldControlAnalysis.ControlImplementationStatus.TitleEn,
			"control_test":                  oldControlAnalysis.ControlTestID,
			"score":                         oldControlAnalysis.Score,
		}

		newData := map[string]interface{}{
			"control_efficiency_status":     existingControlAnalysis.ControlEfficiencyStatus.TitleEn,
			"control_implementation_status": existingControlAnalysis.ControlImplementationStatus.TitleEn,
			"control_test":                  existingControlAnalysis.ControlTestID,
			"score":                         existingControlAnalysis.Score,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			logErr := svc.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskControlAnalysis{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset", "asset_risk", "control_analysis", nil,
			)
			if logErr != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.RiskControlAnalysis{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     logErr.Error(),
				})
			}
		}
	}()

	return nil
}

func (svc *AssetRiskInherentControlAnalysisService) Delete(riskID uint, ids []uint, currentUserID uint) (failedIds []uint, err error) {
	tx := svc.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIds) > 0 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	var successfulIds []uint
	var deletedControlAnalysis []map[string]interface{}

	for _, id := range ids {
		riskControlAnalysis, getErr := svc.GetByID(riskID, id)
		if getErr != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if deleteErr := svc.repo.Delete(riskControlAnalysis.ID, tx); deleteErr != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedControlAnalysis = append(deletedControlAnalysis, map[string]interface{}{
			"id":               id,
			"risk_register_id": riskID,
		})

		successfulIds = append(successfulIds, id)
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	timelineType, err := svc.assetRiskRegisterSvc.GetTimelineTypeByType("control_analysis", tx)
	if err != nil {
		return nil, err
	}
	timelineStatus, err := svc.assetRiskRegisterSvc.GetTimelineStatusByType("pending", tx)
	if err != nil {
		return nil, err
	}
	riskTimeline, err := svc.assetRiskRegisterSvc.GetRiskTimelineByRiskAndType(riskID, timelineType.ID, tx)
	if err != nil {
		return nil, err
	}
	riskTimeline.TimelineStatusID = timelineStatus.ID
	riskTimeline.PerformedByID = nil
	riskTimeline.PerformedAt = nil
	if err = svc.assetRiskRegisterSvc.UpdateRiskTimeline(riskTimeline, tx); err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			logErr := svc.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskControlAnalysis{}.TableName(),
				successfulIds,
				deletedControlAnalysis,
				currentUserID,
				nil,
			)
			if logErr != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", logErr)
			}
		}()
	}

	return failedIds, nil
}

func (svc *AssetRiskInherentControlAnalysisService) GetControlSummary(riskID uint) (map[string]interface{}, error) {
	controlTests, err := svc.repo.GetControlSummary(riskID)
	if err != nil {
		return nil, err
	}

	if len(controlTests) == 0 {
		return map[string]interface{}{
			"total_control_score":          0,
			"implementation_status":        []map[string]interface{}{},
			"efficiency_percentage":        0,
			"implementation_percentage":    0,
			"total_efficiency":             []map[string]interface{}{},
			"total_control_count":          0,
			"control_effectiveness_rating": []map[string]interface{}{},
		}, nil
	}

	implementationStatusCount := make(map[string]int)
	efficiencyStatusCount := make(map[string]int)
	effectivenessRatingCount := make(map[string]int)
	implementationDetails := make(map[string]models.ControlImplementationStatus)
	efficiencyDetails := make(map[string]models.ControlEfficiencyStatus)
	effectivenessRatingDetails := make(map[uint]models.ControlEffectivenessRating)

	for _, ct := range controlTests {
		if ct.ControlImplementationStatus != nil {
			implStatus := ct.ControlImplementationStatus.Type
			implementationStatusCount[implStatus]++
			implementationDetails[implStatus] = *ct.ControlImplementationStatus
		}

		if ct.ControlEfficiencyStatus != nil {
			effStatus := ct.ControlEfficiencyStatus.Type
			efficiencyStatusCount[effStatus]++
			efficiencyDetails[effStatus] = *ct.ControlEfficiencyStatus
		}

		if ct.ControlEffectivenessRating != nil {
			ratingType := ct.ControlEffectivenessRating.Type
			effectivenessRatingCount[ratingType]++
			effectivenessRatingDetails[ct.ControlEffectivenessRating.ID] = *ct.ControlEffectivenessRating
		}
	}

	totalControls := len(controlTests)

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

	effectivenessRatings := make([]map[string]interface{}, 0)
	for _, rating := range effectivenessRatingDetails {
		count := effectivenessRatingCount[rating.Type]
		effectivenessRatings = append(effectivenessRatings, map[string]interface{}{
			"id":         rating.ID,
			"type":       rating.Type,
			"title_en":   rating.TitleEn,
			"title_ar":   rating.TitleAr,
			"count":      count,
			"min_score":  rating.MinScore,
			"color_code": rating.ColorCode,
			"max_score":  rating.MaxScore,
		})
	}

	summary := map[string]interface{}{
		"total_control_score":          (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"implementation_status":        implementationStatus,
		"efficiency_percentage":        (float64(efficiencyStatusCount["efficient"]) / float64(totalControls)) * 100,
		"implementation_percentage":    (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"total_efficiency":             efficiencyStatus,
		"total_control_count":          totalControls,
		"control_effectiveness_rating": effectivenessRatings,
	}

	return summary, nil
}
