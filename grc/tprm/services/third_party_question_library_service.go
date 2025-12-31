package services

import (
	"errors"
	"fmt"
	"io"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThirdPartyQuestionLibraryServiceInterface interface {
	GetAll(params *dtos.ThirdPartyQuestionLibraryListQueryParamsDTO) (*dtos.ThirdPartyQuestionLibraryListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyQuestionLibrary, error)
	GetFullDetailsByID(id uint) (*dtos.ThirdPartyQuestionLibraryDetailDTO, error)
	Create(currentUserID uint, form *dtos.ThirdPartyQuestionLibraryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.ThirdPartyQuestionLibraryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	ExportThirdPartyQuestionLibraries(data []*dtos.ThirdPartyQuestionLibraryListDTO, currentUserID uint) (*os.File, error)

	GenerateQuestionLibraryTemplate() (*os.File, error)
	ValidateQuestionLibraryTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportQuestionLibraries(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportQuestionLibraries(input []map[string]interface{}, currentUserID uint) ([]uint, error)
}

type ThirdPartyQuestionLibraryService struct {
	repo                          tprmRepo.ThirdPartyQuestionLibraryRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	thirdPartyQuestionCategorySvc ThirdPartyQuestionCategoryServiceInterface
	thirdPartyQuestionTypeSvc     ThirdPartyQuestionTypeServiceInterface
	tprmSettingSvc                TPRMSettingServiceInterface
	riskLevelSvc                  riskServices.RiskLevelServiceInterface
	complianceResponseSvc         TPRMComplianceResponseServiceInterface
	userSettingSvc                organizationServices.UserSettingServiceInterface
}

func NewThirdPartyQuestionLibraryService(
	repo tprmRepo.ThirdPartyQuestionLibraryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	thirdPartyQuestionCategorySvc ThirdPartyQuestionCategoryServiceInterface,
	thirdPartyQuestionTypeSvc ThirdPartyQuestionTypeServiceInterface,
	tprmSettingSvc TPRMSettingServiceInterface,
	riskLevelSvc riskServices.RiskLevelServiceInterface,
	complianceResponseSvc TPRMComplianceResponseServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,

) *ThirdPartyQuestionLibraryService {
	return &ThirdPartyQuestionLibraryService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		thirdPartyQuestionCategorySvc: thirdPartyQuestionCategorySvc,
		thirdPartyQuestionTypeSvc:     thirdPartyQuestionTypeSvc,
		tprmSettingSvc:                tprmSettingSvc,
		riskLevelSvc:                  riskLevelSvc,
		complianceResponseSvc:         complianceResponseSvc,
		userSettingSvc:                userSettingSvc,
	}
}

func (tpqls *ThirdPartyQuestionLibraryService) GetAll(params *dtos.ThirdPartyQuestionLibraryListQueryParamsDTO) (*dtos.ThirdPartyQuestionLibraryListResponseDTO, error) {
	return tpqls.repo.GetAll(params)
}

func (tpqls *ThirdPartyQuestionLibraryService) GetByID(id uint) (*models.ThirdPartyQuestionLibrary, error) {
	entity, err := tpqls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return entity, nil
}

func (tpqls *ThirdPartyQuestionLibraryService) GetFullDetailsByID(id uint) (*dtos.ThirdPartyQuestionLibraryDetailDTO, error) {
	entity, err := tpqls.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var dto dtos.ThirdPartyQuestionLibraryDetailDTO
	err = mapstructure.Decode(entity, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (tpqls *ThirdPartyQuestionLibraryService) Create(currentUserID uint, form *dtos.ThirdPartyQuestionLibraryFormDTO) (uint, error) {
	settings, err := tpqls.tprmSettingSvc.GetTPRMSetting()
	if err != nil {
		return 0, err
	}

	dropDownType, err := tpqls.thirdPartyQuestionTypeSvc.GetByType("dropdown-single")
	if err != nil {
		return 0, err
	}

	tx := tpqls.repo.BeginTransaction()

	if err := tpqls.validateForm(form); err != nil {
		return 0, err
	}
	refNo, err := tpqls.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	questionType := form.ThirdPartyQuestionTypeID
	if settings.IsExport {
		questionType = &dropDownType.ID
	}

	entity := &models.ThirdPartyQuestionLibrary{
		ReferenceNumber:              refNo,
		QuestionTitleEn:              form.QuestionTitleEn,
		QuestionTitleAr:              form.QuestionTitleAr,
		ThirdPartyQuestionCategoryID: form.ThirdPartyQuestionCategoryID,
		Mandatory:                    form.Mandatory,
		ThirdPartyQuestionTypeID:     *questionType,
		CreatedBYID:                  &currentUserID,
	}

	id, err := tpqls.repo.Create(entity, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !settings.IsExport {
		if len(form.AnswerOptions) > 0 {
			var answerOptions []*models.ThirdPartyQuestionLibraryAnswerOption
			for _, opt := range form.AnswerOptions {
				option := &models.ThirdPartyQuestionLibraryAnswerOption{
					ThirdPartyQuestionLibraryID: id,
					AnswerOptionTitleEn:         opt.AnswerOptionTitleEn,
					AnswerOptionTitleAr:         opt.AnswerOptionTitleAr,
					RiskLevelID:                 opt.RiskLevelID,
					RangeFrom:                   opt.RangeFrom,
					RangeTo:                     opt.RangeTo,
					CreatedBYID:                 &currentUserID,
				}
				answerOptions = append(answerOptions, option)
			}

			if len(answerOptions) > 0 {
				err := tpqls.repo.CreateAnswerOptions(answerOptions, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	} else {
		complianceResponses, err := tpqls.complianceResponseSvc.GetAll("", "", "", "active")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if len(complianceResponses) == 0 {
			return 0, fmt.Errorf("no compliance response found")
		}

		var answerOptions []*models.ThirdPartyQuestionLibraryAnswerOption
		for _, opt := range complianceResponses {
			option := &models.ThirdPartyQuestionLibraryAnswerOption{
				ThirdPartyQuestionLibraryID: id,
				AnswerOptionTitleEn:         &opt.TitleEn,
				AnswerOptionTitleAr:         &opt.TitleEn,
				TPRMComplianceResponseID:    &opt.ID,
				CreatedBYID:                 &currentUserID,
			}
			answerOptions = append(answerOptions, option)
		}

		if len(answerOptions) > 0 {
			err := tpqls.repo.CreateAnswerOptions(answerOptions, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	tx.Commit()
	return id, nil
}

func (tpqls *ThirdPartyQuestionLibraryService) Update(id uint, currentUserID uint, form *dtos.ThirdPartyQuestionLibraryFormDTO) error {

	used, err := tpqls.repo.IsQuestionLibraryUsed(id)
	if err != nil {
		return err
	}
	if used {
		return utility.NewForbiddenError("This question library is already used in a question set or assessment, so it cannot be edited.")
	}

	tx := tpqls.repo.BeginTransaction()

	existing, err := tpqls.GetByID(id)
	if err != nil {
		return err
	}

	if err := tpqls.validateForm(form); err != nil {
		return err
	}
	settings, err := tpqls.tprmSettingSvc.GetTPRMSetting()
	if err != nil {
		return err
	}
	// Update base fields
	existing.QuestionTitleEn = form.QuestionTitleEn
	existing.QuestionTitleAr = form.QuestionTitleAr
	existing.ThirdPartyQuestionCategoryID = form.ThirdPartyQuestionCategoryID
	existing.Mandatory = form.Mandatory
	existing.UpdatedBYID = &currentUserID

	if settings.IsExternalVendor {
		existing.ThirdPartyQuestionTypeID = *form.ThirdPartyQuestionTypeID
	}

	// Save main question
	err = tpqls.repo.Update(existing, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if settings.IsExternalVendor {
		if len(form.AnswerOptions) > 0 {
			// Delete old answer options
			err = tpqls.repo.DeleteAnswerOptionsByLibraryID(id, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Prepare and insert new answer options
			var answerOptions []*models.ThirdPartyQuestionLibraryAnswerOption
			for _, opt := range form.AnswerOptions {
				answerOptions = append(answerOptions, &models.ThirdPartyQuestionLibraryAnswerOption{
					ThirdPartyQuestionLibraryID: id,
					AnswerOptionTitleEn:         opt.AnswerOptionTitleEn,
					AnswerOptionTitleAr:         opt.AnswerOptionTitleAr,
					RiskLevelID:                 opt.RiskLevelID,
					RangeFrom:                   opt.RangeFrom,
					RangeTo:                     opt.RangeTo,
					CreatedBYID:                 &currentUserID,
				})
			}

			if len(answerOptions) > 0 {
				err := tpqls.repo.CreateAnswerOptions(answerOptions, tx)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	tx.Commit()
	return nil
}

func (tpqls *ThirdPartyQuestionLibraryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint

	for _, id := range ids {
		used, err := tpqls.repo.IsQuestionLibraryUsed(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		if used {
			failedIDs = append(failedIDs, id)
			continue
		}

		tx := tpqls.repo.BeginTransaction()

		_, err = tpqls.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		// First, delete all related Answer Options
		if err := tpqls.repo.DeleteAnswerOptionsByLibraryID(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}

		// Then, delete the main Question Library
		if err := tpqls.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}

		tx.Commit()
	}

	return failedIDs, nil
}

func (tpqls *ThirdPartyQuestionLibraryService) validateForm(form *dtos.ThirdPartyQuestionLibraryFormDTO) error {
	// Validate Third Party Question Category
	if _, err := tpqls.thirdPartyQuestionCategorySvc.GetByID(form.ThirdPartyQuestionCategoryID); err != nil {
		return utility.NewValidationError(errors.New("third party question category not found"), "not_exist", "third_party_question_category")
	}

	// Validate Third Party Question Type
	// if _, err := tpqls.thirdPartyQuestionTypeSvc.GetByID(form.ThirdPartyQuestionTypeID); err != nil {
	// 	return utility.NewValidationError(errors.New("third party question type not found"), "not_exist", "third_party_question_type")
	// }

	settings, err := tpqls.tprmSettingSvc.GetTPRMSetting()
	if err != nil {
		return err
	}

	if settings.IsExternalVendor {
		if form.ThirdPartyQuestionTypeID == nil {
			return utility.NewValidationError(fmt.Errorf("question type IDs are required"), "required", "question type IDs")
		}
	}

	return nil
}

func (tpqls *ThirdPartyQuestionLibraryService) ExportThirdPartyQuestionLibraries(data []*dtos.ThirdPartyQuestionLibraryListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := tpqls.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string

	// Select headers based on language
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Question (AR)",
			"Category",
			"Mandatory",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Question (EN)",
			"Category",
			"Mandatory",
		}
	}

	rows := make([][]interface{}, len(data))
	var mandatory string
	for i, item := range data {
		if item.Mandatory {
			mandatory = "Yes"
		} else {
			mandatory = "No"
		}
		if userSetting.Language == "ar" {

			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.QuestionTitleAr,
				item.ThirdPartyQuestionCategoryTitleAr,
				mandatory,
			}

		} else {

			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.QuestionTitleEn,
				item.ThirdPartyQuestionCategoryTitleEn,
				mandatory,
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (rrs *ThirdPartyQuestionLibraryService) GenerateQuestionLibraryTemplate() (*os.File, error) {
	templateHeaders := []string{
		"Category", "Question EN", "Question AR", "Mandatory",
	}

	templateData := [][]interface{}{}

	// Fetch Risk Categories
	categories, err := rrs.thirdPartyQuestionCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch risk categories: %v", err)
	}
	categoryData := [][]interface{}{{"Category"}}
	for _, cat := range categories {
		categoryData = append(categoryData, []interface{}{cat.TitleEn})
	}

	mandatoryHeaders := []string{"IsEvidenceMandatory"}
	mandatoryData := [][]interface{}{
		{"Yes"},
		{"No"},
	}

	dropDownConfig := map[string]bool{
		"Category":  true,
		"Mandatory": true,
	}

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":  append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"Category":  categoryData,
		"Mandatory": append([][]interface{}{utility.ConvertToInterfaceSlice(mandatoryHeaders)}, mandatoryData...),
	}, dropDownConfig)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (rrs *ThirdPartyQuestionLibraryService) ValidateQuestionLibraryTemplate(file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Category", "Question EN", "Question AR", "Mandatory",
	}

	snakeHeaders := []string{
		"category", "question_en", "question_ar", "mandatory",
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}
	if !utility.ValidateHeaders(headers, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Risk Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, header := range snakeHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = ""
			}
		}
		input = append(input, record)
	}

	validatedData, err := rrs.validateQuestionLibraries(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (cs *ThirdPartyQuestionLibraryService) ValidateImportQuestionLibraries(input []map[string]interface{}) ([]map[string]interface{}, error) {
	validatedData, err := cs.validateQuestionLibraries(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (rs *ThirdPartyQuestionLibraryService) ImportQuestionLibraries(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := rs.validateQuestionLibraries(input)
	if err != nil {
		return nil, err
	}

	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validatedData {
		isMandatory := false

		rawVal := record["mandatory"]
		switch v := rawVal.(type) {
		case string:
			isMandatory = strings.ToLower(strings.TrimSpace(v)) == "yes"
		case bool:
			isMandatory = v
		}

		dto := &dtos.ThirdPartyQuestionLibraryFormDTO{
			ThirdPartyQuestionCategoryID: record["category_id"].(uint),
			QuestionTitleEn:              record["question_en"].(string),
			QuestionTitleAr:              record["question_ar"].(string),
			Mandatory:                    isMandatory,
		}

		id, err := rs.Create(currentUserID, dto)
		if err != nil {
			return nil, fmt.Errorf("failed to create question library for record: %+v, error: %v", record, err)
		}
		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

func (rs *ThirdPartyQuestionLibraryService) validateQuestionLibraries(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"category", "question_en", "question_ar"}

	// Fetch risk categories
	riskCategories, err := rs.thirdPartyQuestionCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching risk categories: %v", err)
	}
	riskCategoryTitles := make(map[string]uint)
	for _, cat := range riskCategories {
		riskCategoryTitles[strings.TrimSpace(cat.TitleEn)] = cat.ID
	}

	var validatedData []map[string]interface{}
	for _, record := range input {
		rowErrors := make(map[string]interface{})
		riskRecord := make(map[string]interface{})

		// Copy original data
		for key, value := range record {
			riskRecord[key] = value
		}

		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {
			case "category":
				if !ok || value == "" {
					rowErrors[header] = "Question Category is required"
				} else if id, exists := riskCategoryTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid question category: %s", value)
				} else {
					riskRecord["category_id"] = id
				}

			case "question_en":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Question EN is required"
				}
			case "question_ar":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Question AR is required"
				}
			case "mandatory":
				trimmed := strings.ToLower(strings.TrimSpace(value))
				if trimmed == "yes" {
					riskRecord["mandatory"] = true
				} else if trimmed == "no" {
					riskRecord["mandatory"] = false
				} else {
					rowErrors["mandatory"] = "Value must be Yes or No"
				}

			}

		}

		if len(rowErrors) == 0 {
			riskRecord["errors"] = nil
		} else {
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = ""
				}
			}
			riskRecord["errors"] = rowErrors
		}

		validatedData = append(validatedData, riskRecord)
	}

	return validatedData, nil
}
