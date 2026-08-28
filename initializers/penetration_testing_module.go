package initializers

import (
	"gorm.io/gorm"
	"jethur/grc/penetration-testing/controllers"
	"jethur/grc/penetration-testing/repositories"
	"jethur/grc/penetration-testing/services"
)

type PenetrationTestingModule struct {
	PenetrationTestTypeController                 *controllers.PenetrationTestTypeController
	PenetrationTestPlanController                 *controllers.PenetrationTestPlanController
	PenetrationTestExecutorController             *controllers.PenetrationTestExecutorController
	PenetrationTestEnvironmentReadinessController *controllers.PenetrationTestEnvironmentReadinessController
	PenetrationTestClassificationController       *controllers.PenetrationTestClassificationController
	PenetrationTestExecutionStatusController      *controllers.PenetrationTestExecutionStatusController
	PenetrationTestSecurityResultController       *controllers.PenetrationTestSecurityResultController
}

func NewPenetrationTestingModule(db *gorm.DB, sharedModule *SharedModule) *PenetrationTestingModule {
	typeService := services.NewPenetrationTestTypeService(repositories.NewPenetrationTestTypeRepository(db), sharedModule.ActivityLogService)
	planService := services.NewPenetrationTestPlanService(repositories.NewPenetrationTestPlanRepository(db), sharedModule.ActivityLogService)
	executorService := services.NewPenetrationTestExecutorService(repositories.NewPenetrationTestExecutorRepository(db), sharedModule.ActivityLogService)
	readinessService := services.NewPenetrationTestEnvironmentReadinessService(repositories.NewPenetrationTestEnvironmentReadinessRepository(db), sharedModule.ActivityLogService)
	classificationService := services.NewPenetrationTestClassificationService(repositories.NewPenetrationTestClassificationRepository(db))
	executionStatusService := services.NewPenetrationTestExecutionStatusService(repositories.NewPenetrationTestExecutionStatusRepository(db))
	securityResultService := services.NewPenetrationTestSecurityResultService(repositories.NewPenetrationTestSecurityResultRepository(db))
	return &PenetrationTestingModule{
		PenetrationTestTypeController:                 controllers.NewPenetrationTestTypeController(typeService),
		PenetrationTestPlanController:                 controllers.NewPenetrationTestPlanController(planService),
		PenetrationTestExecutorController:             controllers.NewPenetrationTestExecutorController(executorService),
		PenetrationTestEnvironmentReadinessController: controllers.NewPenetrationTestEnvironmentReadinessController(readinessService),
		PenetrationTestClassificationController:       controllers.NewPenetrationTestClassificationController(classificationService),
		PenetrationTestExecutionStatusController:      controllers.NewPenetrationTestExecutionStatusController(executionStatusService),
		PenetrationTestSecurityResultController:       controllers.NewPenetrationTestSecurityResultController(securityResultService),
	}
}
