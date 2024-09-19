package services

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditInsightsAuditServiceInterface interface {
	GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error)
	GetByID(id uint) (*models.Audit, error)
}

type AuditInsightsAuditService struct {
	repo           repositories.AuditInsightsAuditRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewAuditInsightsAuditService(
	repo repositories.AuditInsightsAuditRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AuditInsightsAuditServiceInterface {
	return &AuditInsightsAuditService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (aias *AuditInsightsAuditService) GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditListResponseDTO, error) {
	orgUnitID, err := aias.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aias.repo.GetAll(params, orgUnitID)
}

func (aias *AuditInsightsAuditService) GetByID(id uint) (*models.Audit, error) {
	audit, err := aias.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return audit, nil
}

func (aias *AuditInsightsAuditService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error) {
	var auditDTO auditDTOs.AuditDetailDTO

	audit, err := aias.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(audit, &auditDTO)
	if err != nil {
		return nil, err
	}

	auditDTO.StartDateString = audit.StartDate.Format("2006-01-02")
	auditDTO.EndDateString = audit.EndDate.Format("2006-01-02")

	auditDTO.CreatedAt = audit.CreatedAt
	auditDTO.UpdatedAt = audit.UpdatedAt

	return &auditDTO, nil
}
