package services

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	"jethur/utility"
)

type ComplianceCategoryServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]models.ComplianceCategory, error)
	GetByID(id uint) (models.ComplianceCategory, error)
	Create(category models.ComplianceCategory) (models.ComplianceCategory, error)
	Update(id uint, category models.ComplianceCategory) (models.ComplianceCategory, error)
	Delete(ids []uint) ([]uint, error)
	Activate(ids []uint) error
	Deactivate(ids []uint) error
}

// Service and implementations
type ComplianceCategoryService struct {
	repo repositories.ComplianceCategoryRepositoryInterface
}

func NewComplianceCategoryService(repo repositories.ComplianceCategoryRepositoryInterface) *ComplianceCategoryService {
	return &ComplianceCategoryService{
		repo: repo,
	}
}

func (ccs *ComplianceCategoryService) GetAll(sort string, sortBy string, keyword string, status string) ([]models.ComplianceCategory, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status)
}

func (ccs *ComplianceCategoryService) GetByID(id uint) (models.ComplianceCategory, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *ComplianceCategoryService) Create(category models.ComplianceCategory) (models.ComplianceCategory, error) {
	if exists, err := ccs.repo.ExistsByTitle(category.TitleEn, category.TitleAr); err != nil {
		return models.ComplianceCategory{}, err
	} else if exists {
		return models.ComplianceCategory{}, utility.NewValidationError(errors.New("the title must be unique"), "unique", "title")
	}

	return ccs.repo.Create(category)
}

func (ccs *ComplianceCategoryService) Update(id uint, updatedCategory models.ComplianceCategory) (models.ComplianceCategory, error) {
	existingCategory, err := ccs.repo.GetByID(id)
	if err != nil {
		return models.ComplianceCategory{}, err
	}

	if updatedCategory.TitleEn != existingCategory.TitleEn || updatedCategory.TitleAr != existingCategory.TitleAr {
		if exists, err := ccs.repo.ExistsByTitleExcludingID(id, updatedCategory.TitleEn, updatedCategory.TitleAr); err != nil {
			return models.ComplianceCategory{}, err
		} else if exists {
			return models.ComplianceCategory{}, utility.NewValidationError(errors.New("the title must be unique"), "unique", "title")
		}
	}

	return ccs.repo.Update(id, updatedCategory)
}

func (ccs *ComplianceCategoryService) Delete(ids []uint) ([]uint, error) {
	var failedDeletes []uint
	for _, id := range ids {
		err := ccs.repo.Delete(id)
		if err != nil {
			failedDeletes = append(failedDeletes, id)
		}
	}

	return failedDeletes, nil
}

func (service *ComplianceCategoryService) Activate(ids []uint) error {
	return service.repo.Activate(ids)
}

func (service *ComplianceCategoryService) Deactivate(ids []uint) error {
	return service.repo.Deactivate(ids)
}
