package services

import (
	"errors"
	"fmt"
	"time"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type AssetRiskRegisterServiceInterface interface {
	GetAll(params *assetRiskDTOs.AssetRiskRegisterListQueryParamsDTO) (*assetRiskDTOs.AssetRiskRegisterListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskRegisterDetailDTO, error)

	Create(currentUserID uint, dto *assetRiskDTOs.AssetRiskRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error)
}

type AssetRiskRegisterService struct {
	repo             repositories.AssetRiskRegisterRepositoryInterface
	organizationSvc  settingsServices.OrganizationUnitServiceInterface
	userSvc          organizationServices.UserServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
	vulnerabilitySvc VulnerabilityServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
}

func NewAssetRiskRegisterService(
	repo repositories.AssetRiskRegisterRepositoryInterface,
	organizationSvc settingsServices.OrganizationUnitServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	vulnerabilitySvc VulnerabilityServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) AssetRiskRegisterServiceInterface {
	return &AssetRiskRegisterService{
		repo:             repo,
		organizationSvc:  organizationSvc,
		userSvc:          userSvc,
		assetRegisterSvc: assetRegisterSvc,
		vulnerabilitySvc: vulnerabilitySvc,
		activityLogSvc:   activityLogSvc,
	}
}

func (arrs *AssetRiskRegisterService) GetAll(params *assetRiskDTOs.AssetRiskRegisterListQueryParamsDTO) (*assetRiskDTOs.AssetRiskRegisterListResponseDTO, error) {
	return arrs.repo.GetAll(params)
}

func (arrs AssetRiskRegisterService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	risk, err := arrs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return risk, nil
}

func (arrs AssetRiskRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskRegisterDetailDTO, error) {
	var riskDTO assetRiskDTOs.AssetRiskRegisterDetailDTO

	risk, err := arrs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(risk, &riskDTO)
	if err != nil {
		return nil, err
	}

	riskDTO.CreatedAt = risk.CreatedAt
	riskDTO.UpdatedAt = risk.UpdatedAt

	return &riskDTO, nil
}

func (arrs *AssetRiskRegisterService) Create(currentUserID uint, dto *assetRiskDTOs.AssetRiskRegisterFormDTO) (id uint, err error) {
	if err := arrs.validateAssetRiskRegisterFormDTO(dto); err != nil {
		return 0, err
	}

	dates, err := arrs.parseAssetRiskRegisterDates(dto)
	if err != nil {
		return 0, err
	}

	tx := arrs.repo.BeginTransaction()

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

	referenceNumber, err := arrs.repo.GenerateReferenceNumber(tx)
	if err != nil {
		return 0, err
	}

	riskStatus, err := arrs.repo.GetRiskStatusByType("identified", tx)
	if err != nil {
		return 0, err
	}

	impactScore := max(dto.ConfidentialityScore, dto.IntegrityScore, dto.AvailabilityScore)

	impact, err := arrs.repo.GetImpactByScore(impactScore, tx)
	if err != nil {
		return 0, err
	}

	likelihood, err := arrs.repo.GetLikelihoodByID(dto.CurrentLikelihoodID, tx)
	if err != nil {
		err = utility.NewValidationError(fmt.Errorf("likelihood does not exist"), "not_exist", "likelihood")
		return 0, err
	}

	riskValue := impact.Score * likelihood.Score

	riskRatingValue, err := arrs.repo.GetRiskRatingValueByValue(int(riskValue), tx)
	if err != nil {
		return 0, err
	}

	riskRating, err := arrs.repo.GetRiskRatingByID(riskRatingValue.RiskMatrixRiskRatingID, tx)
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		ReferenceNumber:          &referenceNumber,
		RiskStatusID:             riskStatus.ID,
		RiskCategoryID:           &dto.RiskCategoryID,
		RiskMatrixID:             riskRating.RiskMatrixID,
		Title:                    dto.Title,
		RegistrationDate:         &dates.RegistrationDate,
		RiskCode:                 &dto.RiskCode,
		DivisionID:               dto.DivisionID,
		SectionID:                &dto.SectionID,
		ResponsibleUserID:        &dto.ResponsibleUserID,
		IsStrategic:              *dto.IsStrategic,
		IsOperational:            *dto.IsOperational,
		IsInternal:               *dto.IsInternal,
		IsExternal:               *dto.IsExternal,
		ExpectedCause:            &dto.ExpectedCause,
		Justification:            &dto.Justification,
		ConfidentialityScore:     &dto.ConfidentialityScore,
		IntegrityScore:           &dto.IntegrityScore,
		AvailabilityScore:        &dto.AvailabilityScore,
		CurrentLikelihoodID:      &likelihood.ID,
		CurrentImpactID:          &impact.ID,
		CurrentRiskRatingID:      &riskRating.ID,
		CurrentRiskRatingValueID: &riskRatingValue.ID,
		CreatedBYID:              currentUserID,
		CreatedAt:                time.Now(),
		IsCyberRisk:              true,
	}

	id, err = arrs.repo.Create(riskRegister, tx)
	if err != nil {
		return 0, err
	}

	if dto.AssetID != nil {
		riskAsset := &models.RiskAsset{
			RiskRegisterID:  id,
			AssetRegisterID: *dto.AssetID,
			CreatedByID:     currentUserID,
		}

		if err = arrs.repo.CreateRiskAsset(riskAsset, tx); err != nil {
			return 0, err
		}
	}

	riskVulnerability := &models.RiskVulnerability{
		RiskRegisterID:  id,
		VulnerabilityID: dto.VulnerabilityID,
		CreatedByID:     currentUserID,
	}

	if err := arrs.repo.CreateRiskVulnerability(riskVulnerability, tx); err != nil {
		return 0, err
	}

	timelineTypes, err := arrs.repo.GetAllTimelineTypes(tx)
	if err != nil {
		return 0, err
	}

	completedStatus, err := arrs.repo.GetTimelineStatusByType("completed", tx)
	if err != nil {
		return 0, err
	}

	pendingStatus, err := arrs.repo.GetTimelineStatusByType("pending", tx)
	if err != nil {
		return 0, err
	}

	var timelineModels []models.RiskTimeline

	for _, timelineType := range timelineTypes {
		var (
			statusID      uint
			performedByID *uint
			performedAt   *time.Time
		)

		if timelineType.Type == "identified" {
			statusID = completedStatus.ID
			performedByID = &currentUserID
			now := time.Now()
			performedAt = &now
		} else {
			statusID = pendingStatus.ID
		}

		timelineModels = append(timelineModels, models.RiskTimeline{
			RiskRegisterID:     id,
			RiskTimelineTypeID: timelineType.ID,
			TimelineStatusID:   statusID,
			PerformedByID:      performedByID,
			PerformedAt:        performedAt,
		})
	}

	if err := arrs.repo.CreateRiskTimeline(timelineModels, tx); err != nil {
		return 0, err
	}

	go func() {
		existingRisk, logErr := arrs.GetFullDetailsByID(id)
		if logErr != nil {
			utility.LogError("Failed to fetch risk", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = arrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			existingRisk,
			currentUserID,
			"asset",
			"asset_risk",
			"risk",
			nil,
		)

		if logErr != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": id,
				"error":    logErr.Error(),
			})
		}
	}()

	return id, err
}

func (arrs *AssetRiskRegisterService) Update(id uint, currentUserID uint, dto *assetRiskDTOs.AssetRiskRegisterFormDTO) (err error) {
	if err = arrs.validateAssetRiskRegisterFormDTO(dto); err != nil {
		return err
	}

	dates, err := arrs.parseAssetRiskRegisterDates(dto)
	if err != nil {
		return err
	}

	tx := arrs.repo.BeginTransaction()

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

	existingRiskRegister, err := arrs.repo.GetByID(id, tx)
	if err != nil {
		return err
	}

	oldRiskRegister, err := arrs.GetFullDetailsByID(id, tx)
	if err != nil {
		return err
	}

	impactScore := max(dto.ConfidentialityScore, dto.IntegrityScore, dto.AvailabilityScore)

	impact, err := arrs.repo.GetImpactByScore(impactScore, tx)
	if err != nil {
		return err
	}

	likelihood, err := arrs.repo.GetLikelihoodByID(dto.CurrentLikelihoodID, tx)
	if err != nil {
		err = utility.NewValidationError(fmt.Errorf("likelihood does not exist"), "not_exist", "likelihood")
		return err
	}

	riskValue := impact.Score * likelihood.Score

	riskRatingValue, err := arrs.repo.GetRiskRatingValueByValue(int(riskValue), tx)
	if err != nil {
		return err
	}

	riskRating, err := arrs.repo.GetRiskRatingByID(riskRatingValue.RiskMatrixRiskRatingID, tx)
	if err != nil {
		return err
	}

	now := time.Now()

	existingRiskRegister.RiskCategoryID = &dto.RiskCategoryID
	existingRiskRegister.Title = dto.Title
	existingRiskRegister.RegistrationDate = &dates.RegistrationDate
	existingRiskRegister.RiskCode = &dto.RiskCode
	existingRiskRegister.DivisionID = dto.DivisionID
	existingRiskRegister.SectionID = &dto.SectionID
	existingRiskRegister.ResponsibleUserID = &dto.ResponsibleUserID
	existingRiskRegister.IsStrategic = *dto.IsStrategic
	existingRiskRegister.IsOperational = *dto.IsOperational
	existingRiskRegister.IsInternal = *dto.IsInternal
	existingRiskRegister.IsExternal = *dto.IsExternal
	existingRiskRegister.ExpectedCause = &dto.ExpectedCause
	existingRiskRegister.Justification = &dto.Justification

	existingRiskRegister.ConfidentialityScore = &dto.ConfidentialityScore
	existingRiskRegister.IntegrityScore = &dto.IntegrityScore
	existingRiskRegister.AvailabilityScore = &dto.AvailabilityScore

	existingRiskRegister.CurrentLikelihoodID = &likelihood.ID
	existingRiskRegister.CurrentImpactID = &impact.ID
	existingRiskRegister.CurrentRiskRatingID = &riskRating.ID
	existingRiskRegister.CurrentRiskRatingValueID = &riskRatingValue.ID

	existingRiskRegister.UpdatedBYID = &currentUserID
	existingRiskRegister.UpdatedAt = &now

	if err := arrs.repo.Update(existingRiskRegister, tx); err != nil {
		return err
	}

	existingRisk, err := arrs.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		return err
	}

	// Update Asset
	if dto.AssetID != nil {
		if len(existingRisk.Assets) == 0 {
			asset := &models.RiskAsset{
				RiskRegisterID:  id,
				AssetRegisterID: *dto.AssetID,
				CreatedByID:     currentUserID,
			}

			if err := arrs.repo.CreateRiskAsset(asset, tx); err != nil {
				return err
			}
		} else if existingRisk.Assets[0].AssetRegisterID != *dto.AssetID {
			asset := &models.RiskAsset{
				ID:              existingRisk.Assets[0].ID,
				RiskRegisterID:  id,
				AssetRegisterID: *dto.AssetID,
			}

			if err := arrs.repo.UpdateRiskAssetByRiskId(asset, tx); err != nil {
				return err
			}
		}
	} else if len(existingRisk.Assets) > 0 {
		if err := arrs.repo.DeleteRiskAssetByRiskId(id); err != nil {
			return err
		}
	}
	// Update Vulnerability
	if len(existingRisk.Vulnerabilities) == 0 {
		vulnerability := &models.RiskVulnerability{
			RiskRegisterID:  id,
			VulnerabilityID: dto.VulnerabilityID,
			CreatedByID:     currentUserID,
		}

		if err := arrs.repo.CreateRiskVulnerability(vulnerability, tx); err != nil {
			return err
		}
	} else if existingRisk.Vulnerabilities[0].VulnerabilityID != dto.VulnerabilityID {
		vulnerability := &models.RiskVulnerability{
			ID:              existingRisk.Vulnerabilities[0].ID,
			RiskRegisterID:  id,
			VulnerabilityID: dto.VulnerabilityID,
		}

		if err := arrs.repo.UpdateRiskVulnerabilityByRiskId(vulnerability, tx); err != nil {
			return err
		}
	}
	go func() {
		newRisk, logErr := arrs.GetFullDetailsByID(id)
		if logErr != nil {
			utility.LogError("Failed to fetch risk", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		oldData := map[string]interface{}{
			"risk_code":         oldRiskRegister.RiskCode,
			"title":             oldRiskRegister.Title,
			"registration_date": oldRiskRegister.RegistrationDate,
			"justification":     oldRiskRegister.Justification,
			"expected_cause":    oldRiskRegister.ExpectedCause,
			"division":          oldRiskRegister.Division.TitleEn,
			"section":           oldRiskRegister.Section.TitleEn,
			"risk_category":     oldRiskRegister.RiskCategory.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s",
				oldRiskRegister.ResponsibleUser.FirstNameEN,
				oldRiskRegister.ResponsibleUser.LastNameEN),
			"asset": func() interface{} {
				if len(oldRiskRegister.Assets) > 0 && oldRiskRegister.Assets[0] != nil {
					return oldRiskRegister.Assets[0].AssetRegister.TitleEn
				}
				return nil
			}(),
			"vulnerability": func() interface{} {
				if len(oldRiskRegister.Vulnerabilities) > 0 && oldRiskRegister.Vulnerabilities[0] != nil {
					return oldRiskRegister.Vulnerabilities[0].Vulnerability.TitleEn
				}
				return nil
			}(),
			"is_strategic":          oldRiskRegister.IsStrategic,
			"is_operational":        oldRiskRegister.IsOperational,
			"is_internal":           oldRiskRegister.IsInternal,
			"is_external":           oldRiskRegister.IsExternal,
			"confidentiality_score": oldRiskRegister.ConfidentialityScore,
			"integrity_score":       oldRiskRegister.IntegrityScore,
			"availability_score":    oldRiskRegister.AvailabilityScore,
			"current_likelihood":    oldRiskRegister.CurrentLikelihood.TitleEn,
			"current_impact":        oldRiskRegister.CurrentImpact.TitleEn,
			"current_risk_rating":   oldRiskRegister.CurrentRiskRating.TitleEn,
		}

		newData := map[string]interface{}{
			"risk_code":         newRisk.RiskCode,
			"title":             newRisk.Title,
			"registration_date": newRisk.RegistrationDate,
			"justification":     newRisk.Justification,
			"expected_cause":    newRisk.ExpectedCause,
			"division":          newRisk.Division.TitleEn,
			"section":           newRisk.Section.TitleEn,
			"risk_category":     newRisk.RiskCategory.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s",
				newRisk.ResponsibleUser.FirstNameEN,
				newRisk.ResponsibleUser.LastNameEN),
			"asset": func() interface{} {
				if len(newRisk.Assets) > 0 && newRisk.Assets[0] != nil {
					return newRisk.Assets[0].AssetRegister.TitleEn
				}
				return nil
			}(),
			"vulnerability": func() interface{} {
				if len(newRisk.Vulnerabilities) > 0 && newRisk.Vulnerabilities[0] != nil {
					return newRisk.Vulnerabilities[0].Vulnerability.TitleEn
				}
				return nil
			}(),
			"is_strategic":          newRisk.IsStrategic,
			"is_operational":        newRisk.IsOperational,
			"is_internal":           newRisk.IsInternal,
			"is_external":           newRisk.IsExternal,
			"confidentiality_score": newRisk.ConfidentialityScore,
			"integrity_score":       newRisk.IntegrityScore,
			"availability_score":    newRisk.AvailabilityScore,
			"current_likelihood":    newRisk.CurrentLikelihood.TitleEn,
			"current_impact":        newRisk.CurrentImpact.TitleEn,
			"current_risk_rating":   newRisk.CurrentRiskRating.TitleEn,
		}

		changes := utility.CompareFieldChanges(oldData, newData)

		if len(changes) > 0 {
			logErr = arrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskRegister{}.TableName(),
				id,
				changes,
				currentUserID,
				"asset",
				"asset_risk",
				"risk",
				nil,
			)

			if logErr != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"recordID": id,
					"error":    logErr.Error(),
				})
			}
		}
	}()

	return err
}

func (arrs *AssetRiskRegisterService) Delete(ids []uint, currentUserID uint) (failedIDs []uint, err error) {
	tx := arrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil || len(failedIDs) > 0 {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	var successfulIDs []uint
	var deletedRiskRegisters []map[string]interface{}

	for _, id := range ids {
		risk, err := arrs.GetFullDetailsByID(id, tx)
		if err != nil {
			return failedIDs, err
		}

		deletedRiskRegisters = append(deletedRiskRegisters, map[string]interface{}{
			"id":        id,
			"title":     risk.Title,
			"risk_code": risk.RiskCode,
			"reference": risk.ReferenceNumber,
		})

		if err = arrs.repo.DeleteRiskTimeline(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err = arrs.repo.DeleteRiskVulnerabilityByRiskId(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err = arrs.repo.DeleteRiskAssetByRiskId(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err = arrs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(failedIDs) > 0 {
		return failedIDs, err
	}

	if len(successfulIDs) > 0 {
		go func() {
			logErr := arrs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskRegister{}.TableName(),
				successfulIDs,
				deletedRiskRegisters,
				currentUserID,
				nil,
			)

			if logErr != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": logErr.Error(),
				})
			}
		}()
	}

	return failedIDs, err
}

func (arrs *AssetRiskRegisterService) GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error) {
	return arrs.repo.GetRiskStatusByType(status, tx...)
}

func (arrs *AssetRiskRegisterService) validateAssetRiskRegisterFormDTO(dto *assetRiskDTOs.AssetRiskRegisterFormDTO) error {

	if _, err := arrs.organizationSvc.GetByID(dto.DivisionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("division does not exist"), "not_exist", "division")
	}

	if _, err := arrs.organizationSvc.GetByID(dto.SectionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("section does not exist"), "not_exist", "section")
	}

	if _, err := arrs.userSvc.GetByID(dto.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if dto.AssetID != nil {
		if _, err := arrs.assetRegisterSvc.GetByID(*dto.AssetID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset does not exist"), "not_exist", "asset")
		}
	}

	if _, err := arrs.vulnerabilitySvc.GetByID(dto.VulnerabilityID); err != nil {
		return utility.NewValidationError(fmt.Errorf("vulnerability does not exist"), "not_exist", "vulnerability")
	}

	return nil
}

func (arrs *AssetRiskRegisterService) parseAssetRiskRegisterDates(dto *assetRiskDTOs.AssetRiskRegisterFormDTO) (*assetRiskDTOs.AssetRiskRegisterParsedDates, error) {
	parsed := &assetRiskDTOs.AssetRiskRegisterParsedDates{}

	registrationDate, err := utility.ParseDate(dto.RegistrationDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid registration date"), "invalid_date", "registration date")
	}

	parsed.RegistrationDate = *registrationDate

	return parsed, nil
}
