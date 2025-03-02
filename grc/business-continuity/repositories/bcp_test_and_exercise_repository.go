package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpTestAndExerciseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.BcpTestAndExerciseListResponseDTO, error)
	GetByID(id uint) (*models.BcpTestAndExercise, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpTestAndExercise, error)
	Create(BcpTestAndExercise *models.BcpTestAndExercise, tx ...utility.DatabaseTransaction) (uint, error)
	Update(BcpTestAndExercise *models.BcpTestAndExercise, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateBcpTestAndExerciseBcpProcesses(bcpProcesses []*models.BcpTestAndExerciseBcpProcess, tx ...utility.DatabaseTransaction) error
	DeleteBcpTestAndExerciseBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error
}

type BcpTestAndExerciseRepository struct {
	db *gorm.DB
}

func NewBcpTestAndExerciseRepository(db *gorm.DB) *BcpTestAndExerciseRepository {
	return &BcpTestAndExerciseRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BcpTestAndExerciseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpTestAndExerciseRepository) GetAll(bcpId uint, params *businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.BcpTestAndExerciseListResponseDTO, error) {
	var BcpTestAndExercises []*businessContinuityDTOs.BcpTestAndExerciseListDTO
	var totalCount int64

	query := repo.db.Table("bcp_test_and_exercises").
		Select(`
            bcp_test_and_exercises.*,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path  
        `).
		Joins(`
			LEFT JOIN review_frequencies ON review_frequencies.id = bcp_test_and_exercises.review_frequency_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = bcp_test_and_exercises.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = bcp_test_and_exercises.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = bcp_test_and_exercises.updated_by_id
        `).Where("bcp_test_and_exercises.business_continuity_plan_id = ?", bcpId)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(bcp_test_and_exercises.title LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&BcpTestAndExercises).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &businessContinuityDTOs.BcpTestAndExerciseListResponseDTO{
		Items: BcpTestAndExercises,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BcpTestAndExerciseRepository) GetByID(id uint) (*models.BcpTestAndExercise, error) {
	var BcpTestAndExercise *models.BcpTestAndExercise

	result := repo.db.First(&BcpTestAndExercise, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return BcpTestAndExercise, nil

}

func (repo *BcpTestAndExerciseRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpTestAndExercise, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var BcpTestAndExercise *models.BcpTestAndExercise

	result := db.
		Preload("BusinessContinuityPlan").
		Preload("BusinessContinuityPlan.BusinessFunction").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ReviewFrequency").
		Preload("BcpTestAndExerciseBcpProcesses").
		Preload("BcpTestAndExerciseBcpProcesses.BcpProcess").
		Preload("BcpTestAndExerciseBcpProcesses.BcpProcess.Process").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&BcpTestAndExercise, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return BcpTestAndExercise, nil
}

func (repo *BcpTestAndExerciseRepository) Create(BcpTestAndExercise *models.BcpTestAndExercise, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(BcpTestAndExercise).Error; err != nil {
		return 0, err
	}

	return BcpTestAndExercise.ID, nil
}

func (repo *BcpTestAndExerciseRepository) Update(BcpTestAndExercise *models.BcpTestAndExercise, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(BcpTestAndExercise).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpTestAndExerciseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Communication Plan
	if err := db.Delete(&models.BcpTestAndExercise{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpTestAndExerciseRepository) CreateBcpTestAndExerciseBcpProcesses(bcpProcesses []*models.BcpTestAndExerciseBcpProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&bcpProcesses).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BcpTestAndExerciseRepository) DeleteBcpTestAndExerciseBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", bcpProcessIds).Delete(&models.BcpTestAndExerciseBcpProcess{}).Error; err != nil {
		return err
	}

	return nil
}
