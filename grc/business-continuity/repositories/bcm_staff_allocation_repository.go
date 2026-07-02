package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BCMStaffAllocationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetFullDetailsByProcessID(processID uint) ([]*models.BCMStaffAllocation, error)
	Create(item *models.BCMStaffAllocation, tx ...utility.DatabaseTransaction) (uint, error)
	CreateOutageBucket(item *models.BCMStaffAllocationOutageBucket, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteUsersByOutageBucketID(outage_bucket_id uint, tx ...utility.DatabaseTransaction) error
	DeleteUsersByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	DeleteOutageBucketsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	CreateUser(user *models.BCMStaffAllocationOutageBucketUser, tx ...utility.DatabaseTransaction) error
	DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	GetRatingByProcessID(processID uint) (*models.BusinessImpactAssessmentRating, error)
	GetOutageBuckets() ([]*models.BCMBIAOutageBucket, error)
	GetAnalysisByProcessID(processID uint) ([]*models.ProcessEmployeeAnalysis, error)
}

type BCMStaffAllocationRepository struct {
	db *gorm.DB
}

func NewBCMStaffAllocationRepository(db *gorm.DB) *BCMStaffAllocationRepository {
	return &BCMStaffAllocationRepository{db: db}
}

func (repo *BCMStaffAllocationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BCMStaffAllocationRepository) GetFullDetailsByProcessID(processID uint) ([]*models.BCMStaffAllocation, error) {
	var allocations []*models.BCMStaffAllocation
	err := repo.db.
		Preload("ProcessEmployeeAnalysis").
		Preload("ProcessEmployeeAnalysis.Designation").
		Preload("ProcessEmployeeAnalysis.Employees").
		Preload("ProcessEmployeeAnalysis.Employees.User").
		Preload("ProcessEmployeeAnalysis.Employees.Tools").
		Preload("ProcessEmployeeAnalysis.Employees.Tools.EquipmentTool").
		Preload("OutageBuckets").
		Preload("OutageBuckets.OutageBucket").
		Preload("OutageBuckets.Users").
		Preload("OutageBuckets.Users.User").
		Where("process_id = ?", processID).
		First(&allocations).Error
	return allocations, err
}

func (repo *BCMStaffAllocationRepository) CreateOutageBucket(item *models.BCMStaffAllocationOutageBucket, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(item).Error
	return item.ID, err
}

func (repo *BCMStaffAllocationRepository) DeleteUsersByOutageBucketID(outageBucketRowID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("bcm_staff_allocation_outage_bucket_id = ?", outageBucketRowID).
		Delete(&models.BCMStaffAllocationOutageBucketUser{}).Error
}

func (repo *BCMStaffAllocationRepository) CreateUser(user *models.BCMStaffAllocationOutageBucketUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(user).Error
}

func (repo *BCMStaffAllocationRepository) Create(item *models.BCMStaffAllocation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(item).Error
	return item.ID, err
}

func (repo *BCMStaffAllocationRepository) DeleteUsersByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("bcm_staff_allocation_outage_bucket_id IN (SELECT id FROM bcm_staff_allocation_outage_buckets WHERE bcm_staff_allocation_id IN (SELECT id FROM bcm_staff_allocations WHERE process_id = ?))", processID).Delete(&models.BCMStaffAllocationOutageBucketUser{}).Error
}

func (repo *BCMStaffAllocationRepository) DeleteOutageBucketsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("bcm_staff_allocation_id IN (SELECT id FROM bcm_staff_allocations WHERE process_id = ?)", processID).Delete(&models.BCMStaffAllocationOutageBucket{}).Error
}

func (repo *BCMStaffAllocationRepository) DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	if err := repo.DeleteUsersByProcessID(processID, tx...); err != nil {
		return err
	}
	if err := repo.DeleteOutageBucketsByProcessID(processID, tx...); err != nil {
		return err
	}
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_id = ?", processID).Delete(&models.BCMStaffAllocation{}).Error
}

func (repo *BCMStaffAllocationRepository) GetRatingByProcessID(processID uint) (*models.BusinessImpactAssessmentRating, error) {
	var rating models.BusinessImpactAssessmentRating
	err := repo.db.Where("process_id = ?", processID).First(&rating).Error
	return &rating, err
}

func (repo *BCMStaffAllocationRepository) GetOutageBuckets() ([]*models.BCMBIAOutageBucket, error) {
	var buckets []*models.BCMBIAOutageBucket
	err := repo.db.Where("status = 'active'").Order("[order] ASC").Find(&buckets).Error
	return buckets, err
}

func (repo *BCMStaffAllocationRepository) GetAnalysisByProcessID(processID uint) ([]*models.ProcessEmployeeAnalysis, error) {
	var analyses []*models.ProcessEmployeeAnalysis
	err := repo.db.
		Where("process_id = ?", processID).
		Order("id ASC").
		Find(&analyses).Error
	return analyses, err
}
