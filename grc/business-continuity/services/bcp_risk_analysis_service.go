package services

import (
	"errors"
	"jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpRiskAnalysisServiceInterface interface {
	GetAllRisksByBcpId(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error)
	AddRiskToBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm, currentUserID uint) error
	DeleteRiskFromBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm) error
	GetAllProcessRisks(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error)
}

type BcpRiskAnalysisService struct {
	repo           repositories.BcpRiskAnalysisRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	bcpSvc         BusinessContinuityPlanServiceInterface
}

func NewBcpRiskAnalysisService(repo repositories.BcpRiskAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	bcpSvc BusinessContinuityPlanServiceInterface,
) *BcpRiskAnalysisService {
	return &BcpRiskAnalysisService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		bcpSvc:         bcpSvc,
	}
}

func (bras *BcpRiskAnalysisService) AddRiskToBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm, currentUserID uint) error {
	// Check if BCP exists
	_, err := bras.bcpSvc.GetFullDetailsByID(bcpId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utility.NewResourceNotFoundError(err, bcpId)
		}
		return nil
	}

	// Prepare bulk insert data
	bcpRisks := make([]models.BcpRisk, len(bcpRiskForm.RiskRegisterIDs))
	for i, riskID := range bcpRiskForm.RiskRegisterIDs {
		bcpRisks[i] = models.BcpRisk{
			BusinessContinuityPlanID: bcpId,
			RiskRegisterID:           riskID,
			CreatedBYID:              currentUserID,
		}
	}

	// Bulk insert
	return bras.repo.AddRiskToBcp(bcpRisks)
}

func (bras *BcpRiskAnalysisService) GetAllRisksByBcpId(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error) {
	return bras.repo.GetAllRisksByBcpId(bcpId, params)
}

func (bras *BcpRiskAnalysisService) DeleteRiskFromBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm) error {
	return bras.repo.DeleteRiskFromBcp(bcpId, bcpRiskForm)
}

func (bras *BcpRiskAnalysisService) GetAllProcessRisks(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error) {
	return bras.repo.GetAllProcessRisks(bcpId, params)
}
