package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
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
	"github.com/xuri/excelize/v2"

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

	ShareLink(id uint, currentUserID uint) error
	ExportTPRMAssessment(data *dtos.ThirdPartyAssessmentDetailDTO, currentUserID uint) (*os.File, error)
	ValidateTPRMAssessmentTemplate(assessmentID uint, file io.Reader) ([]map[string]interface{}, error)
	ValidateImportTPRMAssessment(assessmentID uint, input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportTPRMAssessment(assessmentID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error)
	MarkAsSend(assessmentID uint, currentUserID uint) error

	GetAssessmentAnswerSummary(assessmentID uint) ([]dtos.AnswerSummaryItemDTO, error)
	ExportTPRMAssessmentReport(assessmentID uint) (string, error)
	ExportThirdPartyAssessments(data []*dtos.ThirdPartyAssessmentListDTO, currentUserID uint) (*os.File, error)
	UpdateAnswer(id uint, questionId, currentUserID uint, dto *dtos.ThirdPartyAssessmentAnswerFormDTO) error

	GetAssessmentSectionDetails(assessmentID uint, sectionID uint) (*dtos.TPRMAssessmentSectionDTO, error)
	GetAssessmentSectionQuestionDetails(assessmentId uint, sectionId uint, questionId uint) (*dtos.TPRMAssessmentQuestionDetailDTO, error)

	GetAllAnswerDocuments(assessmentID uint, questionID uint, answerID uint) ([]dtos.AnswerDocumentListDTO, error)
	CreateAnswerDocument(assessmentID uint, questionID uint, answerID uint, userID uint, dto *dtos.AnswerDocumentFormDTO) error
	DeleteAnswerDocument(assessmentID uint, questionID uint, answerID uint, documentID uint) error
	GetDocumentPath(assessmentID, questionID, answerID uint, token string) (string, error)
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
	assessmentWorkflowSvc         ThirdPartyAssessmentWorkflowServiceInterface
	complianceResponseSvc         TPRMComplianceResponseServiceInterface
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
	assessmentWorkflowSvc ThirdPartyAssessmentWorkflowServiceInterface,
	complianceResponseSvc TPRMComplianceResponseServiceInterface,
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
		assessmentWorkflowSvc:         assessmentWorkflowSvc,
		complianceResponseSvc:         complianceResponseSvc,
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

	if tpr.ImportedAt != nil {
		dto.ImportedAt = *tpr.ImportedAt
	}

	for i := range dto.Sections {
		dto.Sections[i].QuestionCount = len(dto.Sections[i].Questions)
	}

	// for i, sec := range tpr.Sections {
	// 	for j, qstn := range sec.Questions {
	// 		if qstn.QuestionLibrary.QuestionType != nil {
	// 			dto.Sections[i].Questions[j].ReferenceNumber = qstn.QuestionLibrary.ReferenceNumber
	// 			dto.Sections[i].Questions[j].QuestionTypeTitleEn = qstn.QuestionLibrary.QuestionType.TitleEn
	// 			dto.Sections[i].Questions[j].QuestionTypeTitleAr = qstn.QuestionLibrary.QuestionType.TitleAr
	// 			dto.Sections[i].Questions[j].QuestionTypeType = qstn.QuestionLibrary.QuestionType.Type
	// 			dto.Sections[i].Questions[j].Mandatory = qstn.QuestionLibrary.Mandatory
	// 		}
	// 	}

	// }

	dto.CreatedAt = tpr.CreatedAt

	return &dto, nil
}

func (tpas *ThirdPartyAssessmentService) Create(currentUserID uint, form *dtos.ThirdPartyAssessmentFormDTO) (uint, error) {
	tx := tpas.repo.BeginTransaction()

	err := tpas.assessmentWorkflowSvc.ValidateWorkflowBeforeCreate("tprm_assessment", tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tpas.validateForm(form); err != nil {
		return 0, err
	}

	questionSet, err := tpas.questionSetSvc.GetFullDetailsByID(form.QuestionSetID)
	if err != nil {
		return 0, err
	}

	if len(questionSet.Sections) == 0 {
		return 0, utility.NewForbiddenError("There is no questions under this question set")
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
	appURL := utility.GetVendorAssessmentURL()
	link := fmt.Sprintf("%s/vendor-assessment/%s", appURL, uid)
	tpr := &models.ThirdPartyAssessment{
		ReferenceNumber:              referenceNumber,
		Title:                        form.Title,
		ThirdPartyRegisterID:         form.ThirdPartyRegisterID,
		QuestionSetID:                form.QuestionSetID,
		Message:                      form.Message,
		DueDate:                      dates.DueDate,
		ThirdPartyAssessmentStatusID: assessmentStatus.ID,
		UID:                          &uid,
		Link:                         &link,
		CreatedBYID:                  currentUserID,
	}

	id, err := tpas.repo.Create(tpr, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tpas.assessmentWorkflowSvc.CreateTprmAssessmentWorkflow(id, "tprm_assessment", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, category := range questionSet.Sections {
		sectionDTO := &models.TPRMAssessmentSection{
			ThirdPartyAssessmentID:       id,
			ThirdPartyQuestionCategoryID: category.ThirdPartyQuestionCategoryID,
			TitleEn:                      category.TitleEn,
			TitleAr:                      category.TitleAr,
			Order:                        category.Order,
			CreatedBYID:                  currentUserID,
		}
		sectionID, err := tpas.repo.CreateAssessmentSection(sectionDTO, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, question := range category.Questions {
			questionDTO := &models.TPRMAssessmentSectionQuestion{
				ThirdPartyAssessmentID:      id,
				TPRMAssessmentSectionID:     sectionID,
				ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
				QuestionTitleEn:             question.QuestionTitleEn,
				QuestionTitleAr:             question.QuestionTitleAr,
				Order:                       question.Order,
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
					Score:                       option.Score,
					TPRMComplianceResponseID:    option.TPRMComplianceResponseID,
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
	go func() {
		if form.IsSentEmail {
			err := tpas.ShareLink(id, currentUserID)
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}()
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
		for _, section := range existingAssessment.Sections {
			GetAssessmentSectionDetails, err := tpas.GetAssessmentSectionDetails(id, section.ID)
			if err != nil {
				utility.LogError("Failed to fetch section", map[string]interface{}{
					"error": err.Error(),
				})
			}
			for _, question := range GetAssessmentSectionDetails.Questions {
				if err := tpas.repo.DeleteAnswerOptionsByQuestionID(question.ID, tx); err != nil {
					tx.Rollback()
				}
			}
		}
		if err := tpas.repo.DeleteQuestionsByAssessmentID(id, tx); err != nil {
			tx.Rollback()
		}

		if err := tpas.repo.DeleteSectionsByAssessmentID(id, tx); err != nil {
			tx.Rollback()
		}
		for _, category := range questionSet.Sections {
			sectionDTO := &models.TPRMAssessmentSection{
				ThirdPartyAssessmentID: id,
				TitleEn:                category.TitleEn,
				TitleAr:                category.TitleAr,
				Order:                  category.Order,
			}
			sectionID, err := tpas.repo.CreateAssessmentSection(sectionDTO, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, question := range category.Questions {
				questionDTO := &models.TPRMAssessmentSectionQuestion{
					ThirdPartyAssessmentID:      id,
					ThirdPartyQuestionLibraryID: question.ThirdPartyQuestionLibraryID,
					TPRMAssessmentSectionID:     sectionID,
					QuestionTitleEn:             question.QuestionTitleEn,
					QuestionTitleAr:             question.QuestionTitleAr,
					Order:                       question.Order,
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
		for _, section := range existingAssessment.Sections {
			for _, question := range section.Questions {
				if err := tpas.repo.DeleteAnswerOptionsByQuestionID(question.ID, tx); err != nil {
					tx.Rollback()
					failedIDs = append(failedIDs, id)
					continue
				}
			}
		}
		if err := tpas.repo.DeleteQuestionsByAssessmentID(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}

		var deleteWorkflowIds []*uint
		var deleteWorkflowParticipantIds []*uint
		var deleteWorkflowReviewIds []*uint

		if len(existingAssessment.ThirdPartyAssessmentWorkflows) > 0 {
			for _, workflow := range existingAssessment.ThirdPartyAssessmentWorkflows {
				deleteWorkflowIds = append(deleteWorkflowIds, &workflow.ID)

				if len(workflow.ThirdPartyAssessmentWorkflowParticipants) > 0 {
					for _, participant := range workflow.ThirdPartyAssessmentWorkflowParticipants {
						deleteWorkflowParticipantIds = append(deleteWorkflowParticipantIds, &participant.ID)
					}
				}

				if len(workflow.ThirdPartyAssessmentWorkflowReviews) > 0 {
					for _, reviews := range workflow.ThirdPartyAssessmentWorkflowReviews {
						deleteWorkflowReviewIds = append(deleteWorkflowReviewIds, &reviews.ID)
					}
				}
			}
		}

		if err := tpas.assessmentWorkflowSvc.DeleteTprmAssessmentWorkflowParticipants(deleteWorkflowParticipantIds); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := tpas.assessmentWorkflowSvc.DeleteTprmAssessmentReviews(deleteWorkflowReviewIds); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := tpas.assessmentWorkflowSvc.DeleteTprmAssessmentWorkflows(deleteWorkflowIds); err != nil {
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

func (tpas *ThirdPartyAssessmentService) ShareLink(id uint, currentUserID uint) error {

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

			emailTo := []string{*thirdPartyRegister.Email}
			subject := "Assessment"

			emailData := map[string]interface{}{
				"VendorName":  thirdPartyRegister.Title,
				"ExpiryDate":  thirdPartyAssessment.DueDate.Format("2006-01-02"),
				"RedirectURL": thirdPartyAssessment.Link,
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

func (tpas *ThirdPartyAssessmentService) ExportTPRMAssessment(data *dtos.ThirdPartyAssessmentDetailDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference

	// Define headers based on user language setting
	var headers []string

	headers = []string{"SL No.", "Section", "Reference Number", "Question (EN)", "Question (AR)", "Mandatory", "Answer", "Comment"}

	// Prepare the data rows
	slNo := 1
	rows := [][]interface{}{}
	for _, sec := range data.Sections {
		GetAssessmentSectionDetails, err := tpas.GetAssessmentSectionDetails(data.ID, sec.ID)
		if err != nil {
			utility.LogError("Failed to fetch kri register", map[string]interface{}{
				"error": err.Error(),
			})
		}
		for _, qstn := range GetAssessmentSectionDetails.Questions {
			// Initialize data variables with default values
			var section, referenceNumber, questionEn, questionAr string
			var mandatory bool

			// Assign values based on language preference and check for nil pointers
			section = sec.TitleEn
			referenceNumber = qstn.ReferenceNumber
			questionEn = qstn.QuestionTitleEn
			questionAr = qstn.QuestionTitleAr
			mandatory = qstn.Mandatory

			rows = append(rows, []interface{}{
				slNo,
				section,
				referenceNumber,
				questionEn,
				questionAr,
				mandatory,
				"", // Answer column
				"", //comment column
			})
			slNo++
		}
	}

	dropDownConfig := map[string]bool{
		"Answer": true,
	}

	readOnlyConfig := map[string]bool{
		"SL No.":           true,
		"Section":          true,
		"Reference Number": true,
		"Question (EN)":    true,
		"Question (AR)":    true,
		"Mandatory":        true,
	}

	responses, _ := tpas.complianceResponseSvc.GetAll("", "", "", "active")
	responseHeaders := []string{"TitleEn"}
	responseData := [][]interface{}{}
	for _, cat := range responses {
		responseData = append(responseData, []interface{}{cat.TitleEn})
	}
	dueDate := ""
	if data.DueDateString != nil {
		dueDate = *data.DueDateString
	}
	introductionHeaders := []string{"1. Assessment Details", ""}
	introductionRows := [][]interface{}{
		{"Assessment Title", data.Title},
		{"Third Party / Vendor Name", data.ThirdPartyRegister.Title},
		{"Question Set", data.QuestionSet.TitleEn},
		{"Assessment Reference Number", data.ReferenceNumber},
		{"Due Date", dueDate},
		{"Assessment Created On", "" + data.CreatedAt.Format("02-01-2006 15:04")},
		{"Assessment Created By",
			utility.GetStringValue(data.CreatedByUser.FirstNameEN) + " " +
				utility.GetStringValue(data.CreatedByUser.LastNameEN),
		},

		{"", ""}, // blank row

		{"2. Instructions for Completing the Assessment", ""},
		{"1. All assessment questions are listed in Sheet 2 – Assessment Questions.", ""},
		{"2. Enter your responses in the Answer column using the dropdown provided (e.g., Comply / Partially Comply / Not Comply).", ""},
		{"3. Provide explanations or attach supporting evidence references in the Comment column where applicable.", ""},
		{"4. Mandatory questions (marked as TRUE) must be answered; the system may reject submissions with missing mandatory answers.", ""},
		{"5. Do not modify the sheet structure (do not delete or rename columns, questions, or reference numbers).", ""},
		{"6. After completing all sections, save the file and return it as instructed by your assessment coordinator.", ""},
	}

	// answerOptionHeaders := []interface{}{"Option (EN)", "Option (AR)"}
	// answerOptionRows := [][]interface{}{
	// 	{"Comply", "امتثال"},
	// 	{"Partially Comply", "امتثال جزئي"},
	// 	{"Not Comply", "عدم امتثال"},
	// 	{"Not Applicable", "غير قابل للتطبيق"},
	// }

	// Generate Excel file
	file, err := utility.GenerateExcelFileWithSheetsAndDropDownAnswer(map[string][][]interface{}{
		"Introduction": append([][]interface{}{utility.ConvertToInterfaceSlice(introductionHeaders)}, introductionRows...),
		"Template":     append([][]interface{}{utility.ConvertToInterfaceSlice(headers)}, rows...),
		"Answer":       append([][]interface{}{utility.ConvertToInterfaceSlice(responseHeaders)}, responseData...),
	}, dropDownConfig, readOnlyConfig)
	if err != nil {
		return nil, err
	}

	filePath := file.Name()
	file.Close()

	// Re-open for merge + style
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
	})

	instructionTitleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  11,
			Color: "#FFFFFF", // white
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#165D31"}, // green
			Pattern: 1,
		},
	})

	// Merge A10:B10 for clean heading
	f.MergeCell("Introduction", "A1", "B1")

	// Apply bold style
	f.SetCellStyle("Introduction", "A1", "B1", instructionTitleStyle)

	// Merge A10:B10 for clean heading
	f.MergeCell("Introduction", "A10", "B10")

	// Apply bold style
	f.SetCellStyle("Introduction", "A10", "B10", instructionTitleStyle)

	f.SetCellStyle("Introduction", "A2", "A8", labelStyle)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   false,
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})

	// Apply style to header row (Row 1)
	headerStart, _ := excelize.CoordinatesToCellName(1, 1)
	headerEnd, _ := excelize.CoordinatesToCellName(len(headers), 1)
	f.SetCellStyle("Template", headerStart, headerEnd, headerStyle)

	currentSection := ""
	startRow := 2

	for i, row := range rows {
		rowIndex := i + 2
		section := row[1].(string)

		if section != currentSection {
			if currentSection != "" && startRow < rowIndex-1 {
				startCell, _ := excelize.CoordinatesToCellName(2, startRow)
				endCell, _ := excelize.CoordinatesToCellName(2, rowIndex-1)
				f.MergeCell("Template", startCell, endCell)
			}
			currentSection = section
			startRow = rowIndex
		}
	}

	// Merge last group
	lastRow := len(rows) + 1
	if startRow < lastRow {
		startCell, _ := excelize.CoordinatesToCellName(2, startRow)
		endCell, _ := excelize.CoordinatesToCellName(2, lastRow)
		f.MergeCell("Template", startCell, endCell)
	}

	sectionStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})

	startCell, _ := excelize.CoordinatesToCellName(2, 2)
	endCell, _ := excelize.CoordinatesToCellName(2, lastRow)
	f.SetCellStyle("Template", startCell, endCell, sectionStyle)

	// Save
	if err := f.SaveAs(filePath); err != nil {
		return nil, err
	}

	f.Close()

	// Return file
	finalFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return finalFile, nil
}

// Validate Excel Template
func (tpas *ThirdPartyAssessmentService) ValidateTPRMAssessmentTemplate(assessmentID uint, file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFileWithIntroduction(file)
	if err != nil {
		return nil, err
	}
	inputHeaders := []string{"SL No.", "Section", "Reference Number", "Question (EN)", "Question (AR)", "Mandatory", "Answer", "Comment"}
	snakeHeaders := []string{"sl_no", "section", "reference_number", "question_en", "question_ar", "mandatory", "answer", "comment"}

	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Assessment Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, h := range snakeHeaders {
			if i < len(row) {
				record[h] = row[i]
			} else {
				record[h] = ""
			}
		}
		input = append(input, record)
	}

	return tpas.validateTPRMAssessment(assessmentID, input)
}

// Validate JSON input
func (tpas *ThirdPartyAssessmentService) ValidateImportTPRMAssessment(assessmentID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	return tpas.validateTPRMAssessment(assessmentID, input)
}

// Import
func (tpas *ThirdPartyAssessmentService) ImportTPRMAssessment(assessmentID uint, input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validated, err := tpas.validateTPRMAssessment(assessmentID, input)
	if err != nil {
		return nil, err
	}

	assessment, err := tpas.GetByID(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("error fetching assessment details: %v", err)
	}

	for _, r := range validated {
		if errs, ok := r["errors"].(map[string]interface{}); ok && len(errs) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validated {

		answer := &models.TPRMAssessmentAnswer{
			TPRMAssessmentQuestionID: record["tprm_assessment_question_id"].(uint),
			Comment:                  record["comment"].(string),
		}

		answerID, err := tpas.repo.CreateAssessmentAnswer(answer)
		if err != nil {
			return nil, err
		}

		if record["answer_option_id"] == "" {
			optionMap := make(map[uint]map[string]uint)

			GetAssessmentQuestionDetails, err := tpas.repo.GetAssessmentSectionQuestionDetails(record["tprm_assessment_question_id"].(uint))
			if err != nil {
				utility.LogError("Failed to fetch question", map[string]interface{}{
					"error": err.Error(),
				})
			}

			if _, exists := optionMap[record["tprm_assessment_question_id"].(uint)]; !exists {
				optionMap[record["tprm_assessment_question_id"].(uint)] = make(map[string]uint)
			}

			for _, ans := range GetAssessmentQuestionDetails.AnswerOptions {
				if ans.AnswerOptionTitleEn != nil {
					optionMap[record["tprm_assessment_question_id"].(uint)][strings.TrimSpace(*ans.AnswerOptionTitleEn)] = ans.ID
				}
			}

			if id, exists := optionMap[record["tprm_assessment_question_id"].(uint)]["Not Comply"]; exists {
				record["answer"] = "Not Comply"
				record["answer_option_id"] = id
			}
		}

		questionAnswerOption, err := tpas.repo.GetAnswerOptionByID(record["answer_option_id"].(uint))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, utility.NewResourceNotFoundError(err, record["answer_option_id"].(uint))
			}
			return nil, err
		}

		answerOption := &models.TPRMAssessmentAnswerOption{
			TPRMAssessmentQuestionID: record["tprm_assessment_question_id"].(uint),
			AnswerOptionID:           record["answer_option_id"].(uint),
			TPRMAssessmentAnswerID:   answerID,
			TPRMComplianceResponseID: questionAnswerOption.TPRMComplianceResponseID,
		}

		_, err = tpas.repo.CreateAssessmentAnswerOption(answerOption)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		assessment.ImportedBYID = &currentUserID
		assessment.ImportedAt = &now

		if err := tpas.repo.Update(assessment); err != nil {
			return nil, err
		}

		createdIDs = append(createdIDs, answerID)
	}

	return createdIDs, nil
}

func (tpas *ThirdPartyAssessmentService) validateTPRMAssessment(assessmentID uint, input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"sl_no", "section", "reference_number", "question_en", "question_ar", "mandatory", "answer", "comment"}

	assessmentDetails, err := tpas.GetFullDetailsByID(assessmentID)
	if err != nil {
		return nil, fmt.Errorf("error fetching assessment details: %v", err)
	}

	// --- Validate each record ---
	var validated []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processed := make(map[string]interface{})

		// Copy original input
		for k, v := range record {
			processed[k] = v
		}

		questionMap := make(map[string]uint)
		// optionMap := make(map[uint]map[string]uint)
		optionMap := make(map[uint]map[string]uint)
		for _, sec := range assessmentDetails.Sections {
			GetAssessmentSectionDetails, err := tpas.GetAssessmentSectionDetails(assessmentID, sec.ID)
			if err != nil {
				utility.LogError("Failed to fetch section", map[string]interface{}{
					"error": err.Error(),
				})
			}
			for _, qstn := range GetAssessmentSectionDetails.Questions {
				GetAssessmentQuestionDetails, err := tpas.GetAssessmentSectionQuestionDetails(assessmentID, sec.ID, qstn.ID)
				if err != nil {
					utility.LogError("Failed to fetch question", map[string]interface{}{
						"error": err.Error(),
					})
				}
				questionMap[strings.TrimSpace(qstn.ReferenceNumber)] = qstn.ID

				if _, exists := optionMap[qstn.ID]; !exists {
					optionMap[qstn.ID] = make(map[string]uint)
				}

				for _, ans := range GetAssessmentQuestionDetails.AnswerOptions {
					if ans.AnswerOptionTitleEn != nil {
						optionMap[qstn.ID][strings.TrimSpace(*ans.AnswerOptionTitleEn)] = ans.ID
					}
				}
			}
		}
		// question id
		if val, ok := record["reference_number"].(string); !ok || val == "" {
			rowErrors["reference_number"] = "Field is required"
		} else if qid, exists := questionMap[strings.TrimSpace(val)]; exists {
			processed["tprm_assessment_question_id"] = qid
		} else {
			rowErrors["reference_number"] = fmt.Sprintf("Invalid reference number: %s", val)
		}

		if val, ok := record["question_ar"].(string); !ok || val == "" {
			rowErrors["question_ar"] = "Field is required"
		}

		// var questionID uint

		// questionID = processed["tprm_assessment_question_id"]
		var questionID uint
		switch v := processed["tprm_assessment_question_id"].(type) {
		case uint:
			questionID = v
		case float64:
			questionID = uint(v)
		}

		if &questionID != nil {
			if val, ok := record["answer"].(string); !ok || val == "" {
				if record["mandatory"] == "TRUE" {
					rowErrors["answer"] = "Answer is Mandatory"
				}
				// if id, exists := optionMap[questionID]["Not Comply"]; exists {
				// 	record["answer"] = "Not Comply"
				// 	processed["answer"] = "Not Comply"
				// 	processed["answer_option_id"] = id
				// }
			} else if id, exists := optionMap[questionID][strings.TrimSpace(val)]; exists {
				processed["answer_option_id"] = id
			} else {
				rowErrors["answer"] = fmt.Sprintf("Invalid answer: %s", val)
			}
		}

		// Collect errors
		if len(rowErrors) == 0 {
			processed["errors"] = nil
		} else {
			for _, h := range headers {
				if _, exists := rowErrors[h]; !exists {
					rowErrors[h] = "" // keep consistent structure
				}
			}
			processed["errors"] = rowErrors
		}

		validated = append(validated, processed)
	}

	return validated, nil
}

func (tpas *ThirdPartyAssessmentService) MarkAsSend(assessmentID uint, currentUserID uint) error {
	tx := tpas.repo.BeginTransaction()

	assessment, err := tpas.repo.GetByID(assessmentID)
	if err != nil {
		return err
	}

	now := time.Now()
	assessment.IsSend = true
	assessment.SentBYID = &currentUserID
	assessment.SentAt = &now
	assessment.UpdatedBYID = &currentUserID

	if err := tpas.repo.Update(assessment, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (s *ThirdPartyAssessmentService) GetAssessmentAnswerSummary(assessmentID uint) ([]dtos.AnswerSummaryItemDTO, error) {
	return s.repo.GetAnswerSummary(assessmentID)
}

func (tpas *ThirdPartyAssessmentService) ExportTPRMAssessmentReport(assessmentID uint) (string, error) {
	assessment, err := tpas.GetFullDetailsByID(assessmentID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch assessment: %w", err)
	}

	var assessmentSections []dtos.TPRMAssessmentSectionDTO

	for _, section := range assessment.Sections {
		sectionDetails, err := tpas.GetAssessmentSectionDetails(assessmentID, section.ID)
		if err != nil {
			return "", err
		}

		for idx, sectionQuestion := range sectionDetails.Questions {
			questionDetails, err := tpas.GetAssessmentSectionQuestionDetails(assessmentID, section.ID, sectionQuestion.ID)
			if err != nil {
				return "", err
			}

			sectionDetails.Questions[idx].Answer = questionDetails.Answer

		}

		assessmentSections = append(assessmentSections, *sectionDetails)
	}

	assessment.Sections = assessmentSections

	pdfData := map[string]interface{}{
		"Assessment":  assessment,
		"GeneratedAt": time.Now().Format("02 Jan 2006"),
		"Sections":    assessment.Sections,
		"AppURL":      utility.GetAppURL(),
	}

	coverHTML, err := utility.GeneratePDFContent(pdfData, []string{"tprm/templates/pdf/tprm_assessment_cover_page.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate assessment cover page: %w", err)
	}

	contentHTML, err := utility.GeneratePDFContent(pdfData, []string{"tprm/templates/pdf/tprm_assessment_content.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate assessment content page: %w", err)
	}

	headerHTML, err := utility.GeneratePDFContent(pdfData, []string{"tprm/templates/pdf/tprm_assessment_header.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate assessment header: %w", err)
	}

	footerHTML, err := utility.GeneratePDFContent(pdfData, []string{"tprm/templates/pdf/tprm_assessment_footer.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate assessment footer: %w", err)
	}

	fileName := fmt.Sprintf("TPRMAssessmentReport_%d.pdf", assessmentID)

	if err := utility.GeneratePDFWithCoverHeaderFooter(coverHTML, contentHTML, fileName, headerHTML, footerHTML); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return fileName, nil
}

func (tpas *ThirdPartyAssessmentService) ExportThirdPartyAssessments(data []*dtos.ThirdPartyAssessmentListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := tpas.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string

	// Select headers based on language
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title (AR)",
			"Third Party Title",
			"Question Set (AR)",
			"Status (AR)",
			"Created By (AR)",
			"Due Date",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title (EN)",
			"Third Party Title",
			"Question Set (EN)",
			"Status (EN)",
			"Created By (EN)",
			"Due Date",
		}
	}

	rows := make([][]interface{}, len(data))

	for i, item := range data {
		if userSetting.Language == "ar" {

			createdByAr := utility.GetStringValue(item.CreatedUserFirstNameAr) + " " + utility.GetStringValue(item.CreatedUserLastNameAr)

			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.Title,
				item.ThirdPartyRegisterTitle,
				item.QuestionSetTitleAr,
				item.ThirdPartyAssessmentStatusTitleAr,
				createdByAr,
				utility.GetStringValue(item.DueDate),
			}

		} else {

			createdByEn := utility.GetStringValue(item.CreatedUserFirstNameEn) + " " + utility.GetStringValue(item.CreatedUserLastNameEn)

			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.Title,
				item.ThirdPartyRegisterTitle,
				item.QuestionSetTitleEn,
				item.ThirdPartyAssessmentStatusTitleEn,
				createdByEn,
				utility.GetStringValue(item.DueDate),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (tpas *ThirdPartyAssessmentService) GetQuestionByID(id uint) (*models.TPRMAssessmentSectionQuestion, error) {
	question, err := tpas.repo.GetQuestionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return question, nil
}

func (tpas *ThirdPartyAssessmentService) UpdateAnswer(id uint, questionId uint, currentUserID uint, form *dtos.ThirdPartyAssessmentAnswerFormDTO) error {
	existing, err := tpas.GetQuestionByID(questionId)
	if err != nil {
		return err
	}
	tx := tpas.repo.BeginTransaction()

	if existing.Answer != nil {
		existing.Answer.Comment = *form.Comment

		err = tpas.repo.UpdateAnswer(existing.Answer, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tpas.repo.DeleteAnswerOptionsByAnswerID(existing.Answer.ID, tx); err != nil {
			tx.Rollback()
			return err
		}

		questionAnswerOption, err := tpas.repo.GetAnswerOptionByQuestionID(questionId, form.OptionId)
		if err != nil {
			return err
		}

		answerOption := &models.TPRMAssessmentAnswerOption{
			TPRMAssessmentQuestionID: questionId,
			AnswerOptionID:           questionAnswerOption.ID,
			TPRMComplianceResponseID: &form.OptionId,
			TPRMAssessmentAnswerID:   existing.Answer.ID,
		}

		_, err = tpas.repo.CreateAssessmentAnswerOption(answerOption, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		answer := &models.TPRMAssessmentAnswer{
			TPRMAssessmentQuestionID: questionId,
			Comment:                  *form.Comment,
		}

		answerID, err := tpas.repo.CreateAssessmentAnswer(answer)
		if err != nil {
			return err
		}

		questionAnswerOption, err := tpas.repo.GetAnswerOptionByQuestionID(questionId, form.OptionId)
		if err != nil {
			return err
		}

		answerOption := &models.TPRMAssessmentAnswerOption{
			TPRMAssessmentQuestionID: questionId,
			AnswerOptionID:           questionAnswerOption.ID,
			TPRMComplianceResponseID: &form.OptionId,
			TPRMAssessmentAnswerID:   answerID,
		}

		_, err = tpas.repo.CreateAssessmentAnswerOption(answerOption, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (tpas *ThirdPartyAssessmentService) GetAssessmentSectionDetails(assessmentId uint, sectionId uint) (*dtos.TPRMAssessmentSectionDTO, error) {
	var dto dtos.TPRMAssessmentSectionDTO

	tpr, err := tpas.repo.GetAssessmentSectionDetails(sectionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionId)
		}
		return nil, err
	}

	err = mapstructure.Decode(tpr, &dto)
	if err != nil {
		return nil, err
	}

	for j, qstn := range tpr.Questions {
		if qstn.QuestionLibrary.QuestionType != nil {
			dto.Questions[j].ReferenceNumber = qstn.QuestionLibrary.ReferenceNumber
			dto.Questions[j].QuestionTypeTitleEn = qstn.QuestionLibrary.QuestionType.TitleEn
			dto.Questions[j].QuestionTypeTitleAr = qstn.QuestionLibrary.QuestionType.TitleAr
			dto.Questions[j].QuestionTypeType = qstn.QuestionLibrary.QuestionType.Type
			dto.Questions[j].Mandatory = qstn.QuestionLibrary.Mandatory
		}
	}

	return &dto, nil
}

func (tpas *ThirdPartyAssessmentService) GetAssessmentSectionQuestionDetails(assessmentId uint, sectionId uint, questionId uint) (*dtos.TPRMAssessmentQuestionDetailDTO, error) {
	var dto dtos.TPRMAssessmentQuestionDetailDTO

	tpr, err := tpas.repo.GetAssessmentSectionQuestionDetails(questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, sectionId)
		}
		return nil, err
	}

	err = mapstructure.Decode(tpr, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (tpas *ThirdPartyAssessmentService) GetAllAnswerDocuments(assessmentID uint, questionID uint, answerID uint) ([]dtos.AnswerDocumentListDTO, error) {
	if _, err := tpas.GetByID(assessmentID); err != nil {
		return nil, err
	}
	return tpas.repo.GetAllAnswerDocuments(answerID)
}

func (tpas *ThirdPartyAssessmentService) CreateAnswerDocument(assessmentID uint, questionID uint, answerID uint, userID uint, formDTO *dtos.AnswerDocumentFormDTO) error {

	if _, err := tpas.GetByID(assessmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("assessment does not exist for ID: %d", assessmentID), "not_exist", "assessment_id")
	}

	tx := tpas.repo.BeginTransaction()

	url := utility.GenerateDocumentURL("third-party", "assessment", strconv.FormatUint(uint64(assessmentID), 10), "question", strconv.FormatUint(uint64(questionID), 10), "documents")

	for _, doc := range formDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		userDocument := &models.TPRMAssessmentAnswerDocument{
			TPRMAssessmentAnswerID:   answerID,
			TPRMAssessmentQuestionID: questionID,
			URL:                      url,
			Name:                     &doc.FileName,
			Size:                     &doc.Size,
			Type:                     &doc.Type,
			Token:                    &encryptedFilename,
			OriginalName:             &doc.OriginalName,
		}

		_, err = tpas.repo.CreateAnswerDocument(userDocument)
		if err != nil {
			tx.Rollback()
			return err
		}

		storagePath := filepath.Join("third-party", "assessment", strconv.FormatUint(uint64(assessmentID), 10), "question", strconv.FormatUint(uint64(questionID), 10), "documents")
		if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (tpas *ThirdPartyAssessmentService) DeleteAnswerDocument(assessmentID uint, questionID uint, answerID uint, documentID uint) error {
	// Validate control existence
	if _, err := tpas.GetByID(assessmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("assessment does not exist for ID: %d", assessmentID), "not_exist", "assessment_id")
	}

	tx := tpas.repo.BeginTransaction()

	if err := tpas.repo.DeleteAnswerDocument(documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpas *ThirdPartyAssessmentService) GetDocumentPath(assessmentID, questionID, answerID uint, token string) (string, error) {
	question, err := tpas.repo.GetQuestionByID(questionID)
	if err != nil {
		fmt.Printf("Error fetching status update by ID: %v\n", err)
		return "", err
	}

	if len(question.Answer.Documents) == 0 {
		return "", fmt.Errorf("no documents found")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		fmt.Printf("Token decryption failed: %v\n", err)
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)
	fmt.Printf("Decrypted filename: %s\n", decryptedFilename)

	storagePath := filepath.Join("third-party", "assessment", strconv.FormatUint(uint64(assessmentID), 10), "question", strconv.FormatUint(uint64(questionID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	fmt.Printf("Constructed file path: %s\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("File does not exist at path: %s\n", fullPath)
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		fmt.Printf("Error accessing file at path %s: %v\n", fullPath, err)
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
