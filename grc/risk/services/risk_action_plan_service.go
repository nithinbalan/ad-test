package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskActionPlanServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskActionPlan, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error)
	Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, risk *riskDTOs.RiskActionPlanFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error
}

// Service and implementations
type RiskActionPlanService struct {
	repo                riskRepositories.RiskActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	riskRegisterSvc     RiskRegisterServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
}

func NewRiskActionPlanService(
	repo riskRepositories.RiskActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	) RiskActionPlanServiceInterface {
	return &RiskActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		riskRegisterSvc:     riskRegisterSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		notificationSvc: notificationSvc,
	}
}

func (raps *RiskActionPlanService) GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error) {
	if riskID != 0 {
		if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}	
	}
	return raps.repo.GetAll(riskID, params)
}

func (raps RiskActionPlanService) GetByID(riskID uint, id uint) (*models.RiskActionPlan, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlan, nil
}

func (raps RiskActionPlanService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error) {
	var riskActionPlanDTO riskDTOs.RiskActionPlanDetailDTO

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskActionPlan, &riskActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskActionPlanDTO, nil
}

func (raps *RiskActionPlanService) Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error) {
	tx := raps.repo.BeginTransaction()

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	riskActionPlan := &models.RiskActionPlan{
		RiskRegisterID:     riskID,
		Title:              RiskActionPlanFormDTO.Title,
		Description:        RiskActionPlanFormDTO.Description,
		SeverityLevelID:    RiskActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  RiskActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         RiskActionPlanFormDTO.Percentage,
		Comments:           RiskActionPlanFormDTO.Comments,
		ActionPlanStatusID: RiskActionPlanFormDTO.ActionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := raps.repo.Create(riskActionPlan)
	if err != nil {
		return 0, err
	}
	actionPlanStatusID := RiskActionPlanFormDTO.ActionPlanStatusID
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   id,
		ActionPlanStatusID: actionPlanStatusID,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}
	_, err = raps.repo.CreateInitialRiskActionPlanStatusUpdate(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := raps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlan{}.TableName(),
			id,
			RiskActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

	//sends notification and email 
	raps.sendRiskActionPlanNotification(riskActionPlan,currentUserID,"create")

	}()

	return id, nil
}

func (raps *RiskActionPlanService) Update(riskID uint, id uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	tx := raps.repo.BeginTransaction()

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskActionPlan, err := raps.GetByID(riskID, id)

	if err != nil {
		return err
	}

	//Store the old responsible user ID for notification
	oldResponsibleUserID := existingRiskActionPlan.ResponsibleUserID

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Risk ActionPlan
	existingRiskActionPlan.Title = RiskActionPlanFormDTO.Title
	existingRiskActionPlan.Description = RiskActionPlanFormDTO.Description
	existingRiskActionPlan.SeverityLevelID = RiskActionPlanFormDTO.SeverityLevelID
	existingRiskActionPlan.ResponsibleUserID = RiskActionPlanFormDTO.ResponsibleUserID
	existingRiskActionPlan.TargetDate = dates.TargetDate
	existingRiskActionPlan.Percentage = RiskActionPlanFormDTO.Percentage
	existingRiskActionPlan.Comments = RiskActionPlanFormDTO.Comments
	existingRiskActionPlan.ActionPlanStatusID = RiskActionPlanFormDTO.ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func(){
		//sends notification if responsible user has changed
		if oldResponsibleUserID !=existingRiskActionPlan.ResponsibleUserID {
		raps.sendRiskActionPlanNotification(existingRiskActionPlan,currentUserID,"update")
		}
	}()

	return nil
}

func (raps *RiskActionPlanService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		riskActionPlan, err := raps.repo.GetFullDetailsByID(id)

		if err != nil {
			continue
		}
		if len(riskActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range riskActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := raps.repo.DeleteRiskActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = raps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			go raps.sendRiskActionPlanNotification(riskActionPlan,currentUserID,"delete")
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := raps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			
		}()
	}

	return failedIds, nil
}

func (raps *RiskActionPlanService) StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {

	existingRiskActionPlan, err := raps.GetByID(riskID, id)

	if err != nil {
		return err
	}

	existingRiskActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (raps *RiskActionPlanService) validateRiskActionPlanFormDTO(RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	if RiskActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := raps.severityLevelSvc.GetByID(*RiskActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := raps.actionPlanStatusSvc.GetByID(RiskActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (raps *RiskActionPlanService) parseRiskActionPlanDates(dto *riskDTOs.RiskActionPlanFormDTO) (*riskDTOs.RiskActionPlanParsedDates, error) {
	parsed := &riskDTOs.RiskActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}


//Methods to send notifications for risk action plans
func (raps *RiskActionPlanService) sendRiskActionPlanNotification(riskActionPlan *models.RiskActionPlan, currentUserID uint,action string) {

// Fetch responsible user details to get the email
responsibleUser, err:=raps.riskRegisterSvc.GetByID(riskActionPlan.RiskRegisterID)
if err != nil {
	utility.LogError("Failed to fetch responsible user details", map[string]interface{}{
		"error": err.Error(),
	})
	return
}

//prepare email and notification data
emailTo:=[]string{responsibleUser.ResponsibleUser.Email}
subject:=fmt.Sprintf("Risk Action Plan/%s: Responsibility %s",riskActionPlan.Title,action)
appURL:=utility.GetAppURL()
emailData:= map[string]interface{}{
"RiskTitle":              riskActionPlan.Title,
		"ActionPlanDescription":  riskActionPlan.Description,
		"TargetDate":             riskActionPlan.TargetDate.Format("2006-01-02"),
		"AssignedOn":             time.Now().Format("2006-01-02"),
		"RedirectURL":            fmt.Sprintf("%s/risk-action-plan/%d", appURL, riskActionPlan.ID),
}

//email template
emailNotificationDTO:= &commonDTOs.EmailNotificationDTO{
	EmailTo: emailTo,
	Subject: subject,
	TemplateRelativePath: "risk/templates/email/risk_action_plan_notification.html",
	EmailData: emailData,
}

// Notifications DTO
dbNotificationDTO := &commonDTOs.DBNotificationDTO{
	UserID:               riskActionPlan.ResponsibleUserID,
	ModuleIdentifier:     "risk",
	NotificationType:     string(commonDTOs.AssignmentType),
	HeadingEn:            "Risk Action Plan Assignment",
	TitleEn:              riskActionPlan.Title,
	IsSystemNotification: false,
	CreatedByID:          &currentUserID,
}

// Send notification and email
err = raps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
if err != nil {
	utility.LogError("Failed to send notifications", map[string]interface{}{
		"error": err.Error(),
	})
}


}