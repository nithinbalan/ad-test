package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KPIRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *kpiDTOs.KPIRegisterListQueryParamsDTO) (*kpiDTOs.KPIRegisterListResponseDTO, error)
	GetByID(id uint) (*models.KPIRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIRegister, error)
	Create(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(ids []uint, tx ...utility.DatabaseTransaction) error
	SetAllIsLatestFalse(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error
	CreateKPIValueUpdate(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error)

	CreateKPIRegisterInpuData(inputs []*models.KPIRegisterDataInput, tx ...utility.DatabaseTransaction) error
	DeleteKPIRegisterInputData(kpiRegisterId uint, tx ...utility.DatabaseTransaction) error

	GetTrendAnalysis(kpiId uint, limit int) ([]*models.KPIValueUpdate, error)
	GetLatestValueUpdate(kpiId uint) (*models.KPIValueUpdate, error)
	GenerateReferenceNumber() (string, error)
	GetKPICategorySummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPICategorySummaryDTO, error)
	GetKPIStatusSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPIStatusSummaryDTO, error)
	GetKPIOverallSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIOverallSummaryDTO, error)
	GetKPIByReviewStatus(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPIReviewStatusDistributionDTO, error)
	CreateKPIValueUpdateInputData(inputDatas []*models.KPIValueUpdateDataInput, tx ...utility.DatabaseTransaction) error
	UpdateKPIValueUpdate(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) error
	DeleteKPIValueUpdateInputData(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error
	ActivateKPIRegisters(kpiRegisterIDs []uint, archiveStatusID uint, activeStatusID uint) ([]uint, error)
	ArchiveKPIRegisters(kpiRegisterIDs []uint, archiveStatusID uint, activeStatusID uint) ([]uint, error)
	GetAllResponsibleUsers(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
}

type KPIRegisterRepository struct {
	db *gorm.DB
}

func NewKPIRegisterRepository(db *gorm.DB) *KPIRegisterRepository {
	return &KPIRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *KPIRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KPIRegisterRepository) GetAll(params *kpiDTOs.KPIRegisterListQueryParamsDTO) (*kpiDTOs.KPIRegisterListResponseDTO, error) {
	var kpiRegister []*kpiDTOs.KPIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kpi_register").
		Select(`
            kpi_register.*,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS end_date,
			CONVERT(VARCHAR, kpi_register.target_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			kpi_review_statuses.title_en AS kpi_review_status_title_en,
			kpi_review_statuses.title_ar AS kpi_review_status_title_ar,
			kpi_review_statuses.type AS kpi_review_status_type,
			kpi_review_statuses.color_code AS kpi_review_status_color_code,
			lifecycle_statuses.title_en AS status_title_en,
			lifecycle_statuses.title_ar AS status_title_ar,
			lifecycle_statuses.type AS status_type,
			lifecycle_statuses.color_code AS status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
			kpi_achievement_logics.title_en AS kpi_achievement_logic_title_en,
			kpi_achievement_logics.title_ar AS kpi_achievement_logic_title_ar,
			kpi_achievement_logics.type AS kpi_achievement_logic_type,
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
            updated_user.profile_image_path AS updated_user_profile_image_path,
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path
        `).
		Joins(`
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN kpi_review_statuses ON kpi_review_statuses.id = kpi_register.kpi_review_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = kpi_register.status_id
			LEFT JOIN kpi_achievement_logics ON kpi_achievement_logics.id = kpi_register.kpi_achievement_logic_id
			LEFT JOIN users AS review_user ON review_user.id = kpi_register.review_user_id
        `)
	// Apply dynamic filters to the query

	for columnName, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		switch columnName {
		// ✅ KPI ID
		case "kpi_id":
			query = query.Where("kpi_register.id IN ?", values)

		// ✅ KPI Register Type (manual / automatic)
		case "type":
			query = query.Where("kpi_register.type IN ?", values)
		// ✅ KPI Status Type
		// case "kpi_status_type":
		// 	query = query.Where("kpi_statuses.type IN ?", values)
		// ✅ Default
		default:
			query = query.Where(columnName+" IN ?", values)
		}
	}
	// Calculate total count for pagination
	// if err := query.Count(&totalCount).Error; err != nil {
	// 	return nil, err
	// }

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"

		query = query.Where(`
		kpi_register.title LIKE ? OR
		kpi_register.reference_number LIKE ? OR

		kpi_categories.title_en LIKE ? OR
		kpi_categories.title_ar LIKE ? OR

		kpi_statuses.title_en LIKE ? OR
		kpi_statuses.title_ar LIKE ? OR

		kpi_review_statuses.title_en LIKE ? OR
		kpi_review_statuses.title_ar LIKE ? OR

		review_frequencies.title_en LIKE ? OR
		review_frequencies.title_ar LIKE ? OR

		calculation_types.title_en LIKE ? OR
		calculation_types.title_ar LIKE ? OR

		units.title_en LIKE ? OR
		units.title_ar LIKE ? OR
		units.code_en LIKE ? OR
		units.code_ar LIKE ? OR

		kpi_achievement_logics.title_en LIKE ? OR
		kpi_achievement_logics.title_ar LIKE ? OR

		responsible_user.first_name_en LIKE ? OR
		responsible_user.first_name_ar LIKE ? OR
		responsible_user.last_name_en LIKE ? OR
		responsible_user.last_name_ar LIKE ? OR

		review_user.first_name_en LIKE ? OR
		review_user.first_name_ar LIKE ? OR
		review_user.last_name_en LIKE ? OR
		review_user.last_name_ar LIKE ? OR

		created_user.first_name_en LIKE ? OR
		created_user.first_name_ar LIKE ? OR
		created_user.last_name_en LIKE ? OR
		created_user.last_name_ar LIKE ? OR

		updated_user.first_name_en LIKE ? OR
		updated_user.first_name_ar LIKE ? OR
		updated_user.last_name_en LIKE ? OR
		updated_user.last_name_ar LIKE ?
	`,
			keyword, keyword,

			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword, keyword, keyword,

			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
		)
	}
	fmt.Println("params.CurrentUserID=", params.CurrentUserID)
	for _, filter := range params.QuickFilter {
		switch filter {
		case "my-kpis":
			query = query.Where("kpi_register.responsible_user_id = ?", params.CurrentUserID)
		case "awaiting-update":
			query = query.Where("kpi_review_statuses.type = ?", "awaiting-update")
		case "pending":
			query = query.Where("kpi_review_statuses.type = ?", "awaiting-update")
		case "overdue":
			query = query.Where(`kpi_review_statuses.type = ?`, "overdue")
		case "achieved":
			query = query.Where("kpi_statuses.type IN ?", []string{"achieved", "exceeded"})
		case "at-risk":
			query = query.Where("kpi_statuses.type IN ?", []string{"near-target", "needs-improvement"})
		case "not-achieved":
			query = query.Where("kpi_statuses.is_not_achieved = ?", true)
		case "manual":
			query = query.Where("kpi_register.type = ?", "manual")

		case "automated":
			query = query.Where("kpi_register.type = ?", "automated")

		}
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&kpiRegister).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &kpiDTOs.KPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KPIRegisterRepository) GetByID(id uint) (*models.KPIRegister, error) {
	kpiRegister := &models.KPIRegister{}

	result := repo.db.
		Preload("CalculationType").
		First(kpiRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return kpiRegister, nil
}

func (repo *KPIRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kpiRegister *models.KPIRegister

	result := db.
		Preload("KPICategory").
		Preload("KPIStatus").
		Preload("Status").
		Preload("KPIReviewStatus").
		Preload("ResponsibleUser").
		Preload("CalculationType").
		Preload("ReviewFrequency").
		Preload("Unit").
		Preload("DataInputs").
		Preload("KPIAchievementLogic").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ReviewUser").
		First(&kpiRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return kpiRegister, nil
}

func (repo *KPIRegisterRepository) Create(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kpiRegister).Error; err != nil {
		return 0, err
	}

	return kpiRegister.ID, nil
}

func (repo *KPIRegisterRepository) Update(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kpiRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.KPIRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) CreateKPIRegisterInpuData(inputDatas []*models.KPIRegisterDataInput, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&inputDatas).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) DeleteKPIRegisterInputData(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("kpi_register_id = ?", kpiRegisterID).
		Delete(&models.KPIRegisterDataInput{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) GetTrendAnalysis(kpiID uint, limit int) ([]*models.KPIValueUpdate, error) {
	var kpiValueUpdate []*models.KPIValueUpdate

	result := repo.db.Where("kpi_register_id = ?", kpiID).
		Order("reviewed_date DESC").
		Limit(limit).
		Find(&kpiValueUpdate)

	if result.Error != nil {
		return nil, result.Error
	}

	return kpiValueUpdate, nil
}

func (repo *KPIRegisterRepository) GetLatestValueUpdate(kpiID uint) (*models.KPIValueUpdate, error) {
	var kpiValueUpdate *models.KPIValueUpdate

	result := repo.db.Preload("ReviewUser").
		Preload("Documents").
		Preload("Unit").
		Preload("DataInputs").
		Where("kpi_register_id = ? AND is_latest = 1", kpiID).
		First(&kpiValueUpdate)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when no record is found
		}
		return nil, result.Error
	}

	return kpiValueUpdate, nil
}
func (repo *KPIRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber *string
	err := repo.db.Model(&models.KPIRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Scan(&lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != nil && len(*lastReferenceNumber) >= 4 {
		lastDigits = (*lastReferenceNumber)[len(*lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "KPI", currentYear, lastReferenceID+1)
	return newReferenceNumber, nil
}
func (repo *KPIRegisterRepository) SetAllIsLatestFalse(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.KPIValueUpdate{}).
		Where("kpi_register_id = ? AND is_latest = 1", kpiRegisterID).
		Update("is_latest", false).Error
}
func (repo *KPIRegisterRepository) CreateKPIValueUpdate(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) (uint, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kpiValueUpdate).Error; err != nil {
		return 0, err
	}

	return kpiValueUpdate.ID, nil
}

func (repo *KPIRegisterRepository) GetKPICategorySummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPICategorySummaryDTO, error) {

	// ✅ Get all categories
	var categories []*kpiDTOs.KPICategoryDTO
	if err := repo.db.
		Table("kpi_categories").
		Select("id, title_en, title_ar, status").
		Where("status = ?", "active").
		Order("id").
		Scan(&categories).Error; err != nil {
		return nil, err
	}

	// ✅ Get all review statuses
	var reviewStatuses []*kpiDTOs.KPIReviewStatusDTO
	if err := repo.db.
		Table("kpi_review_statuses").
		Select("id, title_en, title_ar, type, color_code, status").
		Where("status = ?", "active").
		Order("id").
		Scan(&reviewStatuses).Error; err != nil {
		return nil, err
	}

	// ✅ Common filter logic
	applyFilters := func(query *gorm.DB) *gorm.DB {

		if params == nil {
			return query
		}

		// Normal filters
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}

		// Quick filters
		for _, filter := range params.QuickFilter {
			switch filter {
			case "my-kpis":
				query = query.Where("kr.responsible_user_id = ?", currentUserID)
			case "awaiting-update":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "pending":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "overdue":
				query = query.Where(`krs.type = ?`, "overdue")
			case "achieved":
				query = query.Where("ks.type IN ?", []string{"achieved", "exceeded"})
			case "at-risk":
				query = query.Where("ks.type IN ?", []string{"near-target", "needs-improvement"})
			case "not-achieved":
				query = query.Where("ks.is_not_achieved = ?", true)
			case "manual":
				query = query.Where("kr.type = ?", "manual")
			case "automated":
				query = query.Where("kr.type = ?", "automated")

			}
		}

		return query
	}

	// =========================
	// ✅ CATEGORY AGGREGATE
	// =========================

	type KPICategoryAggregateDTO struct {
		CategoryID uint    `gorm:"column:category_id"`
		Count      int64   `gorm:"column:count"`
		Percentage float64 `gorm:"column:percentage"`
	}

	var categoryAggregates []*KPICategoryAggregateDTO

	categoryAggregateQuery := repo.db.
		Table("kpi_register kr").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Select(`
			kr.kpi_category_id AS category_id,
			COUNT(kr.id) AS count,
			COALESCE(ROUND(AVG(CAST(kr.achievement AS FLOAT)), 2), 0) AS percentage
		`).
		Where("kr.kpi_category_id IS NOT NULL").
		Group("kr.kpi_category_id")

	categoryAggregateQuery = applyFilters(categoryAggregateQuery)

	if err := categoryAggregateQuery.Scan(&categoryAggregates).Error; err != nil {
		return nil, err
	}

	// =========================
	// ✅ REVIEW STATUS COUNTS
	// =========================

	var groupedCounts []*kpiDTOs.KPICategoryReviewCountDTO

	countQuery := repo.db.
		Table("kpi_register kr").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Select(`
			kr.kpi_category_id AS category_id,
			kr.kpi_review_status_id AS review_status_id,
			COUNT(kr.id) AS count
		`).
		Where("kr.kpi_category_id IS NOT NULL").
		Where("kr.kpi_review_status_id IS NOT NULL").
		Group("kr.kpi_category_id, kr.kpi_review_status_id")

	countQuery = applyFilters(countQuery)

	if err := countQuery.Scan(&groupedCounts).Error; err != nil {
		return nil, err
	}

	// =========================
	// ✅ MAP DATA
	// =========================

	categoryAggregateMap := make(map[uint]*KPICategoryAggregateDTO)
	for _, row := range categoryAggregates {
		categoryAggregateMap[row.CategoryID] = row
	}

	reviewCountMap := make(map[uint]map[uint]int64)
	for _, row := range groupedCounts {
		if _, ok := reviewCountMap[row.CategoryID]; !ok {
			reviewCountMap[row.CategoryID] = make(map[uint]int64)
		}
		reviewCountMap[row.CategoryID][row.ReviewStatusID] = row.Count
	}

	// =========================
	// ✅ BUILD RESPONSE
	// =========================

	results := make([]*kpiDTOs.KPICategorySummaryDTO, 0)

	for _, cat := range categories {

		var categoryCount int64
		var percentage float64

		if agg, ok := categoryAggregateMap[cat.ID]; ok {
			categoryCount = agg.Count
			percentage = agg.Percentage
		}

		// ✅ Skip empty categories
		if categoryCount == 0 {
			continue
		}

		reviewStatusItems := make([]*kpiDTOs.KPIReviewStatusSummaryDTO, 0)

		for _, rs := range reviewStatuses {

			var statusCount int64
			if catStatusMap, ok := reviewCountMap[cat.ID]; ok {
				statusCount = catStatusMap[rs.ID]
			}

			reviewStatusItems = append(reviewStatusItems, &kpiDTOs.KPIReviewStatusSummaryDTO{
				KPIReviewStatusDTO: kpiDTOs.KPIReviewStatusDTO{
					ID:        rs.ID,
					TitleEn:   rs.TitleEn,
					TitleAr:   rs.TitleAr,
					Type:      rs.Type,
					ColorCode: rs.ColorCode,
					Status:    rs.Status,
				},
				Count: statusCount,
			})
		}

		results = append(results, &kpiDTOs.KPICategorySummaryDTO{
			KPICategoryDTO: kpiDTOs.KPICategoryDTO{
				ID:      cat.ID,
				TitleEn: cat.TitleEn,
				TitleAr: cat.TitleAr,
				Status:  cat.Status,
			},
			Count:          categoryCount,
			Percentage:     percentage,
			ReviewStatuses: reviewStatusItems,
		})
	}

	// =========================
	// ✅ SORT (DESC)
	// =========================

	sort.Slice(results, func(i, j int) bool {
		if results[i].Percentage == results[j].Percentage {
			return results[i].Count > results[j].Count
		}
		return results[i].Percentage > results[j].Percentage
	})

	return results, nil
}

func (repo *KPIRegisterRepository) GetKPIStatusSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPIStatusSummaryDTO, error) {

	var statuses []*kpiDTOs.KPIStatusSummaryDTO

	// ✅ Get all statuses
	if err := repo.db.
		Table("kpi_statuses").
		Select(`
			id,
			title_en,
			title_ar,
			color_code,
			type,
			label,
			is_not_achieved,
			[order],
			status
		`).
		Where("status = ?", "active").
		Order("[order] ASC").
		Scan(&statuses).Error; err != nil {
		return nil, err
	}

	// ✅ Aggregate struct
	type KPIStatusAggregateDTO struct {
		KPIStatusID uint    `gorm:"column:kpi_status_id"`
		Count       int64   `gorm:"column:count"`
		Average     float64 `gorm:"column:average"`
	}

	var aggregates []*KPIStatusAggregateDTO

	query := repo.db.
		Table("kpi_register kr").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Select(`
			kr.kpi_status_id AS kpi_status_id,
			COUNT(kr.id) AS count,
			COALESCE(ROUND(AVG(CAST(kr.achievement AS FLOAT)), 2), 0) AS average
		`).
		Where("kr.kpi_status_id IS NOT NULL").
		Group("kr.kpi_status_id")

	// ✅ Apply Filters
	if params != nil {

		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}

		// ✅ Apply Quick Filters
		for _, filter := range params.QuickFilter {
			switch filter {
			case "my-kpis":
				query = query.Where("kr.responsible_user_id = ?", currentUserID)
			case "awaiting-update":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "pending":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "overdue":
				query = query.Where(`krs.type = ?`, "overdue")
			case "achieved":
				query = query.Where("ks.type IN ?", []string{"achieved", "exceeded"})
			case "at-risk":
				query = query.Where("ks.type IN ?", []string{"near-target", "needs-improvement"})
			case "not-achieved":
				query = query.Where("ks.is_not_achieved = ?", true)
			case "manual":
				query = query.Where("kr.type = ?", "manual")
			case "automated":
				query = query.Where("kr.type = ?", "automated")
			}
		}
	}

	if err := query.Scan(&aggregates).Error; err != nil {
		return nil, err
	}

	// ✅ Map aggregation
	aggregateMap := make(map[uint]*KPIStatusAggregateDTO, len(aggregates))
	for _, row := range aggregates {
		aggregateMap[row.KPIStatusID] = row
	}

	// ✅ Merge into statuses
	for _, status := range statuses {
		if agg, ok := aggregateMap[status.ID]; ok {
			status.Count = agg.Count
			status.Percentage = agg.Average
		} else {
			status.Count = 0
			status.Percentage = 0
		}
	}

	return statuses, nil
}

func (repo *KPIRegisterRepository) GetKPIOverallSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIOverallSummaryDTO, error) {

	type KPIOverallSummaryResult struct {
		TotalKPIs          int64   `gorm:"column:total_kpis"`
		OverallAchievement float64 `gorm:"column:overall_achievement"`
	}

	var result KPIOverallSummaryResult

	query := repo.db.
		Table("kpi_register kr").
		Joins("LEFT JOIN kpi_statuses ks ON ks.id = kr.kpi_status_id").
		Joins("LEFT JOIN kpi_review_statuses krs ON krs.id = kr.kpi_review_status_id").
		Select(`
		COUNT(kr.id) AS total_kpis,
		COALESCE(ROUND(AVG(CAST(kr.achievement AS FLOAT)), 2), 0) AS overall_achievement
	`)

	// ✅ Apply Filters
	if params != nil {
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("kr."+col+" IN ?", vals)
			}
		}

		// ✅ Apply Quick Filters (same as list)
		for _, filter := range params.QuickFilter {
			switch filter {
			case "my-kpis":
				query = query.Where("kr.responsible_user_id = ?", currentUserID)
			case "awaiting-update":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "pending":
				query = query.Where("krs.type = ?", "awaiting-update")
			case "overdue":
				query = query.Where(`krs.type = ?`, "overdue")
			case "achieved":
				query = query.Where("ks.type IN ?", []string{"achieved", "exceeded"})
			case "at-risk":
				query = query.Where("ks.type IN ?", []string{"near-target", "needs-improvement"})
			case "not-achieved":
				query = query.Where("ks.is_not_achieved = ?", true)
			case "manual":
				query = query.Where("kr.type = ?", "manual")
			case "automated":
				query = query.Where("kr.type = ?", "automated")
			}
		}
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}

	return &kpiDTOs.KPIOverallSummaryDTO{
		TotalKPIs:          result.TotalKPIs,
		OverallAchievement: result.OverallAchievement,
	}, nil
}

func (repo *KPIRegisterRepository) GetKPIByReviewStatus(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*kpiDTOs.KPIReviewStatusDistributionDTO, error) {

	var results []*kpiDTOs.KPIReviewStatusDistributionDTO

	query := repo.db.
		Table("kpi_review_statuses rs").
		Select(`
		rs.id,
		rs.title_en,
		rs.title_ar,
		rs.color_code,
		rs.type,
		COUNT(kr.id) AS count
	`).
		Joins(`
		LEFT JOIN kpi_register kr
			ON kr.kpi_review_status_id = rs.id
	`).
		Joins(`
		LEFT JOIN kpi_statuses ks
			ON ks.id = kr.kpi_status_id
	`).
		Where("rs.status = ?", "active")

	if params != nil {
		// ✅ Normal Filters
		for col, vals := range params.Filters {
			if len(vals) > 0 {
				query = query.Where("(kr."+col+" IN ? OR kr.id IS NULL)", vals)
			}
		}

		// ✅ Quick Filters
		for _, filter := range params.QuickFilter {
			switch filter {
			case "my-kpis":
				query = query.Where("(kr.responsible_user_id = ? OR kr.id IS NULL)", currentUserID)
			case "awaiting-update":
				query = query.Where("rs.type = ?", "awaiting-update")
			case "pending":
				query = query.Where("rs.type = ?", "awaiting-update")
			case "overdue":
				query = query.Where(`rs.type = ?`, "overdue")
			case "achieved":
				query = query.Where("(ks.type IN ? OR kr.id IS NULL)", []string{"achieved", "exceeded"})
			case "at-risk":
				query = query.Where("(ks.type IN ? OR kr.id IS NULL)", []string{"near-target", "needs-improvement"})
			case "not-achieved":
				query = query.Where("(ks.is_not_achieved = ? OR kr.id IS NULL)", true)
			case "manual":
				query = query.Where("(kr.type = ? OR kr.id IS NULL)", "manual")
			case "automated":
				query = query.Where("(kr.type = ? OR kr.id IS NULL)", "automated")
			}
		}
	}

	query = query.
		Group("rs.id, rs.title_en, rs.title_ar, rs.color_code, rs.type").
		Order("rs.id")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
func (repo *KPIRegisterRepository) CreateKPIValueUpdateInputData(inputDatas []*models.KPIValueUpdateDataInput, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&inputDatas).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) UpdateKPIValueUpdate(kpiValueUpdate *models.KPIValueUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kpiValueUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KPIRegisterRepository) DeleteKPIValueUpdateInputData(kpiValueUpdateID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("kpi_value_update_id = ?", kpiValueUpdateID).
		Delete(&models.KPIValueUpdateDataInput{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *KPIRegisterRepository) ActivateKPIRegisters(kpiRegisterIDs []uint, archiveStatusID uint, activeStatusID uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.KPIRegister{}).
		Where("id IN ? AND status_id = ?", kpiRegisterIDs, archiveStatusID).
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KPIRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status_id", activeStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *KPIRegisterRepository) ArchiveKPIRegisters(kpiRegisterIDs []uint, archiveStatusID uint, activeStatusID uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.KPIRegister{}).
		Where("id IN ? AND status_id = ?", kpiRegisterIDs, activeStatusID).
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KPIRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status_id", archiveStatusID).Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *KPIRegisterRepository) ForceDelete(ids []uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	queries := []string{

		`DELETE kpi_value_update_history_inputs
		 FROM kpi_value_update_history_inputs
		 JOIN kpi_value_update_histories ON kpi_value_update_histories.id = kpi_value_update_history_inputs.kpi_value_update_history_id
		 JOIN kpi_value_updates ON kpi_value_updates.id = kpi_value_update_histories.kpi_value_update_id
		 WHERE kpi_value_updates.kpi_register_id IN (?)`,

		`DELETE kpi_value_update_histories
		 FROM kpi_value_update_histories
		 JOIN kpi_value_updates ON kpi_value_updates.id = kpi_value_update_histories.kpi_value_update_id
		 WHERE kpi_value_updates.kpi_register_id IN (?)`,

		`DELETE kpi_value_update_documents
		 FROM kpi_value_update_documents
		 JOIN kpi_value_updates ON kpi_value_updates.id = kpi_value_update_documents.kpi_value_update_id
		 WHERE kpi_value_updates.kpi_register_id IN (?)`,

		`DELETE kpi_value_update_data_inputs
		 FROM kpi_value_update_data_inputs
		 JOIN kpi_value_updates ON kpi_value_updates.id = kpi_value_update_data_inputs.kpi_value_update_id
		 WHERE kpi_value_updates.kpi_register_id IN (?)`,

		`DELETE FROM strategy_profile_strategic_initiative_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM strategy_profile_kpis WHERE kpi_id IN (?)`,
		`DELETE FROM strategic_objective_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM meeting_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM kri_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_value_updates WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_third_parties WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_risks WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_register_data_inputs WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_projects WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_processes WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_controls WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_compliances WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_associated_documents WHERE kpi_register_id IN (?)`,
		`DELETE FROM kpi_assets WHERE kpi_register_id IN (?)`,
		`DELETE FROM incident_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM bcp_kpis WHERE kpi_register_id IN (?)`,
		`DELETE FROM audit_finding_kpis WHERE kpi_register_id IN (?)`,
	}

	for _, q := range queries {
		if err := db.Exec(q, ids).Error; err != nil {
			return err
		}
	}

	return db.Where("id IN ?", ids).Delete(&models.KPIRegister{}).Error
}

func (repo *KPIRegisterRepository) GetAllResponsibleUsers(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	var users []*userDTOs.UserListDTO
	var totalCount int64

	query := repo.db.Table("users").
		Select(`
			users.*,
			org1.title_en AS organization_unit_1_title_en,
			org2.title_en AS organization_unit_2_title_en,
			org3.title_en AS organization_unit_3_title_en,
			org4.title_en AS organization_unit_4_title_en,
			org5.title_en AS organization_unit_5_title_en,
			org1.title_ar AS organization_unit_1_title_ar,
			org2.title_ar AS organization_unit_2_title_ar,
			org3.title_ar AS organization_unit_3_title_ar,
			org4.title_ar AS organization_unit_4_title_ar,
			org5.title_ar AS organization_unit_5_title_ar,
			designations.title_en AS designation_title_en,
			designations.title_ar AS designation_title_ar,
			license_types.title_en as license_type_title_en,
			license_types.title_ar as license_type_title_ar,
			license_types.type as license_type_type,
			license_types.color_code as license_type_color_code,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.is_active AS user_status_is_active
		`).
		Joins("INNER JOIN organization_units org1 ON org1.id = users.organization_unit_1_id").
		Joins("LEFT JOIN organization_units org2 ON org2.id = users.organization_unit_2_id").
		Joins("LEFT JOIN organization_units org3 ON org3.id = users.organization_unit_3_id").
		Joins("LEFT JOIN organization_units org4 ON org4.id = users.organization_unit_4_id").
		Joins("LEFT JOIN organization_units org5 ON org5.id = users.organization_unit_5_id").
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins("LEFT JOIN license_types ON license_types.id = users.license_type_id").
		Joins("LEFT JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Where(`
			EXISTS (
				SELECT 1
				FROM kpi_register
				WHERE kpi_register.responsible_user_id = users.id
			)
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if params.Keyword != "" {
		query = query.Where(`
        (users.first_name_en LIKE ? OR
         users.last_name_en LIKE ? OR
         users.first_name_ar LIKE ? OR
         users.last_name_ar LIKE ? OR
         users.email LIKE ? OR
         users.mobile_number LIKE ? OR
         designations.title_en LIKE ? OR
         designations.title_ar LIKE ? OR
         user_statuses.title_en LIKE ? OR
         user_statuses.title_ar LIKE ?)
        `,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &userDTOs.UserListResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
