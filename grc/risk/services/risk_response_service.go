package services

import (
	"errors"
	"time"

	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskResponseServiceInterface interface {
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskResponseDetailDTO, error)
	GetFullDetailsByRiskID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskResponseDetailDTO, error)
	GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	CreateIfNotExists(riskRegisterID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error
	DeleteByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error
}

type RiskResponseService struct {
	repo riskRepositories.RiskResponseRepositoryInterface
}

func NewRiskResponseService(repo riskRepositories.RiskResponseRepositoryInterface) RiskResponseServiceInterface {
	return &RiskResponseService{
		repo: repo,
	}
}

func (svc *RiskResponseService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	return svc.repo.GetByID(id, tx...)
}

func (svc *RiskResponseService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskResponseDetailDTO, error) {
	var dto riskDTOs.RiskResponseDetailDTO

	riskResponse, err := svc.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskResponse, &dto)
	if err != nil {
		return nil, err
	}

	dto.SubmittedAt = riskResponse.SubmittedAt
	dto.ApprovedAt = riskResponse.ApprovedAt
	dto.CreatedAt = riskResponse.CreatedAt
	dto.UpdatedAt = riskResponse.UpdatedAt

	// Convert ActionPlans dates from time.Time to string
	if dto.RiskRegister != nil && len(dto.RiskRegister.ActionPlans) > 0 {
		for i := range dto.RiskRegister.ActionPlans {
			if riskResponse.RiskRegister != nil && len(riskResponse.RiskRegister.ActionPlans) > i {
				modelActionPlan := riskResponse.RiskRegister.ActionPlans[i]
				if modelActionPlan.StartDate != nil {
					dto.RiskRegister.ActionPlans[i].StartDate = modelActionPlan.StartDate.Format("2006-01-02")
				}
				if modelActionPlan.TargetDate != nil {
					dto.RiskRegister.ActionPlans[i].TargetDate = modelActionPlan.TargetDate.Format("2006-01-02")
				}
			}
		}
	}

	return &dto, nil
}

func (svc *RiskResponseService) GetFullDetailsByRiskID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskResponseDetailDTO, error) {
	var dto riskDTOs.RiskResponseDetailDTO

	riskResponse, err := svc.repo.GetFullDetailsByRiskRegisterID(riskRegisterID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskRegisterID)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskResponse, &dto)
	if err != nil {
		return nil, err
	}

	dto.SubmittedAt = riskResponse.SubmittedAt
	dto.ApprovedAt = riskResponse.ApprovedAt
	dto.CreatedAt = riskResponse.CreatedAt
	dto.UpdatedAt = riskResponse.UpdatedAt

	// Convert ActionPlans dates from time.Time to string
	if dto.RiskRegister != nil && len(dto.RiskRegister.ActionPlans) > 0 {
		for i := range dto.RiskRegister.ActionPlans {
			if riskResponse.RiskRegister != nil && len(riskResponse.RiskRegister.ActionPlans) > i {
				modelActionPlan := riskResponse.RiskRegister.ActionPlans[i]
				if modelActionPlan.StartDate != nil {
					dto.RiskRegister.ActionPlans[i].StartDate = modelActionPlan.StartDate.Format("2006-01-02")
				}
				if modelActionPlan.TargetDate != nil {
					dto.RiskRegister.ActionPlans[i].TargetDate = modelActionPlan.TargetDate.Format("2006-01-02")
				}
			}
		}
	}

	return &dto, nil
}

func (svc *RiskResponseService) GetByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	return svc.repo.GetByRiskRegisterID(riskRegisterID, tx...)
}

func (svc *RiskResponseService) CreateIfNotExists(riskRegisterID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	_, err := svc.repo.GetByRiskRegisterID(riskRegisterID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new RiskResponse with only risk_id, other fields null
			currentTime := time.Now()
			newRiskResponse := &models.RiskResponse{
				RiskRegisterID: riskRegisterID,
				CreatedBYID:    &currentUserID,
				CreatedAt:      currentTime,
			}
			return svc.repo.Create(newRiskResponse, tx...)
		}
		return err
	}
	// Risk response already exists, do nothing
	return nil
}

func (svc *RiskResponseService) Update(riskResponse *models.RiskResponse, tx ...utility.DatabaseTransaction) error {
	return svc.repo.Update(riskResponse, tx...)
}

func (svc *RiskResponseService) DeleteByRiskRegisterID(riskRegisterID uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteByRiskRegisterID(riskRegisterID, tx...)
}
