package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestFindingServiceInterface interface {
	ExportPenetrationTestFinding(data []*penetrationTestDTOs.PenetrationTestFindingListDTO, currentUserID uint) (*os.File, error)
	GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTestFinding, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, statusID uint, percentage *uint, tx utility.DatabaseTransaction) error
}

type PenetrationTestFindingService struct {
	userSettingSvc      organizationServices.UserSettingServiceInterface
	repo                repositories.PenetrationTestFindingRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	penetrationTestSvc  PenetrationTestServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
}

func NewPenetrationTestFindingService(
	repo repositories.PenetrationTestFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) PenetrationTestFindingServiceInterface {
	return &PenetrationTestFindingService{
		userSettingSvc:      userSettingSvc,
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		penetrationTestSvc:  penetrationTestSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
	}
}

func (ptfs *PenetrationTestFindingService) GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error) {
	return ptfs.repo.GetAll(params)
}

func (ptfs *PenetrationTestFindingService) GetByID(id uint) (*models.PenetrationTestFinding, error) {
	value, err := ptfs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (ptfs *PenetrationTestFindingService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error) {
	value, err := ptfs.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestFindingDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	return detail, nil
}

func (ptfs *PenetrationTestFindingService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error) {
	tx := ptfs.repo.BeginTransaction()
	if err := ptfs.validatePenetrationTestFindingFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	penetrationTest, err := ptfs.penetrationTestSvc.GetByID(formDTO.PenetrationTestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if penetrationTest.PenetrationTestStatus != nil && penetrationTest.PenetrationTestStatus.Type == "closed" {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("cannot add findings to a closed penetration test"), "penetration_test_closed", "penetration test")
	}
	var targetRemediationDate *time.Time
	if formDTO.TargetRemediationDate != nil && *formDTO.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*formDTO.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(err, "invalid_date", "target remediation date")
		}
	}
	actionPlanStatus, err := ptfs.actionPlanStatusSvc.GetByType("new")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value := &models.PenetrationTestFinding{
		PenetrationTestID:     formDTO.PenetrationTestID,
		Title:                 formDTO.Title,
		ActionPlanStatusID:    actionPlanStatus.ID,
		SeverityLevelID:       formDTO.SeverityLevelID,
		Description:           formDTO.Description,
		Recommendation:        formDTO.Recommendation,
		OwnerID:               formDTO.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}
	referenceNumber, err := ptfs.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.ReferenceNumber = referenceNumber
	value.CreatedBYID = currentUserID
	value.CreatedAt = time.Now()

	id, err := ptfs.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := ptfs.repo.CreateInitialPenetrationTestFindingStatusUpdate(&models.PenetrationTestFindingStatusUpdate{
		PenetrationTestFindingID: id,
		ActionPlanStatusID:       actionPlanStatus.ID,
		Comment:                  "Initial action plan status",
		CreatedBYID:              currentUserID,
		CreatedAt:                time.Now(),
	}, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingEntry, err := ptfs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := ptfs.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test finding", nil); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return id, nil
}

func (ptfs *PenetrationTestFindingService) Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingFormDTO) error {
	tx := ptfs.repo.BeginTransaction()
	existing, err := ptfs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := ptfs.validatePenetrationTestFindingFormDTO(form); err != nil {
		tx.Rollback()
		return err
	}

	if existing.PenetrationTestID != form.PenetrationTestID {
		penetrationTest, err := ptfs.penetrationTestSvc.GetByID(form.PenetrationTestID)
		if err != nil {
			tx.Rollback()
			return err
		}
		if penetrationTest.PenetrationTestStatus != nil && penetrationTest.PenetrationTestStatus.Type == "closed" {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("cannot add findings to a closed penetration test"), "penetration_test_closed", "penetration test")
		}
	}

	var targetRemediationDate *time.Time
	if form.TargetRemediationDate != nil && *form.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*form.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid_date", "target remediation date")
		}
	}
	updated := &models.PenetrationTestFinding{
		PenetrationTestID:     form.PenetrationTestID,
		Title:                 form.Title,
		ActionPlanStatusID:    existing.ActionPlanStatusID,
		Percentage:            existing.Percentage,
		SeverityLevelID:       form.SeverityLevelID,
		Description:           form.Description,
		Recommendation:        form.Recommendation,
		OwnerID:               form.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}
	updated.ID = existing.ID
	updated.ReferenceNumber = existing.ReferenceNumber
	updated.CreatedBYID = existing.CreatedBYID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedBYID = &currentUserID
	now := time.Now()
	updated.UpdatedAt = &now

	if err := ptfs.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		existingEntry, err := ptfs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := ptfs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test finding", nil); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (ptfs *PenetrationTestFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deleted []map[string]interface{}

	for _, id := range ids {
		value, err := ptfs.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		tx := ptfs.repo.BeginTransaction()
		if err := ptfs.repo.DeletePenetrationTestFindingStatusUpdates(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := ptfs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		folderPath := filepath.Join("penetration-testing", "penetration-test-findings", strconv.FormatUint(uint64(id), 10), "status-updates")
		if err := utility.RemoveFolder(folderPath); err != nil {
			utility.LogError("Failed to remove finding status update storage folder", map[string]interface{}{"findingID": id, "error": err.Error()})
		}

		successfulIDs = append(successfulIDs, id)
		deleted = append(deleted, map[string]interface{}{
			"id":    id,
			"title": value.Title,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := ptfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestFinding{}.TableName(),
				successfulIDs,
				deleted,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return failedIDs, nil
}

func (ptfs *PenetrationTestFindingService) StatusUpdate(id uint, currentUserID uint, statusID uint, percentage *uint, tx utility.DatabaseTransaction) error {
	finding, err := ptfs.repo.GetByID(id, tx)
	if err != nil {
		return err
	}
	penetrationTest, err := ptfs.penetrationTestSvc.GetByID(finding.PenetrationTestID)
	if err != nil {
		return err
	}
	if penetrationTest.PenetrationTestStatus != nil && penetrationTest.PenetrationTestStatus.Type == "closed" {
		status, err := ptfs.actionPlanStatusSvc.GetByID(statusID)
		if err != nil {
			return err
		}
		if status.Type != "close" {
			return utility.NewValidationError(fmt.Errorf("cannot reopen findings of a closed penetration test"), "penetration_test_closed", "penetration test")
		}
	}
	finding.ActionPlanStatusID = statusID
	finding.Percentage = percentage
	finding.UpdatedBYID = &currentUserID
	now := time.Now()
	finding.UpdatedAt = &now
	return ptfs.repo.Update(finding, tx)
}

func (ptfs *PenetrationTestFindingService) validatePenetrationTestFindingFormDTO(formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) error {
	if _, err := ptfs.penetrationTestSvc.GetByID(formDTO.PenetrationTestID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test does not exist"), "not_exist", "penetration test")
	}
	if _, err := ptfs.severityLevelSvc.GetByID(formDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}
	if _, err := ptfs.userSvc.GetByID(formDTO.OwnerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner does not exist"), "not_exist", "owner")
	}
	return nil
}

func (ptfs *PenetrationTestFindingService) ExportPenetrationTestFinding(data []*penetrationTestDTOs.PenetrationTestFindingListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ptfs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"الرقم التسلسلي",
			"الرقم المرجعي",
			"الملاحظة",
			"مرجع اختبار الاختراق",
			"النظام / التطبيق",
			"الوصف",
			"التوصية",
			"درجة الخطورة",
			"المستخدم المسؤول",
			"الوحدة المسؤولة",
			"التاريخ المستهدف",
			"الحالة",
			"نسبة الإنجاز (%)",
			"تاريخ الإنشاء",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Finding",
			"Pentest Reference",
			"System / Application",
			"Description",
			"Recommendation",
			"Severity",
			"Responsible User",
			"Responsible Unit",
			"Target Date",
			"Status",
			"Completion (%)",
			"Created At",
		}
	}
	rows := make([][]interface{}, len(data))
	for i, finding := range data {
		var responsibleUser, responsibleUnit, severity, status string

		if userSetting.Language == "ar" {
			responsibleUser = strings.TrimSpace(finding.OwnerFirstNameAr + " " + finding.OwnerLastNameAr)
			responsibleUnit = utility.GetStringValue(finding.ResponsibleUnitTitleAr)
			severity = finding.SeverityLevelTitleAr
			status = finding.ActionPlanStatusTitleAr
		} else {
			responsibleUser = strings.TrimSpace(finding.OwnerFirstNameEn + " " + finding.OwnerLastNameEn)
			responsibleUnit = utility.GetStringValue(finding.ResponsibleUnitTitleEn)
			severity = finding.SeverityLevelTitleEn
			status = finding.ActionPlanStatusTitleEn
		}

		var percentage uint
		if finding.Percentage != nil {
			percentage = *finding.Percentage
		}

		rows[i] = []interface{}{
			i + 1,
			finding.ReferenceNumber,
			finding.Title,
			finding.PenetrationTestReferenceNumber,
			finding.SystemApplicationName,
			utility.StripHTMLTags(utility.GetStringValue(finding.Description)),
			utility.StripHTMLTags(utility.GetStringValue(finding.Recommendation)),
			severity,
			responsibleUser,
			responsibleUnit,
			utility.GetStringValue(finding.TargetRemediationDate),
			status,
			percentage,
			finding.CreatedAt,
		}
	}
	return utility.GenerateExcelFile(headers, rows)
}
