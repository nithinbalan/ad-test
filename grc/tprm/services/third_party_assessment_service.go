package services

import (
	"errors"
	"fmt"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"

	"github.com/google/uuid"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThirdPartyAssessmentServiceInterface interface {
	GetAll(params *dtos.ThirdPartyAssessmentListQueryParamsDTO) (*dtos.ThirdPartyAssessmentListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyAssessment, error)
	GetFullDetailsByID(id uint) (*dtos.ThirdPartyAssessmentDetailDTO, error)
	Create(currentUserID uint, dto *dtos.ThirdPartyAssessmentFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.ThirdPartyAssessmentFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	ShareLink(id uint, currentUserID uint, dto *dtos.ThirdPartyAssessmentShareLinkFormDTO) error
}

type ThirdPartyAssessmentService struct {
	repo                          tprmRepo.ThirdPartyAssessmentRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	severityLevelSvc              masterServices.SeverityLevelServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	orgUnitSvc                    settingsServices.OrganizationUnitServiceInterface
	thirdPartyRegisterSvc         ThirdPartyRegisterServiceInterface
	questionSetSvc                QuestionSetServiceInterface
	thirdPartyAssessmentStatusSvc ThirdPartyAssessmentStatusServiceInterface
	userSettingSvc                organizationServices.UserSettingServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
}

func NewThirdPartyAssessmentService(
	repo tprmRepo.ThirdPartyAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	orgUnitSvc settingsServices.OrganizationUnitServiceInterface,
	thirdPartyRegisterSvc ThirdPartyRegisterServiceInterface,
	questionSetSvc QuestionSetServiceInterface,
	thirdPartyAssessmentStatusSvc ThirdPartyAssessmentStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) *ThirdPartyAssessmentService {
	return &ThirdPartyAssessmentService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		severityLevelSvc:              severityLevelSvc,
		userSvc:                       userSvc,
		orgUnitSvc:                    orgUnitSvc,
		thirdPartyRegisterSvc:         thirdPartyRegisterSvc,
		questionSetSvc:                questionSetSvc,
		thirdPartyAssessmentStatusSvc: thirdPartyAssessmentStatusSvc,
		userSettingSvc:                userSettingSvc,
		notificationSvc:               notificationSvc,
	}
}

func (tpas *ThirdPartyAssessmentService) GetAll(params *dtos.ThirdPartyAssessmentListQueryParamsDTO) (*dtos.ThirdPartyAssessmentListResponseDTO, error) {
	return tpas.repo.GetAll(params)
}

func (tpas *ThirdPartyAssessmentService) GetByID(id uint) (*models.ThirdPartyAssessment, error) {
	tpr, err := tpas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return tpr, nil
}

func (tpas *ThirdPartyAssessmentService) GetFullDetailsByID(id uint) (*dtos.ThirdPartyAssessmentDetailDTO, error) {
	var dto dtos.ThirdPartyAssessmentDetailDTO

	tpr, err := tpas.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(tpr, &dto)
	if err != nil {
		return nil, err
	}

	if tpr.DueDate != nil {
		formatted := tpr.DueDate.Format("2006-01-02")
		dto.DueDateString = &formatted
	}

	dto.CreatedAt = tpr.CreatedAt

	return &dto, nil
}

func (tpas *ThirdPartyAssessmentService) Create(currentUserID uint, form *dtos.ThirdPartyAssessmentFormDTO) (uint, error) {
	tx := tpas.repo.BeginTransaction()

	if err := tpas.validateForm(form); err != nil {
		return 0, err
	}

	questionSet, err := tpas.questionSetSvc.GetFullDetailsByID(form.QuestionSetID)
	if err != nil {
		return 0, err
	}

	dates, err := tpas.parseDates(form)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := tpas.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	assessmentStatus, err := tpas.thirdPartyAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	uid := uuid.New().String()

	tpr := &models.ThirdPartyAssessment{
		ReferenceNumber:              referenceNumber,
		Title:                        form.Title,
		ThirdPartyRegisterID:         form.ThirdPartyRegisterID,
		QuestionSetID:                form.QuestionSetID,
		Message:                      form.Message,
		DueDate:                      dates.DueDate,
		ThirdPartyAssessmentStatusID: assessmentStatus.ID,
		UID:                          &uid,
		CreatedBYID:                  currentUserID,
	}

	id, err := tpas.repo.Create(tpr, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, category := range questionSet.Categories {
		for _, question := range category.Questions {
			questionDTO := &models.TPRMAssessmentQuestion{
				ThirdPartyAssessmentID:      id,
				ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
				QuestionTitleEn:             question.QuestionTitleEn,
				QuestionTitleAr:             question.QuestionTitleAr,
				CreatedBYID:                 currentUserID,
			}
			questionId, err := tpas.repo.CreateAssessmentQuestion(questionDTO, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			for _, option := range question.AnswerOptions {
				optionDTO := &models.TPRMAssessmentQuestionAnswerOption{
					TPRMAssessmentQuestionID:    questionId,
					ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
					AnswerOptionTitleEn:         option.AnswerOptionTitleEn,
					AnswerOptionTitleAr:         option.AnswerOptionTitleAr,
					RiskLevelID:                 option.RiskLevelID,
					RangeFrom:                   option.RangeFrom,
					RangeTo:                     option.RangeTo,
					CreatedBYID:                 &currentUserID,
				}
				_, err := tpas.repo.CreateAssessmentQuestionAnswerOption(optionDTO, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	tx.Commit()

	return id, nil
}

func (tpas *ThirdPartyAssessmentService) Update(id uint, currentUserID uint, form *dtos.ThirdPartyAssessmentFormDTO) error {
	tx := tpas.repo.BeginTransaction()
	existing, err := tpas.GetByID(id)
	if err != nil {
		return err
	}

	if err := tpas.validateForm(form); err != nil {
		return err
	}
	existingAssessment, err := tpas.repo.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	questionSet, err := tpas.questionSetSvc.GetFullDetailsByID(form.QuestionSetID)
	if err != nil {
		return err
	}

	dates, err := tpas.parseDates(form)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.QuestionSetID != form.QuestionSetID {
		for _, question := range existingAssessment.Questions {
			if err := tpas.repo.DeleteAnswerOptionsByQuestionID(question.ID, tx); err != nil {
				tx.Rollback()
			}
		}
		if err := tpas.repo.DeleteQuestionsByAssessmentID(id, tx); err != nil {
			tx.Rollback()
		}

		for _, category := range questionSet.Categories {
			for _, question := range category.Questions {
				questionDTO := &models.TPRMAssessmentQuestion{
					ThirdPartyAssessmentID:      id,
					ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
					QuestionTitleEn:             question.QuestionTitleEn,
					QuestionTitleAr:             question.QuestionTitleAr,
					CreatedBYID:                 currentUserID,
				}
				questionId, err := tpas.repo.CreateAssessmentQuestion(questionDTO, tx)
				if err != nil {
					tx.Rollback()
					return err
				}

				for _, option := range question.AnswerOptions {
					optionDTO := &models.TPRMAssessmentQuestionAnswerOption{
						TPRMAssessmentQuestionID:    questionId,
						ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
						AnswerOptionTitleEn:         option.AnswerOptionTitleEn,
						AnswerOptionTitleAr:         option.AnswerOptionTitleAr,
						RiskLevelID:                 option.RiskLevelID,
						RangeFrom:                   option.RangeFrom,
						RangeTo:                     option.RangeTo,
						CreatedBYID:                 &currentUserID,
					}
					_, err := tpas.repo.CreateAssessmentQuestionAnswerOption(optionDTO, tx)
					if err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	existing.Title = form.Title
	existing.ThirdPartyRegisterID = form.ThirdPartyRegisterID
	existing.QuestionSetID = form.QuestionSetID
	existing.Message = form.Message
	existing.DueDate = dates.DueDate
	existing.UpdatedBYID = &currentUserID

	err = tpas.repo.Update(existing, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (tpas *ThirdPartyAssessmentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	tx := tpas.repo.BeginTransaction()
	for _, id := range ids {

		existingAssessment, err := tpas.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		for _, question := range existingAssessment.Questions {
			if err := tpas.repo.DeleteAnswerOptionsByQuestionID(question.ID, tx); err != nil {
				tx.Rollback()
				failedIDs = append(failedIDs, id)
				continue
			}
		}
		if err := tpas.repo.DeleteQuestionsByAssessmentID(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		err = tpas.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		tx.Commit()

	}

	return failedIDs, nil
}

func (tpas *ThirdPartyAssessmentService) validateForm(form *dtos.ThirdPartyAssessmentFormDTO) error {
	if _, err := tpas.thirdPartyRegisterSvc.GetByID(form.ThirdPartyRegisterID); err != nil {
		return utility.NewValidationError(errors.New("third party register not found"), "not_exist", "third_party_register")
	}

	if _, err := tpas.questionSetSvc.GetByID(form.QuestionSetID); err != nil {
		return utility.NewValidationError(errors.New("question set not found"), "not_exist", "question set")
	}

	return nil
}

func (tpas *ThirdPartyAssessmentService) parseDates(dto *dtos.ThirdPartyAssessmentFormDTO) (*dtos.ThirdPartyAssessmentParsedDates, error) {
	parsed := &dtos.ThirdPartyAssessmentParsedDates{}
	var err error

	if dto.DueDate != nil {
		parsed.DueDate, err = utility.ParseDate(*dto.DueDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid due date"), "invalid_date", "due_date")
		}
	}
	return parsed, nil
}

func (tpas *ThirdPartyAssessmentService) ShareLink(id uint, currentUserID uint, form *dtos.ThirdPartyAssessmentShareLinkFormDTO) error {

	thirdPartyAssessment, err := tpas.repo.GetFullDetailsByID(id)
	if err != nil {
		utility.LogError("Failed to fetch assessment", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	thirdPartyRegister, err := tpas.thirdPartyRegisterSvc.GetFullDetailsByID(thirdPartyAssessment.ThirdPartyRegisterID)
	if err != nil {
		utility.LogError("Failed to fetch third party register", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	go func() {

		thirdPartyRegister, err := tpas.thirdPartyRegisterSvc.GetFullDetailsByID(thirdPartyRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch kri register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if thirdPartyAssessment.UID != nil {
			// Prepare email

			emailTo := []string{form.Email}
			subject := "Assessment"
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"VendorName":  thirdPartyRegister.Title,
				"ExpiryDate":  "2025-08-30",
				"RedirectURL": fmt.Sprintf("%s/vendor-assessment/%s", appURL, *thirdPartyAssessment.UID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "tprm/templates/email/third_party_assessment_mail.html",
				EmailData:            emailData,
			}

			err = tpas.notificationSvc.SendNotification(nil, emailNotificationDTO, "kri_register")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return err
}
