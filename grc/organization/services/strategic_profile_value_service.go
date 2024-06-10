package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategicProfileValueServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.StrategicProfileValue, error)
	GetByID(id uint) (*models.StrategicProfileValue, error)
	GetFullDetailsByID(id uint) (*dtos.StrategicProfileValueDetailsDTO, error)
	Create(currentUserID uint, strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO) (uint, error)
	Update(id uint, currentUserID uint, strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategicProfileValueService struct {
	repo                 repositories.StrategicProfileValueRepositoryInterface
	activityLogSvc       services.ActivityLogServiceInterface
	strategicProfileRepo repositories.StrategicProfileRepositoryInterface
}

func NewStrategicProfileValueService(repo repositories.StrategicProfileValueRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface, strategicProfileRepo repositories.StrategicProfileRepositoryInterface) StrategicProfileValueServiceInterface {
	return &StrategicProfileValueService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		strategicProfileRepo: strategicProfileRepo,
	}
}

func (spvs *StrategicProfileValueService) GetAll(sort string, sortBy string, keyword string) ([]*models.StrategicProfileValue, error) {
	return spvs.repo.GetAll(sort, sortBy, keyword)
}

func (spvs StrategicProfileValueService) GetByID(id uint) (*models.StrategicProfileValue, error) {
	strategicProfileValue, err := spvs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return strategicProfileValue, nil

}

func (spvs StrategicProfileValueService) GetFullDetailsByID(id uint) (*dtos.StrategicProfileValueDetailsDTO, error) {
	// Initialize your DTO
	var strategicProfileValueDTO dtos.StrategicProfileValueDetailsDTO

	strategicProfileValue, err := spvs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategicProfileValue, &strategicProfileValueDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &strategicProfileValueDTO, nil
}

func (spvs *StrategicProfileValueService) Create(currentUserID uint, strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO) (uint, error) {
	if exists, err := spvs.repo.ExistsByTitle(strategicProfileValueFormDTO.TitleEn, strategicProfileValueFormDTO.TitleAr); err != nil {
		return 0, err
	} else if exists {
		return 0, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}
	existingStrategicProfile, err := spvs.strategicProfileRepo.Get()

	if err != nil {
		return 0, utility.NewResourceNotFoundError(err, nil)
	}

	strategicProfileValue := &models.StrategicProfileValue{
		StrategicProfileID: existingStrategicProfile.ID,
		TitleEn:            strategicProfileValueFormDTO.TitleEn,
		TitleAr:            strategicProfileValueFormDTO.TitleAr,
		DescriptionEn:      strategicProfileValueFormDTO.DescriptionEn,
		DescriptionAr:      strategicProfileValueFormDTO.DescriptionAr,
	}

	createdStrategicProfileValueID, err := spvs.repo.Create(strategicProfileValue)
	if err != nil {
		return 0, err
	}

	go func() {
		err := spvs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategicProfileValue{}.TableName(),
			createdStrategicProfileValueID,
			strategicProfileValueFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdStrategicProfileValueID, nil
}

func (spvs *StrategicProfileValueService) Update(id uint, currentUserID uint, strategicProfileValueFormDTO *dtos.StrategicProfileValueFormDTO) error {
	existingStrategicProfileValue, err := spvs.GetByID(id)

	if err != nil {
		return err
	}

	if strategicProfileValueFormDTO.TitleEn != existingStrategicProfileValue.TitleEn || strategicProfileValueFormDTO.TitleAr != existingStrategicProfileValue.TitleAr {
		if exists, err := spvs.repo.ExistsByTitleExcludingID(id, strategicProfileValueFormDTO.TitleEn, strategicProfileValueFormDTO.TitleAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	existingStrategicProfileValue.TitleEn = strategicProfileValueFormDTO.TitleEn
	existingStrategicProfileValue.TitleAr = strategicProfileValueFormDTO.TitleAr
	existingStrategicProfileValue.DescriptionEn = strategicProfileValueFormDTO.DescriptionEn
	existingStrategicProfileValue.DescriptionAr = strategicProfileValueFormDTO.DescriptionAr

	err = spvs.repo.Update(existingStrategicProfileValue)

	if err != nil {
		return err
	}

	//@todo - need to capture activity log
	return nil
}

func (spvs *StrategicProfileValueService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := spvs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = spvs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := spvs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategicProfileValue{}.TableName(),
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
