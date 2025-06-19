package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"

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
}

type ThirdPartyQuestionLibraryService struct {
	repo                          tprmRepo.ThirdPartyQuestionLibraryRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	thirdPartyQuestionCategorySvc ThirdPartyQuestionCategoryServiceInterface
	thirdPartyQuestionTypeSvc     ThirdPartyQuestionTypeServiceInterface
}

func NewThirdPartyQuestionLibraryService(
	repo tprmRepo.ThirdPartyQuestionLibraryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	thirdPartyQuestionCategorySvc ThirdPartyQuestionCategoryServiceInterface,
	thirdPartyQuestionTypeSvc ThirdPartyQuestionTypeServiceInterface,
) *ThirdPartyQuestionLibraryService {
	return &ThirdPartyQuestionLibraryService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		thirdPartyQuestionCategorySvc: thirdPartyQuestionCategorySvc,
		thirdPartyQuestionTypeSvc:     thirdPartyQuestionTypeSvc,
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
	tx := tpqls.repo.BeginTransaction()

	if err := tpqls.validateForm(form); err != nil {
		return 0, err
	}
	refNo, err := tpqls.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	entity := &models.ThirdPartyQuestionLibrary{
		ReferenceNumber:              refNo,
		QuestionTitleEn:              form.QuestionTitleEn,
		QuestionTitleAr:              form.QuestionTitleAr,
		ThirdPartyQuestionCategoryID: form.ThirdPartyQuestionCategoryID,
		Mandatory:                    form.Mandatory,
		ThirdPartyQuestionTypeID:     form.ThirdPartyQuestionTypeID,
		CreatedBYID:                  &currentUserID,
	}

	id, err := tpqls.repo.Create(entity, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
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
	tx.Commit()
	return id, nil
}

func (tpqls *ThirdPartyQuestionLibraryService) Update(id uint, currentUserID uint, form *dtos.ThirdPartyQuestionLibraryFormDTO) error {
	tx := tpqls.repo.BeginTransaction()

	existing, err := tpqls.GetByID(id)
	if err != nil {
		return err
	}

	if err := tpqls.validateForm(form); err != nil {
		return err
	}
	// Update base fields
	existing.QuestionTitleEn = form.QuestionTitleEn
	existing.QuestionTitleAr = form.QuestionTitleAr
	existing.ThirdPartyQuestionCategoryID = form.ThirdPartyQuestionCategoryID
	existing.ThirdPartyQuestionTypeID = form.ThirdPartyQuestionTypeID
	existing.Mandatory = form.Mandatory
	existing.UpdatedBYID = &currentUserID

	// Save main question
	err = tpqls.repo.Update(existing, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
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
	tx.Commit()
	return nil
}

func (tpqls *ThirdPartyQuestionLibraryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint

	for _, id := range ids {
		tx := tpqls.repo.BeginTransaction()

		_, err := tpqls.GetByID(id)
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
	if _, err := tpqls.thirdPartyQuestionTypeSvc.GetByID(form.ThirdPartyQuestionTypeID); err != nil {
		return utility.NewValidationError(errors.New("third party question type not found"), "not_exist", "third_party_question_type")
	}

	return nil
}
