package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
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
	ExportRiskRegister(data []*assetRiskDTOs.AssetRiskRegisterListDTO, currentUserID uint) (*os.File, error)

	UpdateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error

	GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error)
	GetImpactByScore(score int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error)
	GetImpactByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error)
	GetLikelihoodByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixLikelihood, error)
	GetRiskRatingValueByValue(value int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRatingValue, error)
	GetRiskRatingByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRating, error)

	GetTimelineTypeByType(timelineType string, tx ...utility.DatabaseTransaction) (*models.RiskTimelineType, error)
	GetTimelineStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.TimelineStatus, error)
	GetRiskTimelineByRiskAndType(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error)
	UpdateRiskTimeline(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error
}

type AssetRiskRegisterService struct {
	repo             repositories.AssetRiskRegisterRepositoryInterface
	organizationSvc  settingsServices.OrganizationUnitServiceInterface
	userSvc          organizationServices.UserServiceInterface
	userSettingSvc   organizationServices.UserSettingServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
	vulnerabilitySvc VulnerabilityServiceInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	notificationSvc  commonServices.NotificationServiceInterface
}

func NewAssetRiskRegisterService(
	repo repositories.AssetRiskRegisterRepositoryInterface,
	organizationSvc settingsServices.OrganizationUnitServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	vulnerabilitySvc VulnerabilityServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AssetRiskRegisterServiceInterface {
	return &AssetRiskRegisterService{
		repo:             repo,
		organizationSvc:  organizationSvc,
		userSvc:          userSvc,
		userSettingSvc:   userSettingSvc,
		assetRegisterSvc: assetRegisterSvc,
		vulnerabilitySvc: vulnerabilitySvc,
		activityLogSvc:   activityLogSvc,
		notificationSvc:  notificationSvc,
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
	riskDTO.RegistrationDateStr = risk.RegistrationDate.Format("2006-01-02")

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
			"registration_date": oldRiskRegister.RegistrationDateStr,
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
			"registration_date": newRisk.RegistrationDateStr,
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

		oldResponsibleUserID := uint(oldRiskRegister.ResponsibleUser.ID)
		newResponsibleUserID := dto.ResponsibleUserID

		if oldResponsibleUserID == newResponsibleUserID {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "asset-risk-register"

		oldResponsibleUserAction := "responsible-user-remove"
		oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{oldResponsibleUserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "ISMS Risk Register Responsibility Removed",
			HeadingAr:            "تمت إزالة مسؤولية تسجيل مخاطر ISMS",
			ReferenceNumber:      newRisk.ReferenceNumber,
			TitleEn:              newRisk.Title,
			TitleAr:              newRisk.Title,
			Entity:               &entity,
			Action:               &oldResponsibleUserAction,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		newResponsibleUserAction := "responsible-user-assign"
		newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{newResponsibleUserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "ISMS Risk Register Responsibility Assigned",
			HeadingAr:            "مسؤولية تسجيل مخاطر ISMS المعينة",
			ReferenceNumber:      newRisk.ReferenceNumber,
			TitleEn:              newRisk.Title,
			TitleAr:              newRisk.Title,
			Entity:               &entity,
			Action:               &newResponsibleUserAction,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		oldResponsibleUser, notifyErr := arrs.userSvc.GetByID(oldResponsibleUserID)
		if notifyErr != nil {
			utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
				"error": notifyErr.Error(),
			})
			return
		}

		newResponsibleUser, notifyErr := arrs.userSvc.GetByID(newResponsibleUserID)
		if notifyErr != nil {
			utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
				"error": notifyErr.Error(),
			})
			return
		}

		updatedByUser, notifyErr := arrs.userSvc.GetByID(currentUserID)
		if notifyErr != nil {
			utility.LogError("Failed to fetch updating user", map[string]interface{}{
				"error": notifyErr.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		redirectURL := fmt.Sprintf("%s/asset/risks/%d/overview", appURL, id)

		oldResponsibleUserEmailData := map[string]interface{}{
			"FirstNameEn":          oldResponsibleUser.FirstNameEn,
			"LastNameEn":           oldResponsibleUser.LastNameEn,
			"RiskReferenceNumber":  newRisk.ReferenceNumber,
			"RiskTitle":            newRisk.Title,
			"RiskJustification":    utility.GetStringValue(newRisk.Justification),
			"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
			"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
			"RedirectURL":          redirectURL,
		}

		oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{oldResponsibleUser.Email},
			Subject:              fmt.Sprintf("ISMS Risk Register: %s - Responsibility Removed", newRisk.Title),
			TemplateRelativePath: "asset/templates/email/asset_risk_register_remove_responsible_user_mail.html",
			EmailData:            oldResponsibleUserEmailData,
		}

		newResponsibleUserEmailData := map[string]interface{}{
			"FirstNameEn":           newResponsibleUser.FirstNameEn,
			"LastNameEn":            newResponsibleUser.LastNameEn,
			"RiskReferenceNumber":   newRisk.ReferenceNumber,
			"RiskTitle":             newRisk.Title,
			"RiskJustification":     utility.GetStringValue(newRisk.Justification),
			"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
			"AssignedByLastNameEn":  updatedByUser.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{newResponsibleUser.Email},
			Subject:              fmt.Sprintf("ISMS Risk Register: %s - Responsibility Assigned", newRisk.Title),
			TemplateRelativePath: "asset/templates/email/asset_risk_register_change_responsible_user_mail.html",
			EmailData:            newResponsibleUserEmailData,
		}

		if notifyErr := arrs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "asset_risk"); notifyErr != nil {
			utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
				"error": notifyErr.Error(),
			})
		}

		if notifyErr := arrs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "asset_risk"); notifyErr != nil {
			utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
				"error": notifyErr.Error(),
			})
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
	var deletionNotifications []struct {
		ResponsibleUserID uint
		Title             string
		ReferenceNumber   string
	}

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

		deletionNotifications = append(deletionNotifications, struct {
			ResponsibleUserID uint
			Title             string
			ReferenceNumber   string
		}{
			ResponsibleUserID: uint(risk.ResponsibleUser.ID),
			Title:             risk.Title,
			ReferenceNumber:   risk.ReferenceNumber,
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

			entity := "asset-risk-register"
			action := "delete"

			for _, notification := range deletionNotifications {
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{notification.ResponsibleUserID},
					ModuleIdentifier:     "asset_risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "ISMS Risk Register Deleted",
					HeadingAr:            "تم حذف سجل مخاطر ISMS",
					ReferenceNumber:      notification.ReferenceNumber,
					TitleEn:              notification.Title,
					TitleAr:              notification.Title,
					Entity:               &entity,
					Action:               &action,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, notifyErr := arrs.userSvc.GetByID(notification.ResponsibleUserID)
				if notifyErr != nil {
					utility.LogError("Failed to fetch responsible user for deletion notification", map[string]interface{}{
						"error": notifyErr.Error(),
					})
					continue
				}

				emailData := map[string]interface{}{
					"FirstNameEn":         responsibleUser.FirstNameEn,
					"LastNameEn":          responsibleUser.LastNameEn,
					"RiskReferenceNumber": notification.ReferenceNumber,
					"RiskTitle":           notification.Title,
					"DeletedOn":           time.Now().Format("2006-01-02"),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              []string{responsibleUser.Email},
					Subject:              fmt.Sprintf("ISMS Risk Register: %s - Deleted", notification.Title),
					TemplateRelativePath: "asset/templates/email/asset_risk_register_delete_mail.html",
					EmailData:            emailData,
				}

				if notifyErr := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); notifyErr != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": notifyErr.Error(),
					})
				}
			}
		}()
	}

	return failedIDs, err
}

func (arrs *AssetRiskRegisterService) ExportRiskRegister(data []*assetRiskDTOs.AssetRiskRegisterListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := arrs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Risk Title",
			"Risk Code",
			"Risk Category",
			"Division",
			"Department",
			"Risk Level",
			"Source",
			"Confidentiality Score",
			"Integrity Score",
			"Availability Score",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Registration Date",
			"Created By",
			"Created At",
			"Updated By",
			"Updated At",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Risk Title",
			"Risk Code",
			"Risk Category",
			"Division",
			"Department",
			"Risk Level",
			"Source",
			"Confidentiality Score",
			"Integrity Score",
			"Availability Score",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Registration Date",
			"Created By",
			"Created At",
			"Updated By",
			"Updated At",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, risk := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				risk.ReferenceNumber,
				risk.Title,
				risk.RiskCode,
				risk.RiskCategoryTitleAr,
				risk.DivisionTitleAr,
				risk.SectionTitleAr,
				assetRiskLevel(risk.IsStrategic, risk.IsOperational, userSetting.Language),
				assetRiskSource(risk.IsInternal, risk.IsExternal, userSetting.Language),
				risk.ConfidentialityScore,
				risk.IntegrityScore,
				risk.AvailabilityScore,
				risk.CurrentLikelihoodTitleAr,
				risk.CurrentImpactTitleAr,
				risk.CurrentRiskRatingTitleAr,
				risk.RegistrationDate,
				fmt.Sprintf("%s %s", risk.CreatedUserFirstNameAr, risk.CreatedUserLastNameAr),
				risk.CreatedAt,
				func() interface{} {
					if risk.UpdatedUserFirstNameAr != nil && risk.UpdatedUserLastNameAr != nil {
						return fmt.Sprintf("%s %s", *risk.UpdatedUserFirstNameAr, *risk.UpdatedUserLastNameAr)
					}
					return nil
				}(),
				utility.GetStringValue(risk.UpdatedAt),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				risk.ReferenceNumber,
				risk.Title,
				risk.RiskCode,
				risk.RiskCategoryTitleEn,
				risk.DivisionTitleEn,
				risk.SectionTitleEn,
				assetRiskLevel(risk.IsStrategic, risk.IsOperational, userSetting.Language),
				assetRiskSource(risk.IsInternal, risk.IsExternal, userSetting.Language),
				risk.ConfidentialityScore,
				risk.IntegrityScore,
				risk.AvailabilityScore,
				risk.CurrentLikelihoodTitleEn,
				risk.CurrentImpactTitleEn,
				risk.CurrentRiskRatingTitleEn,
				risk.RegistrationDate,
				fmt.Sprintf("%s %s", risk.CreatedUserFirstNameEn, risk.CreatedUserLastNameEn),
				risk.CreatedAt,
				func() interface{} {
					if risk.UpdatedUserFirstNameEn != nil && risk.UpdatedUserLastNameEn != nil {
						return fmt.Sprintf("%s %s", *risk.UpdatedUserFirstNameEn, *risk.UpdatedUserLastNameEn)
					}
					return nil
				}(),
				utility.GetStringValue(risk.UpdatedAt),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func assetRiskSource(isInternal bool, isExternal bool, language string) string {
	if language == "ar" {
		if isInternal {
			return "داخلي"
		}
		if isExternal {
			return "خارجي"
		}
		return ""
	}

	if isInternal {
		return "Internal"
	}
	if isExternal {
		return "External"
	}
	return ""
}

func assetRiskLevel(isStrategic bool, isOperational bool, language string) string {
	if language == "ar" {
		if isStrategic {
			return "استراتيجي"
		}
		if isOperational {
			return "تشغيلي"
		}
		return ""
	}

	if isStrategic {
		return "Strategic"
	}
	if isOperational {
		return "Operational"
	}
	return ""
}

func (arrs *AssetRiskRegisterService) GetRiskStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.RiskStatus, error) {
	return arrs.repo.GetRiskStatusByType(status, tx...)
}

func (arrs *AssetRiskRegisterService) UpdateRiskRegister(riskRegister *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	return arrs.repo.Update(riskRegister, tx...)
}

func (arrs *AssetRiskRegisterService) GetImpactByScore(score int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error) {
	return arrs.repo.GetImpactByScore(score, tx...)
}

func (arrs *AssetRiskRegisterService) GetImpactByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixImpact, error) {
	return arrs.repo.GetImpactByID(id, tx...)
}

func (arrs *AssetRiskRegisterService) GetLikelihoodByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixLikelihood, error) {
	return arrs.repo.GetLikelihoodByID(id, tx...)
}

func (arrs *AssetRiskRegisterService) GetRiskRatingValueByValue(value int, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRatingValue, error) {
	return arrs.repo.GetRiskRatingValueByValue(value, tx...)
}

func (arrs *AssetRiskRegisterService) GetRiskRatingByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskMatrixRiskRating, error) {
	return arrs.repo.GetRiskRatingByID(id, tx...)
}

func (arrs *AssetRiskRegisterService) GetTimelineTypeByType(timelineType string, tx ...utility.DatabaseTransaction) (*models.RiskTimelineType, error) {
	return arrs.repo.GetTimelineTypeByType(timelineType, tx...)
}

func (arrs *AssetRiskRegisterService) GetTimelineStatusByType(status string, tx ...utility.DatabaseTransaction) (*models.TimelineStatus, error) {
	return arrs.repo.GetTimelineStatusByType(status, tx...)
}

func (arrs *AssetRiskRegisterService) GetRiskTimelineByRiskAndType(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error) {
	return arrs.repo.GetRiskTimelineByRiskAndType(riskID, timelineTypeID, tx...)
}

func (arrs *AssetRiskRegisterService) UpdateRiskTimeline(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error {
	return arrs.repo.UpdateRiskTimeline(riskTimeline, tx...)
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
