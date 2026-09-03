package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestAndExerciseTestResultRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(testAndExerciseID uint) (*models.TestAndExerciseTestResult, error)
	GetFullDetailsByID(testAndExerciseID uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseTestResult, error)
	Create(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) (uint, error)
	Update(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	UpdateTestAndExerciseProcess(process *models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error
	UpdateTestAndExerciseSystemRegister(systemRegister *models.TestAndExerciseSystemRegister, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.TestAndExerciseTestResultDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.TestAndExerciseTestResultDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateTestResultSystemRegisters(items []*models.TestAndExerciseTestResultSystemRegister, tx ...utility.DatabaseTransaction) error
	GetTestResultSystemRegisterByID(id uint) (*models.TestAndExerciseTestResultSystemRegister, error)
	UpdateTestResultSystemRegister(item *models.TestAndExerciseTestResultSystemRegister, tx ...utility.DatabaseTransaction) error
	AreAllTestResultSystemRegistersUpdated(testResultID uint, tx ...utility.DatabaseTransaction) (bool, error)
	UpdateTestAndExercise(testAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) error

	CreatePerformanceMeasure(measure *models.TestResultPerformanceMeasure, tx ...utility.DatabaseTransaction) (uint, error)
	GetPerformanceMeasureByID(id uint) (*models.TestResultPerformanceMeasure, error)
	UpdatePerformanceMeasure(measure *models.TestResultPerformanceMeasure, tx ...utility.DatabaseTransaction) error
	DeletePerformanceMeasure(id uint, tx ...utility.DatabaseTransaction) error
	GetPerformanceMeasuresByTestResultID(testResultID uint) ([]*models.TestResultPerformanceMeasure, error)
}

type TestAndExerciseTestResultRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseTestResultRepository(db *gorm.DB) *TestAndExerciseTestResultRepository {
	return &TestAndExerciseTestResultRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TestAndExerciseTestResultRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TestAndExerciseTestResultRepository) GetByID(testAndExerciseID uint) (*models.TestAndExerciseTestResult, error) {
	var testAndExerciseTestResult *models.TestAndExerciseTestResult

	result := repo.db.Where("test_and_exercise_id = ?", testAndExerciseID).First(&testAndExerciseTestResult)
	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseTestResult, nil

}

func (repo *TestAndExerciseTestResultRepository) GetFullDetailsByID(testAndExerciseID uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseTestResult, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var testAndExerciseTestResult *models.TestAndExerciseTestResult

	result := db.Where("test_and_exercise_id = ?", testAndExerciseID).
		Preload("TestAndExercise").
		Preload("TestAndExercise.BusinessContinuityPlan").
		Preload("TestAndExercise.BusinessContinuityPlan.BusinessFunction").
		Preload("TestAndExercise.TestAndExerciseStatus").
		Preload("TestAndExercise.BcpTestAndExercise").
		Preload("TestAndExercise.BcpTestAndExercise.ResponsibleUser").
		Preload("TestAndExercise.BcpTestAndExercise.ReviewFrequency").
		Preload("TestAndExercise.TestAndExerciseProcesses").
		Preload("TestAndExercise.TestAndExerciseProcesses.Process").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedRtoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedRpoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualRtoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualRpoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.Timeliness").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.SystemRegister").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.PlannedRtoUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.PlannedMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.PlannedRpoUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.ActualRtoUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.ActualMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.ActualRpoUnit").
		Preload("TestAndExercise.TestAndExerciseSystemRegisters.Timeliness").
		Preload("Documents").
		Preload("TestResultSystemRegisters").
		Preload("TestResultSystemRegisters.SystemRegister").
		Preload("TestResultSystemRegisters.CommunicationFeedbackOption").
		Preload("TestResultSystemRegisters.RecoveryFeedbackOption").
		Preload("TestResultSystemRegisters.TestResultOption").
		Preload("PerformanceMeasures").
		Preload("CreatedByUser").
		First(&testAndExerciseTestResult)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseTestResult, nil
}

func (repo *TestAndExerciseTestResultRepository) Create(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(TestAndExerciseTestResult).Error; err != nil {
		return 0, err
	}

	return TestAndExerciseTestResult.ID, nil
}

func (repo *TestAndExerciseTestResultRepository) Update(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(TestAndExerciseTestResult).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseTestResultRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business test restul
	if err := db.Delete(&models.TestAndExerciseTestResult{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseTestResultRepository) UpdateTestAndExerciseProcess(process *models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(process).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) UpdateTestAndExerciseSystemRegister(systemRegister *models.TestAndExerciseSystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(systemRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) CreateDocuments(docs []*models.TestAndExerciseTestResultDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) GetDocumentByID(id uint) (*models.TestAndExerciseTestResultDocument, error) {
	var document models.TestAndExerciseTestResultDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *TestAndExerciseTestResultRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.TestAndExerciseTestResultDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseTestResultRepository) CreateTestResultSystemRegisters(items []*models.TestAndExerciseTestResultSystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&items).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) GetTestResultSystemRegisterByID(id uint) (*models.TestAndExerciseTestResultSystemRegister, error) {
	var item models.TestAndExerciseTestResultSystemRegister
	if err := repo.db.Preload("TestAndExerciseTestResult").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *TestAndExerciseTestResultRepository) UpdateTestResultSystemRegister(item *models.TestAndExerciseTestResultSystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(item).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) AreAllTestResultSystemRegistersUpdated(testResultID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var totalCount, updatedCount int64
	if err := db.Model(&models.TestAndExerciseTestResultSystemRegister{}).
		Where("test_and_exercise_test_result_id = ?", testResultID).
		Count(&totalCount).Error; err != nil {
		return false, err
	}

	if err := db.Model(&models.TestAndExerciseTestResultSystemRegister{}).
		Where("test_and_exercise_test_result_id = ? AND test_result_option_id IS NOT NULL", testResultID).
		Count(&updatedCount).Error; err != nil {
		return false, err
	}

	return totalCount > 0 && totalCount == updatedCount, nil
}

func (repo *TestAndExerciseTestResultRepository) UpdateTestAndExercise(testAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(testAndExercise).Error
}

func (repo *TestAndExerciseTestResultRepository) CreatePerformanceMeasure(measure *models.TestResultPerformanceMeasure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(measure).Error; err != nil {
		return 0, err
	}

	return measure.ID, nil
}

func (repo *TestAndExerciseTestResultRepository) GetPerformanceMeasureByID(id uint) (*models.TestResultPerformanceMeasure, error) {
	var measure *models.TestResultPerformanceMeasure

	result := repo.db.First(&measure, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return measure, nil
}

func (repo *TestAndExerciseTestResultRepository) UpdatePerformanceMeasure(measure *models.TestResultPerformanceMeasure, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(measure).Error
}

func (repo *TestAndExerciseTestResultRepository) DeletePerformanceMeasure(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.TestResultPerformanceMeasure{}, id).Error
}

func (repo *TestAndExerciseTestResultRepository) GetPerformanceMeasuresByTestResultID(testResultID uint) ([]*models.TestResultPerformanceMeasure, error) {
	var measures []*models.TestResultPerformanceMeasure

	result := repo.db.Where("test_result_id = ?", testResultID).Find(&measures)
	if result.Error != nil {
		return nil, result.Error
	}

	return measures, nil
}
