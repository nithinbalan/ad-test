package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SystemRegisterServiceInterface interface {
	GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error)
	GetByID(id uint) (*models.SystemRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.SystemRegisterDetailDTO, error)
	Create(systemRegisterForm *dtos.SystemRegisterFormDTO, currentUserID uint) (uint, error)
	Update(id uint, systemRegisterForm *dtos.SystemRegisterFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ExportSystemRegister(data []*dtos.SystemRegisterListDTO, currentUserID uint) (*os.File, error)
}

type SystemRegisterService struct {
	repo           repositories.SystemRegisterRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewSystemRegisterService(
	repo repositories.SystemRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) SystemRegisterServiceInterface {
	return &SystemRegisterService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (srs *SystemRegisterService) GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error) {
	return srs.repo.GetAll(params)
}

func (srs *SystemRegisterService) GetByID(id uint) (*models.SystemRegister, error) {
	systemRegister, err := srs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return systemRegister, nil
}

func (srs *SystemRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.SystemRegisterDetailDTO, error) {
	systemRegister, err := srs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var detail dtos.SystemRegisterDetailDTO
	if err := mapstructure.Decode(systemRegister, &detail); err != nil {
		return nil, err
	}

	detail.CreatedAt = systemRegister.CreatedAt

	return &detail, nil
}

func (srs *SystemRegisterService) Create(systemRegisterFormDTO *dtos.SystemRegisterFormDTO, currentUserID uint) (uint, error) {
	tx := srs.repo.BeginTransaction()

	systemRegister := &models.SystemRegister{
		TitleEn:                 systemRegisterFormDTO.TitleEn,
		TitleAr:                 systemRegisterFormDTO.TitleAr,
		ServiceDescription:      systemRegisterFormDTO.ServiceDescription,
		ModulesComponents:       systemRegisterFormDTO.ModulesComponents,
		ThirdPartyDependencies:  systemRegisterFormDTO.ThirdPartyDependencies,
		ServiceImportanceID:     systemRegisterFormDTO.ServiceImportanceID,
		ResponsibleITTeamID:     systemRegisterFormDTO.ResponsibleITTeamID,
		OutsourcingPartner:      systemRegisterFormDTO.OutsourcingPartner,
		OperationLocationID:     systemRegisterFormDTO.OperationLocationID,
		BusinessOwnerID:         systemRegisterFormDTO.BusinessOwnerID,
		Customers:               systemRegisterFormDTO.Customers,
		ImpactTypeID:            systemRegisterFormDTO.ImpactTypeID,
		DisruptionImpactLevelID: systemRegisterFormDTO.DisruptionImpactLevelID,
		RTO:                     systemRegisterFormDTO.RTO,
		RTOUnitID:               systemRegisterFormDTO.RTOUnitID,
		MTPD:                    systemRegisterFormDTO.MTPD,
		MTPDUnitID:              systemRegisterFormDTO.MTPDUnitID,
		RPO:                     systemRegisterFormDTO.RPO,
		RPOUnitID:               systemRegisterFormDTO.RPOUnitID,
		ITProcesses:             systemRegisterFormDTO.ITProcesses,
		ITServiceOwnerID:        systemRegisterFormDTO.ITServiceOwnerID,
		ProcessOperationTeamID:  systemRegisterFormDTO.ProcessOperationTeamID,
		ProcessRecoveryTeamID:   systemRegisterFormDTO.ProcessRecoveryTeamID,
		RecoveryManagerID:       systemRegisterFormDTO.RecoveryManagerID,
		RecoveryManagerPhone:    systemRegisterFormDTO.RecoveryManagerPhone,
		RecoveryManagerEmail:    systemRegisterFormDTO.RecoveryManagerEmail,
		RecoveryLocation:        systemRegisterFormDTO.RecoveryLocation,
		VendorRTO:               systemRegisterFormDTO.VendorRTO,
		VendorRTOUnitID:         systemRegisterFormDTO.VendorRTOUnitID,
		CreatedBYID:             currentUserID,
		CreatedAt:               time.Now(),
	}

	systemRegisterID, err := srs.repo.Create(systemRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if systemRegisterFormDTO.SLARegisterID != nil {
		sla := &models.SystemRegisterSLA{
			SystemRegisterID: systemRegisterID,
			SLARegisterID:    *systemRegisterFormDTO.SLARegisterID,
			CreatedBYID:      currentUserID,
			CreatedAt:        time.Now(),
		}
		if err := srs.repo.CreateSLA(sla, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		detail, err := srs.GetFullDetailsByID(systemRegisterID)
		if err != nil {
			utility.LogError("Failed to fetch system register", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := srs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SystemRegister{}.TableName(),
			systemRegisterID,
			detail,
			currentUserID,
			"asset",
			"system_register",
			"",
			nil,
		); err != nil {
			utility.LogError("Failed to log system register activity", map[string]interface{}{
				"action": models.CreateAction, "recordID": systemRegisterID, "error": err.Error(),
			})
		}
	}()

	return systemRegisterID, nil
}

func (srs *SystemRegisterService) Update(id uint, systemRegisterFormDTO *dtos.SystemRegisterFormDTO, currentUserID uint) error {
	tx := srs.repo.BeginTransaction()

	existingSystemRegister, err := srs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldDetail, err := srs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingSystemRegister.TitleEn = systemRegisterFormDTO.TitleEn
	existingSystemRegister.TitleAr = systemRegisterFormDTO.TitleAr
	existingSystemRegister.ServiceDescription = systemRegisterFormDTO.ServiceDescription
	existingSystemRegister.ModulesComponents = systemRegisterFormDTO.ModulesComponents
	existingSystemRegister.ThirdPartyDependencies = systemRegisterFormDTO.ThirdPartyDependencies
	existingSystemRegister.ServiceImportanceID = systemRegisterFormDTO.ServiceImportanceID
	existingSystemRegister.ResponsibleITTeamID = systemRegisterFormDTO.ResponsibleITTeamID
	existingSystemRegister.OutsourcingPartner = systemRegisterFormDTO.OutsourcingPartner
	existingSystemRegister.OperationLocationID = systemRegisterFormDTO.OperationLocationID
	existingSystemRegister.BusinessOwnerID = systemRegisterFormDTO.BusinessOwnerID
	existingSystemRegister.Customers = systemRegisterFormDTO.Customers
	existingSystemRegister.ImpactTypeID = systemRegisterFormDTO.ImpactTypeID
	existingSystemRegister.DisruptionImpactLevelID = systemRegisterFormDTO.DisruptionImpactLevelID
	existingSystemRegister.RTO = systemRegisterFormDTO.RTO
	existingSystemRegister.RTOUnitID = systemRegisterFormDTO.RTOUnitID
	existingSystemRegister.MTPD = systemRegisterFormDTO.MTPD
	existingSystemRegister.MTPDUnitID = systemRegisterFormDTO.MTPDUnitID
	existingSystemRegister.RPO = systemRegisterFormDTO.RPO
	existingSystemRegister.RPOUnitID = systemRegisterFormDTO.RPOUnitID
	existingSystemRegister.ITProcesses = systemRegisterFormDTO.ITProcesses
	existingSystemRegister.ITServiceOwnerID = systemRegisterFormDTO.ITServiceOwnerID
	existingSystemRegister.ProcessOperationTeamID = systemRegisterFormDTO.ProcessOperationTeamID
	existingSystemRegister.ProcessRecoveryTeamID = systemRegisterFormDTO.ProcessRecoveryTeamID
	existingSystemRegister.RecoveryManagerID = systemRegisterFormDTO.RecoveryManagerID
	existingSystemRegister.RecoveryManagerPhone = systemRegisterFormDTO.RecoveryManagerPhone
	existingSystemRegister.RecoveryManagerEmail = systemRegisterFormDTO.RecoveryManagerEmail
	existingSystemRegister.RecoveryLocation = systemRegisterFormDTO.RecoveryLocation
	existingSystemRegister.VendorRTO = systemRegisterFormDTO.VendorRTO
	existingSystemRegister.VendorRTOUnitID = systemRegisterFormDTO.VendorRTOUnitID
	existingSystemRegister.UpdatedBYID = &currentUserID
	now := time.Now()
	existingSystemRegister.UpdatedAt = &now

	if err := srs.repo.Update(existingSystemRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := srs.repo.DeleteSLAs(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if systemRegisterFormDTO.SLARegisterID != nil {
		sla := &models.SystemRegisterSLA{
			SystemRegisterID: id,
			SLARegisterID:    *systemRegisterFormDTO.SLARegisterID,
			CreatedBYID:      currentUserID,
			CreatedAt:        time.Now(),
		}
		if err := srs.repo.CreateSLA(sla, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		newDetail, err := srs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		oldData := map[string]interface{}{
			"title_en":                 oldDetail.TitleEn,
			"title_ar":                 oldDetail.TitleAr,
			"service_description":      oldDetail.ServiceDescription,
			"third_party_dependencies": oldDetail.ThirdPartyDependencies,
			"outsourcing_partner":      oldDetail.OutsourcingPartner,
			"customers":                oldDetail.Customers,
			"recovery_manager_phone":   oldDetail.RecoveryManagerPhone,
			"recovery_manager_email":   oldDetail.RecoveryManagerEmail,
		}

		newData := map[string]interface{}{
			"title_en":                 newDetail.TitleEn,
			"title_ar":                 newDetail.TitleAr,
			"service_description":      newDetail.ServiceDescription,
			"third_party_dependencies": newDetail.ThirdPartyDependencies,
			"outsourcing_partner":      newDetail.OutsourcingPartner,
			"customers":                newDetail.Customers,
			"recovery_manager_phone":   newDetail.RecoveryManagerPhone,
			"recovery_manager_email":   newDetail.RecoveryManagerEmail,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)
		if len(changedData) > 0 {
			if err := srs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SystemRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset",
				"system_register",
				"",
				nil,
			); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action": models.UpdateAction, "recordID": id, "error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (srs *SystemRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedSystemRegisters []map[string]interface{}

	for _, id := range ids {
		tx := srs.repo.BeginTransaction()

		systemRegister, err := srs.repo.GetByID(id)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := srs.repo.DeleteSLAs(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := srs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := tx.Commit(); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedSystemRegisters = append(deletedSystemRegisters, map[string]interface{}{
			"id":       id,
			"title_en": systemRegister.TitleEn,
			"title_ar": systemRegister.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			if err := srs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SystemRegister{}.TableName(),
				successfulIds,
				deletedSystemRegisters,
				currentUserID,
				nil,
			); err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (srs *SystemRegisterService) ExportSystemRegister(data []*dtos.SystemRegisterListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := srs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"الرقم التسلسلي",
			"العنوان",
			"وصف الخدمة",
			"أهمية الخدمة",
			"فريق تكنولوجيا المعلومات المسؤول",
			"مالك العمل",
			"فريق عمليات العملية",
			"فريق استعادة العملية",
			"هاتف مدير الاستعادة",
			"بريد مدير الاستعادة",
		}
	} else {
		headers = []string{
			"SL No.",
			"Title",
			"Service Description",
			"Service Importance",
			"Responsible IT Team",
			"Business Owner",
			"Process Operation Team",
			"Process Recovery Team",
			"Recovery Manager Phone",
			"Recovery Manager Email",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, sr := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				sr.TitleAr,
				utility.GetStringValue(sr.ServiceDescription),
				utility.GetStringValue(sr.ServiceImportanceTitleAr),
				utility.GetStringValue(sr.ResponsibleITTeamTitleAr),
				utility.GetStringValue(sr.BusinessOwnerTitleAr),
				utility.GetStringValue(sr.ProcessOperationTeamTitleAr),
				utility.GetStringValue(sr.ProcessRecoveryTeamTitleAr),
				sr.RecoveryManagerPhone,
				sr.RecoveryManagerEmail,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				sr.TitleEn,
				utility.GetStringValue(sr.ServiceDescription),
				utility.GetStringValue(sr.ServiceImportanceTitleEn),
				utility.GetStringValue(sr.ResponsibleITTeamTitleEn),
				utility.GetStringValue(sr.BusinessOwnerTitleEn),
				utility.GetStringValue(sr.ProcessOperationTeamTitleEn),
				utility.GetStringValue(sr.ProcessRecoveryTeamTitleEn),
				sr.RecoveryManagerPhone,
				sr.RecoveryManagerEmail,
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
