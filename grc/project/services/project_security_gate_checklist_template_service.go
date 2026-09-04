package services

import (
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
)

type ProjectSecurityGateChecklistTemplateServiceInterface interface {
	GetAll() ([]*dtos.ProjectSecurityGateChecklistTemplateDTO, error)
	GetByID(id uint) (*models.ProjectSecurityGateChecklistTemplate, error)
	GetByStageID(stageID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplate, error)
	GetDetailsByID(id uint) (*dtos.ProjectSecurityGateChecklistTemplateDetailDTO, error)
	GetLatestPublishedVersionByTemplateID(templateID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error)
	CreateVersion(templateID uint, currentUserID uint) (uint, error)
	PublishVersion(templateID uint, id uint, currentUserID uint) error
	CreateRequirement(templateID uint, versionID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO) (uint, error)
	UpdateRequirement(templateID uint, versionID uint, id uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO) error
	DeleteRequirement(templateID uint, versionID uint, id uint, currentUserID uint) error
}

type ProjectSecurityGateChecklistTemplateService struct {
	repo        repositories.ProjectSecurityGateChecklistTemplateRepositoryInterface
	activityLog logServices.ActivityLogServiceInterface
}

func NewProjectSecurityGateChecklistTemplateService(
	repo repositories.ProjectSecurityGateChecklistTemplateRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
) ProjectSecurityGateChecklistTemplateServiceInterface {
	return &ProjectSecurityGateChecklistTemplateService{
		repo:        repo,
		activityLog: activityLog,
	}
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) GetAll() ([]*dtos.ProjectSecurityGateChecklistTemplateDTO, error) {
	templates, err := psgcts.repo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*dtos.ProjectSecurityGateChecklistTemplateDTO, 0, len(templates))
	for _, template := range templates {
		var dto dtos.ProjectSecurityGateChecklistTemplateDTO
		if err := mapstructure.Decode(template, &dto); err != nil {
			return nil, err
		}
		result = append(result, &dto)
	}

	return result, nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) GetByID(id uint) (*models.ProjectSecurityGateChecklistTemplate, error) {
	return psgcts.repo.GetByID(id)
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) GetByStageID(stageID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplate, error) {
	return psgcts.repo.GetByStageID(stageID, tx...)
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) GetLatestPublishedVersionByTemplateID(templateID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistTemplateVersion, error) {
	return psgcts.repo.GetLatestPublishedVersionByTemplateID(templateID, tx...)
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) GetDetailsByID(id uint) (*dtos.ProjectSecurityGateChecklistTemplateDetailDTO, error) {
	template, err := psgcts.repo.GetDetailsByID(id)
	if err != nil {
		return nil, err
	}

	var detail dtos.ProjectSecurityGateChecklistTemplateDetailDTO
	if err := mapstructure.Decode(template, &detail); err != nil {
		return nil, err
	}

	for i, versionDTO := range detail.Versions {
		versionDTO.CreatedAt = template.Versions[i].CreatedAt
		versionDTO.PublishedAt = template.Versions[i].PublishedAt
	}

	detail.LatestVersionID = latestPublishedVersionID(template.Versions)

	return &detail, nil
}

func latestPublishedVersionID(versions []models.ProjectSecurityGateChecklistTemplateVersion) *uint {
	for i := range versions {
		if versions[i].Status == "published" {
			id := versions[i].ID
			return &id
		}
	}
	return nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) CreateVersion(templateID uint, currentUserID uint) (id uint, err error) {
	tx := psgcts.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgcts.repo.GetByID(templateID, tx); err != nil {
		return 0, err
	}

	hasDraft, err := psgcts.repo.ExistsVersionByStatus(templateID, "draft", tx)
	if err != nil {
		return 0, err
	}
	if hasDraft {
		err = utility.NewValidationError(
			fmt.Errorf("an unpublished version already exists for this template; publish it before creating a new one"),
			"draft_exists",
			"version",
		)
		return 0, err
	}

	maxVersion, err := psgcts.repo.GetMaxVersionByTemplateID(templateID, tx)
	if err != nil {
		return 0, err
	}

	version := &models.ProjectSecurityGateChecklistTemplateVersion{
		TemplateID:  templateID,
		Version:     maxVersion + 1,
		Status:      "draft",
		CreatedByID: currentUserID,
		CreatedAt:   time.Now(),
	}

	created, err := psgcts.repo.CreateVersion(version, tx)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		err := psgcts.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGateChecklistTemplateVersion{}.TableName(),
			created.ID,
			created,
			currentUserID,
			//@todo: add module id and sub module id once module and submodule are ready
			"", "", "project_security_gate_checklist_template_version_create", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created.ID, nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) PublishVersion(templateID uint, id uint, currentUserID uint) (err error) {
	tx := psgcts.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	if _, err = psgcts.repo.GetByID(templateID, tx); err != nil {
		return err
	}

	version, err := psgcts.repo.GetVersionByID(id, tx)
	if err != nil {
		return err
	}

	if version.TemplateID != templateID {
		err = &utility.ResourceNotFoundError{}
		return err
	}

	if version.Status == "published" {
		err = utility.NewValidationError(
			fmt.Errorf("this version has already been published"),
			"already_published",
			"version",
		)
		return err
	}

	if len(version.Requirements) == 0 {
		err = utility.NewValidationError(
			fmt.Errorf("a version cannot be published without at least one checklist requirement"),
			"no_requirements",
			"requirement",
		)
		return err
	}

	version.Status = "published"
	version.PublishedByID = &currentUserID
	now := time.Now()
	version.PublishedAt = &now

	if err = psgcts.repo.UpdateVersion(version, tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		err := psgcts.activityLog.LogActivity(
			models.UpdateAction,
			models.ProjectSecurityGateChecklistTemplateVersion{}.TableName(),
			version.ID,
			map[string]interface{}{"status": "published"},
			currentUserID,
			//@todo: add module id and sub module id once module and submodule are ready
			"", "", "project_security_gate_checklist_template_version_publish", nil,
		)
		if err != nil {
			log.Printf("Failed to log publish activity: %v", err)
		}
	}()

	return nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) resolveEditableVersion(templateID uint, versionID uint) (*models.ProjectSecurityGateChecklistTemplateVersion, error) {
	version, err := psgcts.repo.GetVersionByID(versionID)
	if err != nil {
		return nil, err
	}

	if version.TemplateID != templateID {
		return nil, &utility.ResourceNotFoundError{}
	}

	if version.Status == "published" {
		return nil, utility.NewValidationError(
			fmt.Errorf("checklist requirements cannot be added, edited or deleted once the version is published"),
			"version_published",
			"version",
		)
	}

	return version, nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) CreateRequirement(templateID uint, versionID uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO) (uint, error) {
	if _, err := psgcts.resolveEditableVersion(templateID, versionID); err != nil {
		return 0, err
	}

	requirement := &models.ProjectSecurityGateChecklistTemplateRequirement{
		TemplateVersionID: versionID,
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		IsCritical:        form.IsCritical,
	}

	created, err := psgcts.repo.CreateRequirement(requirement)
	if err != nil {
		return 0, err
	}

	go func() {
		err := psgcts.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSecurityGateChecklistTemplateRequirement{}.TableName(),
			created.ID,
			created,
			currentUserID,
			//@todo: add module id and sub module id once module and submodule are ready
			"", "", "project_security_gate_checklist_template_requirement_create", nil,
		)
		if err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return created.ID, nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) UpdateRequirement(templateID uint, versionID uint, id uint, currentUserID uint, form *dtos.ProjectSecurityGateChecklistTemplateRequirementFormDTO) error {
	if _, err := psgcts.resolveEditableVersion(templateID, versionID); err != nil {
		return err
	}

	existing, err := psgcts.repo.GetRequirementByID(id)
	if err != nil {
		return err
	}
	if existing.TemplateVersionID != versionID {
		return &utility.ResourceNotFoundError{}
	}

	oldData := map[string]interface{}{
		"title_en":    existing.TitleEn,
		"title_ar":    existing.TitleAr,
		"is_critical": existing.IsCritical,
	}

	existing.TitleEn = form.TitleEn
	existing.TitleAr = form.TitleAr
	existing.IsCritical = form.IsCritical

	if err = psgcts.repo.UpdateRequirement(existing); err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"title_en":    form.TitleEn,
			"title_ar":    form.TitleAr,
			"is_critical": form.IsCritical,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := psgcts.activityLog.LogActivity(
				models.UpdateAction,
				models.ProjectSecurityGateChecklistTemplateRequirement{}.TableName(),
				id,
				changes,
				currentUserID,
				//@todo: add module id and sub module id once module and submodule are ready
				"", "", "project_security_gate_checklist_template_requirement_update", nil,
			)
			if err != nil {
				log.Printf("Failed to log update activity: %v", err)
			}
		}
	}()

	return nil
}

func (psgcts *ProjectSecurityGateChecklistTemplateService) DeleteRequirement(templateID uint, versionID uint, id uint, currentUserID uint) error {
	if _, err := psgcts.resolveEditableVersion(templateID, versionID); err != nil {
		return err
	}

	requirement, err := psgcts.repo.GetRequirementByID(id)
	if err != nil {
		return err
	}
	if requirement.TemplateVersionID != versionID {
		return &utility.ResourceNotFoundError{}
	}

	if err := psgcts.repo.DeleteRequirement(id); err != nil {
		return err
	}

	go func() {
		err := psgcts.activityLog.LogActivity(
			models.DeleteAction,
			models.ProjectSecurityGateChecklistTemplateRequirement{}.TableName(),
			id,
			map[string]interface{}{
				"title_en": requirement.TitleEn,
				"title_ar": requirement.TitleAr,
			},
			currentUserID,
			//@todo: add module id and sub module id once module and submodule are ready
			"", "", "project_security_gate_checklist_template_requirement_delete", nil,
		)
		if err != nil {
			log.Printf("Failed to log delete activity: %v", err)
		}
	}()

	return nil
}
