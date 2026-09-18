package repositories

import (
	"errors"
	"strconv"
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateReportRepositoryInterface interface {
	GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error)
	GetAllChecklists(securityGateID uint) ([]*models.ProjectSecurityGateChecklist, error)
	GetAllChecklistsByResponseStatusIDs(securityGateID uint, responseStatusIDs []uint) ([]*models.ProjectSecurityGateChecklist, error)
	GetAllExceptionWorkflowHistories(securityGateID uint) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error)
	GetAllResponseDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistResponseDocument, error)
	GetAllWorkflowCommentDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistCommentDocument, error)

	GetGatesByIDs(securityGateIDs []uint) ([]*models.ProjectSecurityGate, error)
	GetResponseStatusesByIDs(responseStatusIDs []uint) ([]*models.ProjectSecurityGateChecklistResponseStatus, error)
	GetUserByID(id uint) (*models.User, error)

	BeginTransaction() utility.DatabaseTransaction

	CreateSecurityGateReportSection(section *models.ProjectSecurityGateReportSection, tx ...utility.DatabaseTransaction) error
	GetSecurityGateReportSections(reportID uint) ([]*models.ProjectSecurityGateReportSection, error)
	GetSecurityGateReportSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReportSection, error)
	UpdateSecurityGateReportSection(section *models.ProjectSecurityGateReportSection, tx ...utility.DatabaseTransaction) error
	DeleteSecurityGateReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderSecurityGateReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error

	CreateSecurityGateExceptionReportSection(section *models.ProjectSecurityGateExceptionReportSection, tx ...utility.DatabaseTransaction) error
	GetSecurityGateExceptionReportSections(reportID uint) ([]*models.ProjectSecurityGateExceptionReportSection, error)
	GetSecurityGateExceptionReportSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateExceptionReportSection, error)
	UpdateSecurityGateExceptionReportSection(section *models.ProjectSecurityGateExceptionReportSection, tx ...utility.DatabaseTransaction) error
	DeleteSecurityGateExceptionReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderSecurityGateExceptionReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error

	CreateSecurityGateReport(report *models.ProjectSecurityGateReport, tx ...utility.DatabaseTransaction) (uint, error)
	GetAllSecurityGateReports(params *projectDTOs.ProjectSecurityGateReportListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateReportListResponseDTO, error)
	GetSecurityGateReportByID(id uint) (*models.ProjectSecurityGateReport, error)
	DeleteSecurityGateReport(id uint) error

	CreateSecurityGateExceptionReport(report *models.ProjectSecurityGateExceptionReport, tx ...utility.DatabaseTransaction) (uint, error)
	GetAllSecurityGateExceptionReports(params *projectDTOs.ProjectSecurityGateReportListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateReportListResponseDTO, error)
	GetSecurityGateExceptionReportByID(id uint) (*models.ProjectSecurityGateExceptionReport, error)
	DeleteSecurityGateExceptionReport(id uint) error
}

type ProjectSecurityGateReportRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateReportRepository(db *gorm.DB) ProjectSecurityGateReportRepositoryInterface {
	return &ProjectSecurityGateReportRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateReportRepository) GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error) {
	var projectSecurityGate *models.ProjectSecurityGate
	result := repo.db.
		Preload("Stage").
		Preload("TemplateVersion").
		Preload("BusinessOwner").
		Preload("ReadinessStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Classifications").
		Preload("Classifications.Classification").
		First(&projectSecurityGate, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectSecurityGate, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllChecklists(securityGateID uint) ([]*models.ProjectSecurityGateChecklist, error) {
	var checklists []*models.ProjectSecurityGateChecklist
	result := repo.db.
		Select("project_security_gate_checklists.*").
		Preload("Requirement").
		Preload("ResponseStatus").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Risks").
		Preload("Risks.RiskRegister").
		Preload("Risks.RiskRegister.RiskStatus").
		Preload("Risks.RiskRegister.RiskCategory").
		Preload("Risks.RiskRegister.ResponsibleUnit").
		Preload("Risks.RiskRegister.ResponsibleUser").
		Joins("JOIN project_security_gate_checklist_template_requirements AS requirements ON requirements.id = project_security_gate_checklists.requirement_id").
		Where("project_security_gate_checklists.security_gate_id = ?", securityGateID).
		Order("requirements.is_critical DESC, project_security_gate_checklists.id ASC").
		Find(&checklists)

	if result.Error != nil {
		return nil, result.Error
	}

	return checklists, nil
}

// GetAllChecklistsByResponseStatusIDs keeps only the responses the report asked for
func (repo *ProjectSecurityGateReportRepository) GetAllChecklistsByResponseStatusIDs(securityGateID uint, responseStatusIDs []uint) ([]*models.ProjectSecurityGateChecklist, error) {
	var checklists []*models.ProjectSecurityGateChecklist
	result := repo.db.
		Select("project_security_gate_checklists.*").
		Preload("Requirement").
		Preload("ResponseStatus").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Risks").
		Preload("Risks.RiskRegister").
		Preload("Risks.RiskRegister.RiskStatus").
		Preload("Risks.RiskRegister.RiskCategory").
		Preload("Risks.RiskRegister.ResponsibleUnit").
		Preload("Risks.RiskRegister.ResponsibleUser").
		Joins("JOIN project_security_gate_checklist_template_requirements AS requirements ON requirements.id = project_security_gate_checklists.requirement_id").
		Where("project_security_gate_checklists.security_gate_id = ?", securityGateID).
		Where("project_security_gate_checklists.response_status_id IN (?)", responseStatusIDs).
		Order("requirements.is_critical DESC, project_security_gate_checklists.id ASC").
		Find(&checklists)

	if result.Error != nil {
		return nil, result.Error
	}

	return checklists, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllExceptionWorkflowHistories(securityGateID uint) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error) {
	var histories []*models.ProjectSecurityGateChecklistExceptionWorkflowHistory
	result := repo.db.
		Select("project_security_gate_checklist_exception_workflow_histories.*").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = project_security_gate_checklist_exception_workflow_histories.checklist_id").
		Where("checklists.security_gate_id = ?", securityGateID).
		Order("project_security_gate_checklist_exception_workflow_histories.checklist_id ASC, project_security_gate_checklist_exception_workflow_histories.created_at ASC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllResponseDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistResponseDocument, error) {
	var documents []*models.ProjectSecurityGateChecklistResponseDocument
	result := repo.db.
		Select("project_security_gate_checklist_response_documents.*").
		Preload("Checklist").
		Preload("Checklist.Requirement").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = project_security_gate_checklist_response_documents.project_security_gate_checklist_id").
		Where("checklists.security_gate_id = ?", securityGateID).
		Order("project_security_gate_checklist_response_documents.project_security_gate_checklist_id ASC, project_security_gate_checklist_response_documents.id ASC").
		Find(&documents)

	if result.Error != nil {
		return nil, result.Error
	}

	return documents, nil
}

// GetAllWorkflowCommentDocuments skips plain discussion comments
func (repo *ProjectSecurityGateReportRepository) GetAllWorkflowCommentDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistCommentDocument, error) {
	var documents []*models.ProjectSecurityGateChecklistCommentDocument
	result := repo.db.
		Select("project_security_gate_checklist_comment_documents.*").
		Preload("Comment").
		Preload("Comment.Checklist").
		Preload("Comment.Checklist.Requirement").
		Joins("JOIN project_security_gate_checklist_comments AS comments ON comments.id = project_security_gate_checklist_comment_documents.project_security_gate_checklist_comment_id").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = comments.checklist_id").
		Where("checklists.security_gate_id = ? AND comments.workflow_history_id IS NOT NULL", securityGateID).
		Order("comments.checklist_id ASC, comments.id ASC, project_security_gate_checklist_comment_documents.id ASC").
		Find(&documents)

	if result.Error != nil {
		return nil, result.Error
	}

	return documents, nil
}

func (repo *ProjectSecurityGateReportRepository) GetGatesByIDs(securityGateIDs []uint) ([]*models.ProjectSecurityGate, error) {
	var securityGates []*models.ProjectSecurityGate
	result := repo.db.
		Where("id IN (?)", securityGateIDs).
		Order("project_title_en ASC").
		Find(&securityGates)

	if result.Error != nil {
		return nil, result.Error
	}

	return securityGates, nil
}

func (repo *ProjectSecurityGateReportRepository) GetResponseStatusesByIDs(responseStatusIDs []uint) ([]*models.ProjectSecurityGateChecklistResponseStatus, error) {
	var responseStatuses []*models.ProjectSecurityGateChecklistResponseStatus
	result := repo.db.
		Where("id IN (?)", responseStatusIDs).
		Order("id ASC").
		Find(&responseStatuses)

	if result.Error != nil {
		return nil, result.Error
	}

	return responseStatuses, nil
}

func (repo *ProjectSecurityGateReportRepository) GetUserByID(id uint) (*models.User, error) {
	var user *models.User
	result := repo.db.First(&user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (repo *ProjectSecurityGateReportRepository) CreateSecurityGateReport(report *models.ProjectSecurityGateReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}

	return report.ID, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllSecurityGateReports(params *projectDTOs.ProjectSecurityGateReportListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateReportListResponseDTO, error) {
	query := repo.db.Table("project_security_gate_reports").
		Select(`project_security_gate_reports.id,
			project_security_gate_reports.title,
			project_security_gate_reports.created_at,
			project_security_gate_reports.security_gate_ids,
			project_security_gate_reports.response_status_ids,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_security_gate_reports.created_by_id")

	if params.Page < 1 {
		params.Page = 1
	}

	if params.PageSize < 1 {
		params.PageSize = 10
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_security_gate_reports."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(project_security_gate_reports.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	var rawReports []struct {
		ID                          uint   `gorm:"column:id"`
		Title                       string `gorm:"column:title"`
		CreatedAt                   string `gorm:"column:created_at"`
		SecurityGateIDs             string `gorm:"column:security_gate_ids"`
		ResponseStatusIDs           string `gorm:"column:response_status_ids"`
		CreatedUserID               uint   `gorm:"column:created_user_id"`
		CreatedUserFirstNameEn      string `gorm:"column:created_user_first_name_en"`
		CreatedUserFirstNameAr      string `gorm:"column:created_user_first_name_ar"`
		CreatedUserLastNameEn       string `gorm:"column:created_user_last_name_en"`
		CreatedUserLastNameAr       string `gorm:"column:created_user_last_name_ar"`
		CreatedUserInitial          string `gorm:"column:created_user_initial"`
		CreatedUserColor            string `gorm:"column:created_user_color"`
		CreatedUserEmail            string `gorm:"column:created_user_email"`
		CreatedUserProfileImagePath string `gorm:"column:created_user_profile_image_path"`
	}
	if err := query.
		Order(params.SortBy + " " + params.Sort).
		Offset(offset).
		Limit(params.PageSize).
		Find(&rawReports).Error; err != nil {
		return nil, err
	}

	reports := make([]*projectDTOs.ProjectSecurityGateReportListDTO, len(rawReports))
	for i, raw := range rawReports {
		securityGateIDs := repo.parseIDList(raw.SecurityGateIDs)
		responseStatusIDs := repo.parseIDList(raw.ResponseStatusIDs)

		securityGates := repo.getSecurityGateTitles(securityGateIDs)
		responseStatuses := repo.getResponseStatusTitles(responseStatusIDs)

		reports[i] = &projectDTOs.ProjectSecurityGateReportListDTO{
			ID:                          raw.ID,
			Title:                       raw.Title,
			CreatedAt:                   raw.CreatedAt,
			CreatedUserID:               raw.CreatedUserID,
			CreatedUserFirstNameEn:      raw.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      raw.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       raw.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       raw.CreatedUserLastNameAr,
			CreatedUserInitial:          raw.CreatedUserInitial,
			CreatedUserColor:            raw.CreatedUserColor,
			CreatedUserEmail:            raw.CreatedUserEmail,
			CreatedUserProfileImagePath: raw.CreatedUserProfileImagePath,
			SecurityGates:               securityGates,
			ResponseStatuses:            responseStatuses,
		}
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &projectDTOs.ProjectSecurityGateReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateReportByID(id uint) (*models.ProjectSecurityGateReport, error) {
	var report *models.ProjectSecurityGateReport
	result := repo.db.
		Preload("CreatedByUser").
		First(&report, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return report, nil
}

func (repo *ProjectSecurityGateReportRepository) DeleteSecurityGateReport(id uint) error {
	return repo.db.Where("id = ?", id).Delete(&models.ProjectSecurityGateReport{}).Error
}

func (repo *ProjectSecurityGateReportRepository) CreateSecurityGateExceptionReport(report *models.ProjectSecurityGateExceptionReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}

	return report.ID, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllSecurityGateExceptionReports(params *projectDTOs.ProjectSecurityGateReportListQueryParamsDTO) (*projectDTOs.ProjectSecurityGateReportListResponseDTO, error) {
	query := repo.db.Table("project_security_gate_exception_reports").
		Select(`project_security_gate_exception_reports.id,
			project_security_gate_exception_reports.title,
			project_security_gate_exception_reports.created_at,
			project_security_gate_exception_reports.security_gate_ids,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_security_gate_exception_reports.created_by_id")

	if params.Page < 1 {
		params.Page = 1
	}

	if params.PageSize < 1 {
		params.PageSize = 10
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_security_gate_exception_reports."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(project_security_gate_exception_reports.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	var rawReports []struct {
		ID                          uint   `gorm:"column:id"`
		Title                       string `gorm:"column:title"`
		CreatedAt                   string `gorm:"column:created_at"`
		SecurityGateIDs             string `gorm:"column:security_gate_ids"`
		CreatedUserID               uint   `gorm:"column:created_user_id"`
		CreatedUserFirstNameEn      string `gorm:"column:created_user_first_name_en"`
		CreatedUserFirstNameAr      string `gorm:"column:created_user_first_name_ar"`
		CreatedUserLastNameEn       string `gorm:"column:created_user_last_name_en"`
		CreatedUserLastNameAr       string `gorm:"column:created_user_last_name_ar"`
		CreatedUserInitial          string `gorm:"column:created_user_initial"`
		CreatedUserColor            string `gorm:"column:created_user_color"`
		CreatedUserEmail            string `gorm:"column:created_user_email"`
		CreatedUserProfileImagePath string `gorm:"column:created_user_profile_image_path"`
	}
	if err := query.
		Order(params.SortBy + " " + params.Sort).
		Offset(offset).
		Limit(params.PageSize).
		Find(&rawReports).Error; err != nil {
		return nil, err
	}

	reports := make([]*projectDTOs.ProjectSecurityGateReportListDTO, len(rawReports))
	for i, raw := range rawReports {
		securityGateIDs := repo.parseIDList(raw.SecurityGateIDs)

		securityGates := repo.getSecurityGateTitles(securityGateIDs)

		reports[i] = &projectDTOs.ProjectSecurityGateReportListDTO{
			ID:                          raw.ID,
			Title:                       raw.Title,
			CreatedAt:                   raw.CreatedAt,
			CreatedUserID:               raw.CreatedUserID,
			CreatedUserFirstNameEn:      raw.CreatedUserFirstNameEn,
			CreatedUserFirstNameAr:      raw.CreatedUserFirstNameAr,
			CreatedUserLastNameEn:       raw.CreatedUserLastNameEn,
			CreatedUserLastNameAr:       raw.CreatedUserLastNameAr,
			CreatedUserInitial:          raw.CreatedUserInitial,
			CreatedUserColor:            raw.CreatedUserColor,
			CreatedUserEmail:            raw.CreatedUserEmail,
			CreatedUserProfileImagePath: raw.CreatedUserProfileImagePath,
			SecurityGates:               securityGates,
			ResponseStatuses:            []*projectDTOs.ResponseStatusTitleDTO{},
		}
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &projectDTOs.ProjectSecurityGateReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateExceptionReportByID(id uint) (*models.ProjectSecurityGateExceptionReport, error) {
	var report *models.ProjectSecurityGateExceptionReport
	result := repo.db.
		Preload("CreatedByUser").
		First(&report, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return report, nil
}

func (repo *ProjectSecurityGateReportRepository) DeleteSecurityGateExceptionReport(id uint) error {
	return repo.db.Where("id = ?", id).Delete(&models.ProjectSecurityGateExceptionReport{}).Error
}

func (repo *ProjectSecurityGateReportRepository) parseIDList(value string) []uint {
	ids := []uint{}
	for _, part := range strings.Split(value, ",") {
		id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id))
	}
	return ids
}

func (repo *ProjectSecurityGateReportRepository) getSecurityGateTitles(ids []uint) []*projectDTOs.SecurityGateTitleDTO {
	if len(ids) == 0 {
		return []*projectDTOs.SecurityGateTitleDTO{}
	}
	var gates []struct {
		ID      uint   `gorm:"column:id"`
		TitleEn string `gorm:"column:project_title_en"`
		TitleAr string `gorm:"column:project_title_ar"`
	}
	repo.db.Table("project_security_gates").Select("id, project_title_en, project_title_ar").Where("id IN (?)", ids).Find(&gates)
	result := make([]*projectDTOs.SecurityGateTitleDTO, len(gates))
	for i, g := range gates {
		result[i] = &projectDTOs.SecurityGateTitleDTO{ID: g.ID, TitleEn: g.TitleEn, TitleAr: g.TitleAr}
	}
	return result
}

func (repo *ProjectSecurityGateReportRepository) getResponseStatusTitles(ids []uint) []*projectDTOs.ResponseStatusTitleDTO {
	if len(ids) == 0 {
		return []*projectDTOs.ResponseStatusTitleDTO{}
	}
	var statuses []struct {
		ID        uint   `gorm:"column:id"`
		TitleEn   string `gorm:"column:title_en"`
		TitleAr   string `gorm:"column:title_ar"`
		ColorCode string `gorm:"column:color_code"`
	}
	repo.db.Table("project_security_gate_checklist_response_statuses").Select("id, title_en, title_ar, color_code").Where("id IN (?)", ids).Find(&statuses)
	result := make([]*projectDTOs.ResponseStatusTitleDTO, len(statuses))
	for i, s := range statuses {
		result[i] = &projectDTOs.ResponseStatusTitleDTO{ID: s.ID, TitleEn: s.TitleEn, TitleAr: s.TitleAr, ColorCode: s.ColorCode}
	}
	return result
}

func (repo *ProjectSecurityGateReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProjectSecurityGateReportRepository) CreateSecurityGateReportSection(section *models.ProjectSecurityGateReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(section).Error
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateReportSections(reportID uint) ([]*models.ProjectSecurityGateReportSection, error) {
	sections := []*models.ProjectSecurityGateReportSection{}
	result := repo.db.
		Where("report_id = ?", reportID).
		Order("order_num ASC").
		Find(&sections)

	if result.Error != nil {
		return nil, result.Error
	}

	return sections, nil
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateReportSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section *models.ProjectSecurityGateReportSection
	result := db.First(&section, sectionID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, sectionID)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return section, nil
}

func (repo *ProjectSecurityGateReportRepository) UpdateSecurityGateReportSection(section *models.ProjectSecurityGateReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(section).Error
}

func (repo *ProjectSecurityGateReportRepository) DeleteSecurityGateReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.ProjectSecurityGateReportSection{}).Error
}

// ReorderSecurityGateReportSections closes the gap left by a deleted section
func (repo *ProjectSecurityGateReportRepository) ReorderSecurityGateReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ProjectSecurityGateReportSection{}).
		Where("report_id = ? AND order_num > ?", reportID, deletedOrderNum).
		UpdateColumn("order_num", gorm.Expr("order_num - 1")).Error
}

func (repo *ProjectSecurityGateReportRepository) CreateSecurityGateExceptionReportSection(section *models.ProjectSecurityGateExceptionReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(section).Error
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateExceptionReportSections(reportID uint) ([]*models.ProjectSecurityGateExceptionReportSection, error) {
	sections := []*models.ProjectSecurityGateExceptionReportSection{}
	result := repo.db.
		Where("report_id = ?", reportID).
		Order("order_num ASC").
		Find(&sections)

	if result.Error != nil {
		return nil, result.Error
	}

	return sections, nil
}

func (repo *ProjectSecurityGateReportRepository) GetSecurityGateExceptionReportSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateExceptionReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section *models.ProjectSecurityGateExceptionReportSection
	result := db.First(&section, sectionID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, sectionID)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return section, nil
}

func (repo *ProjectSecurityGateReportRepository) UpdateSecurityGateExceptionReportSection(section *models.ProjectSecurityGateExceptionReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(section).Error
}

func (repo *ProjectSecurityGateReportRepository) DeleteSecurityGateExceptionReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.ProjectSecurityGateExceptionReportSection{}).Error
}

// ReorderSecurityGateExceptionReportSections closes the gap left by a deleted section
func (repo *ProjectSecurityGateReportRepository) ReorderSecurityGateExceptionReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ProjectSecurityGateExceptionReportSection{}).
		Where("report_id = ? AND order_num > ?", reportID, deletedOrderNum).
		UpdateColumn("order_num", gorm.Expr("order_num - 1")).Error
}
