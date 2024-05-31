package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingActionPlanServiceInterface interface {
	GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error)
	GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error)
	GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error)
	Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error)
	Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error
	Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(meetingID uint, id uint, currentUserID uint, actionPlanStatusID uint, percentage uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type MeetingActionPlanService struct {
	repo                meetingRepositories.MeetingActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	meetingSvc          MeetingServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewMeetingActionPlanService(
	repo meetingRepositories.MeetingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingSvc MeetingServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface) MeetingActionPlanServiceInterface {
	return &MeetingActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		meetingSvc:          meetingSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (maps *MeetingActionPlanService) GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return maps.repo.GetAll(meetingID, params)
}

func (maps MeetingActionPlanService) GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingActionPlan, nil
}

func (maps MeetingActionPlanService) GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error) {
	var meetingActionPlanDTO meetingDTOs.MeetingActionPlanDetailDTO

	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meetingActionPlan, &meetingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &meetingActionPlanDTO, nil
}

func (maps *MeetingActionPlanService) Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error) {
	meeting, err := maps.meetingSvc.GetByID(meetingID)
	if err != nil {
		return 0, err
	}

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	meetingActionPlan := &models.MeetingActionPlan{
		MeetingID:          meetingID,
		Title:              MeetingActionPlanFormDTO.Title,
		Description:        MeetingActionPlanFormDTO.Description,
		SeverityLevelID:    MeetingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  MeetingActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         *MeetingActionPlanFormDTO.Percentage,
		ActionPlanStatusID: MeetingActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := maps.repo.Create(meetingActionPlan)
	if err != nil {
		return 0, err
	}

	eventDescriptionEn := "You have been added as responsible for a new Action Plan"
	eventDescriptionAr := "تمت إضافتك كمسؤول لخطة عمل جديدة"

	messageEn := fmt.Sprintf("%s: %s", eventDescriptionEn, MeetingActionPlanFormDTO.Title)
	messageAr := fmt.Sprintf("%s: %s", eventDescriptionAr, MeetingActionPlanFormDTO.Title)

	redirectURL := fmt.Sprintf("meeting/%d/action-plans/%d", meetingID, id)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserID:           MeetingActionPlanFormDTO.ResponsibleUserID,
		ModuleIdentifier: "meeting_action_plan",
		Icon:             models.CreateIcon,
		MessageEn:        messageEn,
		MessageAr:        messageAr,
		RedirectURL:      redirectURL,
		CreatedBy:        currentUserID,
	}

	resposibleUser, err := maps.userSvc.GetByID(MeetingActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, err
	}

	emailTo := []string{resposibleUser.Email}
	subject := fmt.Sprintf("meeting/%s: Action Plan Added", meeting.Title)
	templateRelativePath := "meeting/templates/email/create_action_plan_mail.html"
	emailData := meetingDTOs.SendCreateActionPlanMailDTO{
		FirstNameEn: resposibleUser.FirstNameEn,
		LastNameEn:  resposibleUser.LastNameEn,
		MessageEn:   messageEn,
		RedirectURL: redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              emailTo,
		Subject:              subject,
		TemplateRelativePath: templateRelativePath,
		EmailData:            emailData,
	}

	err = maps.notificationSvc.SendNotification(
		dbNotificationDTO, emailNotificationDTO, "create_meeting_action_plan")
	if err != nil {
		return 0, err
	}

	go func() {
		err := maps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlan{}.TableName(),
			id,
			MeetingActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (maps *MeetingActionPlanService) Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)

	if err != nil {
		return err
	}

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingMeetingActionPlan.Title = MeetingActionPlanFormDTO.Title
	existingMeetingActionPlan.Description = MeetingActionPlanFormDTO.Description
	existingMeetingActionPlan.SeverityLevelID = MeetingActionPlanFormDTO.SeverityLevelID
	existingMeetingActionPlan.ResponsibleUserID = MeetingActionPlanFormDTO.ResponsibleUserID
	existingMeetingActionPlan.TargetDate = dates.TargetDate
	existingMeetingActionPlan.Percentage = *MeetingActionPlanFormDTO.Percentage
	existingMeetingActionPlan.ActionPlanStatusID = MeetingActionPlanFormDTO.ActionPlanStatusID
	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan)

	if err != nil {
		return err
	}

	return nil
}

func (maps *MeetingActionPlanService) Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := maps.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = maps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := maps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingActionPlan{}.TableName(),
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

func (maps *MeetingActionPlanService) StatusUpdate(meetingID uint, id uint, currentUserID uint, actionPlanStatusID uint, percentage uint, tx ...utility.DatabaseTransaction) error {
	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)

	if err != nil {
		return err
	}

	existingMeetingActionPlan.ActionPlanStatusID = actionPlanStatusID
	existingMeetingActionPlan.Percentage = percentage
	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan, tx...)

	if err != nil {
		return err
	}

	return nil
}

func (maps *MeetingActionPlanService) validateMeetingActionPlanFormDTO(meetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	if meetingActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := maps.severityLevelSvc.GetByID(*meetingActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := maps.userSvc.GetByID(uint(meetingActionPlanFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := maps.actionPlanStatusSvc.GetByID(meetingActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (maps *MeetingActionPlanService) parseMeetingActionPlanDates(dto *meetingDTOs.MeetingActionPlanFormDTO) (*meetingDTOs.MeetingActionPlanParsedDates, error) {
	parsed := &meetingDTOs.MeetingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
