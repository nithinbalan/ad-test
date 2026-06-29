package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessEmployeeAnalysisRepositoryInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessEmployeeAnalysisListQueryParamsDTO) (*orgDTOs.ProcessEmployeeAnalysisListResponseDTO, error)
	GetByID(id uint) (*models.ProcessEmployeeAnalysis, error)
	GetFullDetailsByID(id uint) (*models.ProcessEmployeeAnalysis, error)
	Create(item *models.ProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	Update(item *models.ProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateEmployee(employee *models.ProcessEmployeeAnalysisEmployee, tx ...utility.DatabaseTransaction) error
	DeleteEmployeesByAnalysisID(analysisID uint, tx ...utility.DatabaseTransaction) error
	BeginTransaction() utility.DatabaseTransaction
}

type ProcessEmployeeAnalysisRepository struct {
	db *gorm.DB
}

func NewProcessEmployeeAnalysisRepository(db *gorm.DB) *ProcessEmployeeAnalysisRepository {
	return &ProcessEmployeeAnalysisRepository{
		db: db,
	}
}

func (repo *ProcessEmployeeAnalysisRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProcessEmployeeAnalysisRepository) GetAll(processID uint, params *orgDTOs.ProcessEmployeeAnalysisListQueryParamsDTO) (*orgDTOs.ProcessEmployeeAnalysisListResponseDTO, error) {
	var result []*orgDTOs.ProcessEmployeeAnalysisListDTO
	var totalCount int64

	query := repo.db.Table("process_employee_analysis").
		Select(`
			process_employee_analysis.id,
			process_employee_analysis.process_id,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			process_employee_analysis.designation_id,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			process_employee_analysis.normal_count,
			process_employee_analysis.critical_count,
			process_employee_analysis.created_by_id,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			process_employee_analysis.updated_by_id,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, process_employee_analysis.created_at, 23) AS created_at,
			CONVERT(VARCHAR, process_employee_analysis.updated_at, 23) AS updated_at
		`).
		Joins(`LEFT JOIN processes ON processes.id = process_employee_analysis.process_id`).
		Joins(`LEFT JOIN designations ON designations.id = process_employee_analysis.designation_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = process_employee_analysis.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = process_employee_analysis.updated_by_id`).
		Where("process_employee_analysis.process_id = ?", processID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			designations.title_en LIKE ? OR
			designations.title_ar LIKE ?
		`, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}

	if err := repo.populateEmployeeAnalysisListAssociations(result); err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &orgDTOs.ProcessEmployeeAnalysisListResponseDTO{
		Items: result,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ProcessEmployeeAnalysisRepository) populateEmployeeAnalysisListAssociations(analyses []*orgDTOs.ProcessEmployeeAnalysisListDTO) error {
	if len(analyses) == 0 {
		return nil
	}

	analysisIDs := make([]uint, len(analyses))
	for i, a := range analyses {
		analysisIDs[i] = a.ID
	}

	var employees []*models.ProcessEmployeeAnalysisEmployee
	if err := repo.db.Where("process_employee_analysis_id IN ?", analysisIDs).Preload("User").Find(&employees).Error; err != nil {
		return err
	}

	employeeMap := make(map[uint][]orgDTOs.ProcessEmployeeAnalysisEmployeeDTO)
	for _, emp := range employees {
		if emp == nil || emp.User == nil {
			continue
		}
		employeeMap[emp.ProcessEmployeeAnalysisID] = append(employeeMap[emp.ProcessEmployeeAnalysisID], orgDTOs.ProcessEmployeeAnalysisEmployeeDTO{
			ID:                        emp.ID,
			ProcessEmployeeAnalysisID: emp.ProcessEmployeeAnalysisID,
			UserID:                    emp.UserID,
			User: orgDTOs.UserMinimalDTO{
				ID:               int(emp.User.ID),
				FirstNameEN:      emp.User.FirstNameEn,
				FirstNameAR:      emp.User.FirstNameAr,
				LastNameEN:       emp.User.LastNameEn,
				LastNameAR:       emp.User.LastNameAr,
				Initial:          emp.User.Initial,
				Color:            emp.User.Color,
				Email:            emp.User.Email,
				ProfileImagePath: emp.User.ProfileImagePath,
			},
		})
	}

	for _, analysis := range analyses {
		if emps, ok := employeeMap[analysis.ID]; ok {
			analysis.Employees = emps
		} else {
			analysis.Employees = []orgDTOs.ProcessEmployeeAnalysisEmployeeDTO{}
		}
	}

	return nil
}

func (repo *ProcessEmployeeAnalysisRepository) GetByID(id uint) (*models.ProcessEmployeeAnalysis, error) {
	var item *models.ProcessEmployeeAnalysis
	err := repo.db.First(&item, id).Error
	return item, err
}

func (repo *ProcessEmployeeAnalysisRepository) GetFullDetailsByID(id uint) (*models.ProcessEmployeeAnalysis, error) {
	var item *models.ProcessEmployeeAnalysis
	err := repo.db.
		Preload("Process").
		Preload("Designation").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Employees").
		First(&item, id).Error
	return item, err
}

func (repo *ProcessEmployeeAnalysisRepository) Create(item *models.ProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(item).Error
	return item.ID, err
}

func (repo *ProcessEmployeeAnalysisRepository) Update(item *models.ProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProcessEmployeeAnalysis{}).
		Where("id = ?", item.ID).
		Select("DesignationID", "NormalCount", "CriticalCount", "UpdatedBYID").
		Updates(item).Error
}

func (repo *ProcessEmployeeAnalysisRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ProcessEmployeeAnalysis{}, id).Error
}

func (repo *ProcessEmployeeAnalysisRepository) CreateEmployee(employee *models.ProcessEmployeeAnalysisEmployee, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(employee).Error
}

func (repo *ProcessEmployeeAnalysisRepository) DeleteEmployeesByAnalysisID(analysisID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_employee_analysis_id = ?", analysisID).
		Delete(&models.ProcessEmployeeAnalysisEmployee{}).Error
}
