package services

import (
	"errors"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessExternalDependencyServiceInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessExternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessExternalDependencyListResponseDTO, error)
	GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessExternalDependencyDetailDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessExternalDependency, error)
	Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessExternalDependencyFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessExternalDependencyFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessExternalDependencyService struct {
	repo           repositories.ProcessExternalDependencyRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	processSvc     ProcessServiceInterface
}

func NewProcessExternalDependencyService(
	repo repositories.ProcessExternalDependencyRepositoryInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
) ProcessExternalDependencyServiceInterface {
	return &ProcessExternalDependencyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		processSvc:     processSvc,
	}
}

func (peds *ProcessExternalDependencyService) GetAll(processID uint, params *orgDTOs.ProcessExternalDependencyListQueryParamsDTO) (*orgDTOs.ProcessExternalDependencyListResponseDTO, error) {
	if _, err := peds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return peds.repo.GetAll(processID, params)
}

func (peds *ProcessExternalDependencyService) GetByID(processID uint, id uint) (*models.ProcessExternalDependency, error) {
	if _, err := peds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	item, err := peds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return item, nil
}

func (peds *ProcessExternalDependencyService) GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessExternalDependencyDetailDTO, error) {
	var detailDTO orgDTOs.ProcessExternalDependencyDetailDTO

	item, err := peds.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(item, &detailDTO)
	if err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = item.CreatedAt
	detailDTO.UpdatedAt = item.UpdatedAt

	return &detailDTO, nil
}

func (peds *ProcessExternalDependencyService) getVendorRTOFromSystemRegister(slaRegisterID *uint) (*int, *uint, error) {
	if slaRegisterID == nil {
		return nil, nil, nil
	}

	return peds.repo.GetVendorRTOBySLARegisterID(*slaRegisterID)
}

func (peds *ProcessExternalDependencyService) Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessExternalDependencyFormDTO) (uint, error) {
	if _, err := peds.processSvc.GetByID(processID); err != nil {
		return 0, err
	}

	vendorRTO, vendorRTOUnitID, err := peds.getVendorRTOFromSystemRegister(dto.SLARegisterID)
	if err != nil {
		return 0, err
	}

	item := &models.ProcessExternalDependency{
		ProcessID:                processID,
		EntityName:               dto.EntityName,
		EntityTypeID:             dto.EntityTypeID,
		NeededSupportDescription: dto.NeededSupportDescription,
		DependencyTierID:         dto.DependencyTierID,
		SLARegisterID:            dto.SLARegisterID,
		VendorRTO:                vendorRTO,
		VendorRTOUnitID:          vendorRTOUnitID,
		InternalStakeholder:      dto.InternalStakeholder,
		ContactNumber:            dto.ContactNumber,
		ContactEmail:             dto.ContactEmail,
		CreatedBYID:              &currentUserID,
	}

	id, err := peds.repo.Create(item)
	if err != nil {
		return 0, err
	}

	go func() {
		created, err := peds.repo.GetByID(id)
		if err != nil {
			return
		}
		_ = peds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessExternalDependency{}.TableName(),
			id,
			created,
			currentUserID,
			"organization", "process", "process_external_dependency", nil,
		)
	}()

	return id, nil
}

func (peds *ProcessExternalDependencyService) Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessExternalDependencyFormDTO) error {
	if _, err := peds.processSvc.GetByID(processID); err != nil {
		return err
	}

	existing, err := peds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	oldData := *existing

	vendorRTO, vendorRTOUnitID, err := peds.getVendorRTOFromSystemRegister(dto.SLARegisterID)
	if err != nil {
		return err
	}

	existing.EntityName = dto.EntityName
	existing.EntityTypeID = dto.EntityTypeID
	existing.NeededSupportDescription = dto.NeededSupportDescription
	existing.DependencyTierID = dto.DependencyTierID
	existing.SLARegisterID = dto.SLARegisterID
	existing.VendorRTO = vendorRTO
	existing.VendorRTOUnitID = vendorRTOUnitID
	existing.InternalStakeholder = dto.InternalStakeholder
	existing.ContactNumber = dto.ContactNumber
	existing.ContactEmail = dto.ContactEmail
	existing.UpdatedBYID = &currentUserID

	if err := peds.repo.Update(existing); err != nil {
		return err
	}

	go func() {
		oldFields := map[string]interface{}{
			"entity_name":                oldData.EntityName,
			"entity_type_id":             oldData.EntityTypeID,
			"needed_support_description": oldData.NeededSupportDescription,
			"dependency_tier_id":         oldData.DependencyTierID,
			"sla_register_id":            oldData.SLARegisterID,
			"vendor_rto":                 oldData.VendorRTO,
			"vendor_rto_unit_id":         oldData.VendorRTOUnitID,
			"internal_stakeholder":       oldData.InternalStakeholder,
			"contact_number":             oldData.ContactNumber,
			"contact_email":              oldData.ContactEmail,
		}
		newFields := map[string]interface{}{
			"entity_name":                existing.EntityName,
			"entity_type_id":             existing.EntityTypeID,
			"needed_support_description": existing.NeededSupportDescription,
			"dependency_tier_id":         existing.DependencyTierID,
			"sla_register_id":            existing.SLARegisterID,
			"vendor_rto":                 existing.VendorRTO,
			"vendor_rto_unit_id":         existing.VendorRTOUnitID,
			"internal_stakeholder":       existing.InternalStakeholder,
			"contact_number":             existing.ContactNumber,
			"contact_email":              existing.ContactEmail,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)
		if len(changedData) > 0 {
			_ = peds.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessExternalDependency{}.TableName(),
				id,
				changedData,
				currentUserID,
				"organization", "process", "process_external_dependency", nil,
			)
		}
	}()

	return nil
}

func (peds *ProcessExternalDependencyService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := peds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if _, err := peds.repo.GetByID(id); err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := peds.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go peds.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessExternalDependency{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}
