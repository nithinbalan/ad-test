package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type QuestionSetServiceInterface interface {
	GetAll(params *dtos.QuestionSetListQueryParamsDTO) (*dtos.QuestionSetListResponseDTO, error)
	GetByID(id uint) (*models.QuestionSet, error)
	GetFullDetailsByID(id uint) (*dtos.QuestionSetDetailDTO, error)
	Create(currentUserID uint, form *dtos.QuestionSetFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *dtos.QuestionSetFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	AddQuestionsToSet(questionSetID, currentUserID uint, dto *dtos.QuestionSetAddQuestionsDTO) error
	DeleteQuestionFromSet(questionID uint) error
}

type QuestionSetService struct {
	repo                         tprmRepo.QuestionSetRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	thirdPartyQuestionLibrarySvc ThirdPartyQuestionLibraryServiceInterface
}

func NewQuestionSetService(
	repo tprmRepo.QuestionSetRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	thirdPartyQuestionLibrarySvc ThirdPartyQuestionLibraryServiceInterface,
) *QuestionSetService {
	return &QuestionSetService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		thirdPartyQuestionLibrarySvc: thirdPartyQuestionLibrarySvc,
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

	// Temp map for grouping
	categoryMap := make(map[string]*dtos.QuestionSetCategoryGroupDTO)

	for _, q := range entity.Questions {
		var catEn, catAr string
		var catID uint

		questionDTO := dtos.QuestionSetQuestionDetailDTO{
			ID:                          q.ID,
			QuestionSetID:               q.QuestionSetID,
			ThirdPartyQuestionLibraryID: q.ThirdPartyQuestionLibraryID,
			QuestionTitleEn:             q.QuestionTitleEn,
			QuestionTitleAr:             q.QuestionTitleAr,
		}

		if q.QuestionLibrary != nil {
			questionDTO.ReferenceNumber = q.QuestionLibrary.ReferenceNumber
			questionDTO.Mandatory = q.QuestionLibrary.Mandatory
			questionDTO.ThirdPartyQuestionCategoryID = q.QuestionLibrary.ThirdPartyQuestionCategoryID
			questionDTO.ThirdPartyQuestionTypeID = q.QuestionLibrary.ThirdPartyQuestionTypeID

			if q.QuestionLibrary.QuestionCategory != nil {
				catEn = q.QuestionLibrary.QuestionCategory.TitleEn
				catAr = q.QuestionLibrary.QuestionCategory.TitleAr
				catID = q.QuestionLibrary.QuestionCategory.ID
				questionDTO.ThirdPartyQuestionCategoryID = q.QuestionLibrary.QuestionCategory.ID
				questionDTO.QuestionCategoryTitleEn = catEn
				questionDTO.QuestionCategoryTitleAr = catAr
			}

			if q.QuestionLibrary.QuestionType != nil {
				questionDTO.QuestionTypeTitleEn = q.QuestionLibrary.QuestionType.TitleEn
				questionDTO.QuestionTypeTitleAr = q.QuestionLibrary.QuestionType.TitleAr
				questionDTO.QuestionTypeType = q.QuestionLibrary.QuestionType.Type
			}

			for _, opt := range q.QuestionLibrary.AnswerOptions {
				optionDTO := dtos.ThirdPartyQuestionLibraryAnswerOptionDetailDTO{
					ID:                  opt.ID,
					AnswerOptionTitleEn: opt.AnswerOptionTitleEn,
					AnswerOptionTitleAr: opt.AnswerOptionTitleAr,
					RiskLevelID:         opt.RiskLevelID,
					RangeFrom:           opt.RangeFrom,
					RangeTo:             opt.RangeTo,
				}
				if opt.RiskLevel != nil {
					optionDTO.RiskLevel.ID = opt.RiskLevel.ID
					optionDTO.RiskLevel.TitleEn = opt.RiskLevel.TitleEn
					optionDTO.RiskLevel.TitleAr = opt.RiskLevel.TitleAr
					optionDTO.RiskLevel.ColorCode = opt.RiskLevel.ColorCode
					optionDTO.RiskLevel.Type = opt.RiskLevel.Type
				}
				questionDTO.AnswerOptions = append(questionDTO.AnswerOptions, optionDTO)
			}
		}

		groupKey := catEn + "::" + catAr
		if _, exists := categoryMap[groupKey]; !exists {
			categoryMap[groupKey] = &dtos.QuestionSetCategoryGroupDTO{
				CategoryID:      catID,
				CategoryTitleEn: catEn,
				CategoryTitleAr: catAr,
				Questions:       []dtos.QuestionSetQuestionDetailDTO{},
			}
		}
		categoryMap[groupKey].Questions = append(categoryMap[groupKey].Questions, questionDTO)
	}

	// Convert map to slice
	dto.Categories = []dtos.QuestionSetCategoryGroupDTO{}
	for _, v := range categoryMap {
		dto.Categories = append(dto.Categories, *v)
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

	var questions []*models.QuestionSetQuestion

	for _, libID := range dto.QuestionLibraryIDs {
		lib, err := svc.thirdPartyQuestionLibrarySvc.GetByID(libID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "not_exist", fmt.Sprintf("question_library_id_%d", libID))
		}

		questions = append(questions, &models.QuestionSetQuestion{
			QuestionSetID:               questionSetID,
			ThirdPartyQuestionLibraryID: libID,
			QuestionTitleEn:             lib.QuestionTitleEn,
			QuestionTitleAr:             lib.QuestionTitleAr,
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

func (svc *QuestionSetService) DeleteQuestionFromSet(questionID uint) error {
	tx := svc.repo.BeginTransaction()

	if err := svc.repo.DeleteByID(questionID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
