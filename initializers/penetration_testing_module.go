package initializers

import (
	"gorm.io/gorm"
	"jethur/grc/penetration-testing/controllers"
	"jethur/grc/penetration-testing/repositories"
	"jethur/grc/penetration-testing/services"
)

type PenetrationTestingModule struct {
	PenetrationTestController                       *controllers.PenetrationTestController
	PenetrationTestEnvironmentPreparationController *controllers.PenetrationTestEnvironmentPreparationController
	PenetrationTestExecutionController              *controllers.PenetrationTestExecutionController
	PenetrationTestManagementSharingController      *controllers.PenetrationTestManagementSharingController
	PenetrationTestFindingController                *controllers.PenetrationTestFindingController
	PenetrationTestFindingStatusUpdateController    *controllers.PenetrationTestFindingStatusUpdateController
	PenetrationTestRetestController                 *controllers.PenetrationTestRetestController
	PenetrationTestTypeController                   *controllers.PenetrationTestTypeController
	PenetrationTestPlanController                   *controllers.PenetrationTestPlanController
	PenetrationTestExecutorController               *controllers.PenetrationTestExecutorController
	PenetrationTestEnvironmentReadinessController   *controllers.PenetrationTestEnvironmentReadinessController
	PenetrationTestClassificationController         *controllers.PenetrationTestClassificationController
	PenetrationTestExecutionStatusController        *controllers.PenetrationTestExecutionStatusController
	PenetrationTestSecurityResultController         *controllers.PenetrationTestSecurityResultController
	PenetrationTestSupportingDocumentController     *controllers.PenetrationTestSupportingDocumentController
}

func NewPenetrationTestingModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
) *PenetrationTestingModule {
	typeService := services.NewPenetrationTestTypeService(repositories.NewPenetrationTestTypeRepository(db), sharedModule.ActivityLogService)
	planService := services.NewPenetrationTestPlanService(repositories.NewPenetrationTestPlanRepository(db), sharedModule.ActivityLogService)
	executorService := services.NewPenetrationTestExecutorService(repositories.NewPenetrationTestExecutorRepository(db), sharedModule.ActivityLogService)
	readinessService := services.NewPenetrationTestEnvironmentReadinessService(repositories.NewPenetrationTestEnvironmentReadinessRepository(db), sharedModule.ActivityLogService)
	classificationService := services.NewPenetrationTestClassificationService(repositories.NewPenetrationTestClassificationRepository(db))
	executionStatusService := services.NewPenetrationTestExecutionStatusService(repositories.NewPenetrationTestExecutionStatusRepository(db))
	securityResultService := services.NewPenetrationTestSecurityResultService(repositories.NewPenetrationTestSecurityResultRepository(db))
	penetrationTestService := services.NewPenetrationTestService(
		repositories.NewPenetrationTestRepository(db),
		sharedModule.ActivityLogService,
		typeService,
		planService,
		classificationService,
		executorService,
		organizationModule.UserSvc,
	)
	environmentPreparationService := services.NewPenetrationTestEnvironmentPreparationService(
		repositories.NewPenetrationTestEnvironmentPreparationRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestService,
		readinessService,
	)
	executionService := services.NewPenetrationTestExecutionService(
		repositories.NewPenetrationTestExecutionRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestService,
		executionStatusService,
		securityResultService,
	)
	managementSharingService := services.NewPenetrationTestManagementSharingService(
		repositories.NewPenetrationTestManagementSharingRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestService,
	)
	penetrationTestFindingService := services.NewPenetrationTestFindingService(
		repositories.NewPenetrationTestFindingRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestService,
		masterModule.SeverityLevelSvc,
		masterModule.ActionPlanStatusSvc,
		organizationModule.UserSvc,
	)
	penetrationTestFindingStatusUpdateService := services.NewPenetrationTestFindingStatusUpdateService(
		repositories.NewPenetrationTestFindingStatusUpdateRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestFindingService,
		masterModule.ActionPlanStatusSvc,
	)
	penetrationTestRetestService := services.NewPenetrationTestRetestService(
		repositories.NewPenetrationTestRetestRepository(db),
		sharedModule.ActivityLogService,
		penetrationTestService,
		penetrationTestFindingService,
		executorService,
		masterModule.ActionPlanStatusSvc,
		repositories.NewPenetrationTestFindingStatusUpdateRepository(db),
	)
	supportingDocumentService := services.NewPenetrationTestSupportingDocumentService(
		repositories.NewPenetrationTestSupportingDocumentRepository(db),
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
		penetrationTestService,
	)
	return &PenetrationTestingModule{
		PenetrationTestController:                       controllers.NewPenetrationTestController(penetrationTestService),
		PenetrationTestEnvironmentPreparationController: controllers.NewPenetrationTestEnvironmentPreparationController(environmentPreparationService),
		PenetrationTestExecutionController:              controllers.NewPenetrationTestExecutionController(executionService),
		PenetrationTestManagementSharingController:      controllers.NewPenetrationTestManagementSharingController(managementSharingService),
		PenetrationTestFindingController:                controllers.NewPenetrationTestFindingController(penetrationTestFindingService),
		PenetrationTestFindingStatusUpdateController:    controllers.NewPenetrationTestFindingStatusUpdateController(penetrationTestFindingStatusUpdateService),
		PenetrationTestRetestController:                 controllers.NewPenetrationTestRetestController(penetrationTestRetestService),
		PenetrationTestTypeController:                   controllers.NewPenetrationTestTypeController(typeService),
		PenetrationTestPlanController:                   controllers.NewPenetrationTestPlanController(planService),
		PenetrationTestExecutorController:               controllers.NewPenetrationTestExecutorController(executorService),
		PenetrationTestEnvironmentReadinessController:   controllers.NewPenetrationTestEnvironmentReadinessController(readinessService),
		PenetrationTestClassificationController:         controllers.NewPenetrationTestClassificationController(classificationService),
		PenetrationTestExecutionStatusController:        controllers.NewPenetrationTestExecutionStatusController(executionStatusService),
		PenetrationTestSecurityResultController:         controllers.NewPenetrationTestSecurityResultController(securityResultService),
		PenetrationTestSupportingDocumentController:     controllers.NewPenetrationTestSupportingDocumentController(supportingDocumentService, db),
	}
}
