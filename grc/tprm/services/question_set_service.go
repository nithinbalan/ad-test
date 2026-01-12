package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type QuestionSetServiceInterface interface {
	GetAll(params *dtos.QuestionSetListQueryParamsDTO) (*dtos.QuestionSetListResponseDTO, error)
	GetByID(id uint) (*models.QuestionSet, error)
	GetFullDetailsByID(id uint) (*dtos.QuestionSetDetailDTO, error)
	Create(currentUserID uint, form *dtos.QuestionSetFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.QuestionSetFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)

	AddQuestionsToSet(questionSetID, currentUserID uint, dto *dtos.QuestionSetAddQuestionsDTO) error
	DeleteQuestionFromSet(questionSetID uint, questionID uint) error
	ReorderSection(sections []dtos.ReorderDTO) error
	ReorderQuestion(questions []dtos.ReorderDTO) error

	ExportQuestionSet(data *dtos.QuestionSetDetailDTO, currentUserID uint) (*os.File, error)
	ExportQuestionSets(data []*dtos.QuestionSetListDTO, currentUserID uint) (*os.File, error)
}

type QuestionSetService struct {
	repo                         tprmRepo.QuestionSetRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	thirdPartyQuestionLibrarySvc ThirdPartyQuestionLibraryServiceInterface
	userSettingSvc               organizationServices.UserSettingServiceInterface
}

func NewQuestionSetService(
	repo tprmRepo.QuestionSetRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	thirdPartyQuestionLibrarySvc ThirdPartyQuestionLibraryServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) *QuestionSetService {
	return &QuestionSetService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		thirdPartyQuestionLibrarySvc: thirdPartyQuestionLibrarySvc,
		userSettingSvc:               userSettingSvc,
	}
}

func (svc *QuestionSetService) GetAll(params *dtos.QuestionSetListQueryParamsDTO) (*dtos.QuestionSetListResponseDTO, error) {
	return svc.repo.GetAll(params)
}

func (svc *QuestionSetService) GetByID(id uint) (*models.QuestionSet, error) {
	entity, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return entity, nil
}

func (svc *QuestionSetService) GetFullDetailsByID(id uint) (*dtos.QuestionSetDetailDTO, error) {
	entity, err := svc.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	dto := &dtos.QuestionSetDetailDTO{
		ID:          entity.ID,
		TitleEn:     entity.TitleEn,
		TitleAr:     entity.TitleAr,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	dto.Sections = make([]dtos.QuestionSetSectionDTO, 0)

	for _, sec := range entity.Sections {

		sectionDTO := dtos.QuestionSetSectionDTO{
			ID:                           sec.ID,
			TitleEn:                      sec.TitleEn,
			TitleAr:                      sec.TitleAr,
			Order:                        sec.Order,
			ThirdPartyQuestionCategoryID: sec.ThirdPartyQuestionCategoryID,
			Questions:                    []dtos.QuestionSetQuestionDetailDTO{},
		}

		for _, q := range sec.Questions {

			questionDTO := dtos.QuestionSetQuestionDetailDTO{
				ID:                          q.ID,
				QuestionSetID:               q.QuestionSetID,
				ThirdPartyQuestionLibraryID: q.ThirdPartyQuestionLibraryID,
				QuestionTitleEn:             q.QuestionTitleEn,
				QuestionTitleAr:             q.QuestionTitleAr,
				Order:                       q.Order,
			}

			if q.QuestionLibrary != nil {
				questionDTO.ReferenceNumber = q.QuestionLibrary.ReferenceNumber
				questionDTO.Mandatory = q.QuestionLibrary.Mandatory
				questionDTO.ThirdPartyQuestionCategoryID = q.QuestionLibrary.ThirdPartyQuestionCategoryID
				questionDTO.ThirdPartyQuestionTypeID = q.QuestionLibrary.ThirdPartyQuestionTypeID

				if q.QuestionLibrary.QuestionCategory != nil {
					questionDTO.QuestionCategoryTitleEn = q.QuestionLibrary.QuestionCategory.TitleEn
					questionDTO.QuestionCategoryTitleAr = q.QuestionLibrary.QuestionCategory.TitleAr
				}

				if q.QuestionLibrary.QuestionType != nil {
					questionDTO.QuestionTypeTitleEn = q.QuestionLibrary.QuestionType.TitleEn
					questionDTO.QuestionTypeTitleAr = q.QuestionLibrary.QuestionType.TitleAr
					questionDTO.QuestionTypeType = q.QuestionLibrary.QuestionType.Type
				}

				for _, opt := range q.QuestionLibrary.AnswerOptions {
					optDTO := dtos.ThirdPartyQuestionLibraryAnswerOptionDetailDTO{
						ID:                  opt.ID,
						AnswerOptionTitleEn: opt.AnswerOptionTitleEn,
						AnswerOptionTitleAr: opt.AnswerOptionTitleAr,
						RiskLevelID:         opt.RiskLevelID,
						RangeFrom:           opt.RangeFrom,
						RangeTo:             opt.RangeTo,
					}

					if opt.TPRMComplianceResponse != nil {
						optDTO.TPRMComplianceResponse = &dtos.TPRMComplianceResponseDetailDTO{}
						optDTO.Score = &opt.TPRMComplianceResponse.Score
						optDTO.TPRMComplianceResponseID = opt.TPRMComplianceResponseID
						optDTO.TPRMComplianceResponse.ID = opt.TPRMComplianceResponse.ID
						optDTO.TPRMComplianceResponse.TitleEn = opt.TPRMComplianceResponse.TitleEn
						optDTO.TPRMComplianceResponse.TitleAr = opt.TPRMComplianceResponse.TitleAr
						optDTO.TPRMComplianceResponse.ColorCode = *opt.TPRMComplianceResponse.ColorCode
						optDTO.TPRMComplianceResponse.Score = opt.TPRMComplianceResponse.Score
					}

					if opt.RiskLevel != nil {
						optDTO.RiskLevel.ID = opt.RiskLevel.ID
						optDTO.RiskLevel.TitleEn = opt.RiskLevel.TitleEn
						optDTO.RiskLevel.TitleAr = opt.RiskLevel.TitleAr
						optDTO.RiskLevel.ColorCode = opt.RiskLevel.ColorCode
						optDTO.RiskLevel.Type = opt.RiskLevel.Type
					}

					questionDTO.AnswerOptions = append(questionDTO.AnswerOptions, optDTO)
				}
			}

			sectionDTO.Questions = append(sectionDTO.Questions, questionDTO)
		}

		dto.Sections = append(dto.Sections, sectionDTO)
	}

	dto.TotalQuestions = len(entity.Questions)
	return dto, nil
}

func (svc *QuestionSetService) Create(currentUserID uint, form *dtos.QuestionSetFormDTO) (uint, error) {
	tx := svc.repo.BeginTransaction()

	entity := &models.QuestionSet{
		TitleEn:     form.TitleEn,
		TitleAr:     form.TitleAr,
		Description: form.Description,
		CreatedBYID: currentUserID,
	}

	id, err := svc.repo.Create(entity, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (svc *QuestionSetService) Update(id uint, currentUserID uint, form *dtos.QuestionSetFormDTO) error {
	tx := svc.repo.BeginTransaction()

	entity, err := svc.GetByID(id)
	if err != nil {
		return err
	}

	entity.TitleEn = form.TitleEn
	entity.TitleAr = form.TitleAr
	entity.Description = form.Description
	entity.UpdatedBYID = &currentUserID

	err = svc.repo.Update(entity, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (svc *QuestionSetService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failed []uint

	for _, id := range ids {
		tx := svc.repo.BeginTransaction()

		if _, err := svc.GetByID(id); err != nil {
			failed = append(failed, id)
			continue
		}

		// Delete child questions first
		if err := svc.repo.DeleteByQuestionSetID(id, tx); err != nil {
			tx.Rollback()
			failed = append(failed, id)
		}

		if err := svc.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failed = append(failed, id)
			continue
		}

		tx.Commit()
	}

	return failed, nil
}

func (svc *QuestionSetService) AddQuestionsToSet(questionSetID, currentUserID uint, dto *dtos.QuestionSetAddQuestionsDTO) error {

	tx := svc.repo.BeginTransaction()

	questionSet, err := svc.repo.GetFullDetailsByID(questionSetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return utility.NewResourceNotFoundError(err, questionSetID)
		}
		tx.Rollback()
		return err
	}

	sectionMap := make(map[uint]uint)     // categoryID -> sectionID
	sectionOrderMap := make(map[uint]int) // sectionID -> max question order
	maxSectionOrder := 0                  // for new sections

	for _, sec := range questionSet.Sections {
		sectionMap[sec.ThirdPartyQuestionCategoryID] = sec.ID

		if sec.Order > maxSectionOrder {
			maxSectionOrder = sec.Order
		}

		// find max question order in this section
		for _, q := range sec.Questions {
			if q.Order > sectionOrderMap[sec.ID] {
				sectionOrderMap[sec.ID] = q.Order
			}
		}
	}

	var questions []*models.QuestionSetSectionQuestion

	for _, libID := range dto.QuestionLibraryIDs {

		lib, err := svc.thirdPartyQuestionLibrarySvc.GetFullDetailsByID(libID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "not_exist", fmt.Sprintf("question_library_id_%d", libID))
		}

		category := lib.QuestionCategory

		sectionID, exists := sectionMap[category.ID]
		if !exists {

			maxSectionOrder++

			newSection := &models.QuestionSetSection{
				QuestionSetID:                questionSetID,
				ThirdPartyQuestionCategoryID: category.ID,
				TitleEn:                      category.TitleEn,
				TitleAr:                      category.TitleAr,
				Order:                        maxSectionOrder,
				CreatedBYID:                  currentUserID,
			}

			sectionID, err = svc.repo.CreateSection(newSection, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			sectionMap[category.ID] = sectionID
			sectionOrderMap[sectionID] = 0 // initialize
		}

		sectionOrderMap[sectionID]++
		nextOrder := sectionOrderMap[sectionID]

		questions = append(questions, &models.QuestionSetSectionQuestion{
			QuestionSetID:               questionSetID,
			QuestionSetSectionID:        sectionID,
			ThirdPartyQuestionLibraryID: libID,
			QuestionTitleEn:             lib.QuestionTitleEn,
			QuestionTitleAr:             lib.QuestionTitleAr,
			Order:                       nextOrder,
			CreatedBYID:                 currentUserID,
		})
	}

	if len(questions) > 0 {
		if err := svc.repo.CreateBulk(questions, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (svc *QuestionSetService) DeleteQuestionFromSet(questionSetID uint, questionID uint) error {
	sectionID, err := svc.repo.GetSectionIDByQuestionID(questionID)
	if err != nil {
		return err
	}
	tx := svc.repo.BeginTransaction()

	if err := svc.repo.DeleteByID(questionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	questions, err := svc.repo.GetQuestionsBySectionID(sectionID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(questions) == 0 {
		if err := svc.repo.DeleteSection(sectionID, tx); err != nil {
			tx.Rollback()
			return err
		}

		questionSet, err := svc.repo.GetFullDetailsByID(questionSetID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// 5. Reorder remaining sections
		for index, section := range questionSet.Sections {
			newOrder := index + 1
			section.Order = newOrder
			err := svc.repo.UpdateSection(section, tx)
			if err != nil {
				return err
			}
		}
	}
	for i, q := range questions {
		qstn, err := svc.repo.GetByQuestionID(q.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		newOrder := i + 1
		if q.Order != newOrder {
			qstn.Order = newOrder
			if err := svc.repo.UpdateQuestion(qstn, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (svc *QuestionSetService) ReorderSection(sections []dtos.ReorderDTO) error {
	for _, section := range sections {
		existingSection, err := svc.repo.GetSectionByID(section.ID)
		if err != nil {
			return nil
		}
		existingSection.Order = section.Order
		err = svc.repo.UpdateSection(existingSection)
		if err != nil {
			return nil
		}

	}

	return nil
}

func (svc *QuestionSetService) ReorderQuestion(questions []dtos.ReorderDTO) error {
	for _, question := range questions {
		existingQuestion, err := svc.repo.GetByQuestionID(question.ID)
		if err != nil {
			return nil
		}
		existingQuestion.Order = question.Order
		err = svc.repo.UpdateQuestion(existingQuestion)
		if err != nil {
			return nil
		}

	}

	return nil
}

func (svc *QuestionSetService) ExportQuestionSet(data *dtos.QuestionSetDetailDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := svc.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{"SL No.", "Section", "Question", "Mandatory", "Question Type", "Options"}
	} else {
		headers = []string{"SL No.", "Section", "Question", "Mandatory", "Question Type", "Options"}
	}

	// Prepare the data rows
	slNo := 1
	rows := [][]interface{}{}
	for _, sec := range data.Sections {
		for _, qstn := range sec.Questions {
			// Initialize data variables with default values
			var options []string
			var section, question, questionType string
			var mandatory bool

			// Assign values based on language preference and check for nil pointers
			if userSetting.Language == "ar" {
				section = sec.TitleAr
				question = qstn.QuestionTitleAr
				mandatory = qstn.Mandatory
				questionType = qstn.QuestionTypeTitleAr

				for _, opt := range qstn.AnswerOptions {
					options = append(options, *opt.AnswerOptionTitleAr)
				}

			} else {
				section = sec.TitleEn
				question = qstn.QuestionTitleEn
				mandatory = qstn.Mandatory
				questionType = qstn.QuestionTypeTitleEn
				for _, opt := range qstn.AnswerOptions {
					options = append(options, *opt.AnswerOptionTitleEn)
				}
			}

			joinedOptions := strings.Join(options, ", ")

			rows = append(rows, []interface{}{
				slNo,
				section,
				question,
				mandatory,
				questionType,
				joinedOptions,
			})
			slNo++
		}
	}

	// Generate the Excel file
	file, err := utility.GenerateExcelFile(headers, rows)
	if err != nil {
		return nil, err
	}

	filePath := file.Name()
	file.Close()

	// Re-open for merge & styling
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	currentSection := ""
	startRow := 2

	for i, row := range rows {
		rowIndex := i + 2
		section := row[1].(string)

		if section != currentSection {
			if currentSection != "" && startRow < rowIndex-1 {
				// merge previous section block
				startCell, _ := excelize.CoordinatesToCellName(2, startRow)
				endCell, _ := excelize.CoordinatesToCellName(2, rowIndex-1)
				f.MergeCell("Sheet1", startCell, endCell)
			}

			currentSection = section
			startRow = rowIndex
		}
	}

	// Merge final section
	lastRow := len(rows) + 1
	if startRow < lastRow {
		startCell, _ := excelize.CoordinatesToCellName(2, startRow)
		endCell, _ := excelize.CoordinatesToCellName(2, lastRow)
		f.MergeCell("Sheet1", startCell, endCell)
	}

	sectionStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical: "center",
			WrapText: true,
		},
	})

	secStart, _ := excelize.CoordinatesToCellName(2, 2)
	secEnd, _ := excelize.CoordinatesToCellName(2, lastRow)
	f.SetCellStyle("Sheet1", secStart, secEnd, sectionStyle)

	if err := f.SaveAs(filePath); err != nil {
		return nil, err
	}
	f.Close()

	finalFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return finalFile, nil

}

func (svc *QuestionSetService) ExportQuestionSets(data []*dtos.QuestionSetListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := svc.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Title",
			"Description",
			"Created By",
			"Created At",
			"Total Questions",
		}
	} else {
		headers = []string{
			"SL No.",
			"Title",
			"Description",
			"Created By",
			"Created At",
			"Total Questions",
		}
	}

	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, questionSet := range data {
		var title, description, createdBy string

		// Assign values based on language preference
		if userSetting.Language == "ar" {
			title = questionSet.TitleAr
			description = utility.GetStringValue(questionSet.Description)
			createdBy = questionSet.CreatedUserFirstNameAr + " " + questionSet.CreatedUserLastNameAr
		} else {
			title = questionSet.TitleEn
			description = utility.GetStringValue(questionSet.Description)
			createdBy = questionSet.CreatedUserFirstNameEn + " " + questionSet.CreatedUserLastNameEn
		}

		rows[i] = []interface{}{
			i + 1,
			title,
			description,
			createdBy,
			questionSet.CreatedAt,
			questionSet.TotalQuestions,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (svc *QuestionSetService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := svc.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := svc.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	go func() {
		if err := svc.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.QuestionSet{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		); err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}
