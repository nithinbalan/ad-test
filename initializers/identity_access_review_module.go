package initializers

import (
	"jethur/grc/identity-access-review/controllers"
	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/identity-access-review/services"

	"gorm.io/gorm"
)

type IdentityAccessReviewModule struct {
	ReviewTypeController             *controllers.ReviewTypeController
	ReviewCampaignStatusController   *controllers.ReviewCampaignStatusController
	IdentityReviewCampaignController *controllers.IdentityReviewCampaignController
	IdentityReviewUserController     *controllers.IdentityReviewUserController
	PrivilegedAccountController      *controllers.PrivilegedAccountController
	AccountTypeController            *controllers.AccountTypeController
}

func NewIdentityAccessReviewModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	organizationModule *OrganizationModule,
	settingsModule *SettingsModule,
) *IdentityAccessReviewModule {
	reviewTypeRepo := repositories.NewReviewTypeRepository(db)
	reviewTypeSvc := services.NewReviewTypeService(reviewTypeRepo)
	reviewTypeController := controllers.NewReviewTypeController(reviewTypeSvc)

	reviewCampaignStatusRepo := repositories.NewReviewCampaignStatusRepository(db)
	reviewCampaignStatusSvc := services.NewReviewCampaignStatusService(reviewCampaignStatusRepo)
	reviewCampaignStatusController := controllers.NewReviewCampaignStatusController(reviewCampaignStatusSvc)

	identityReviewCampaignRepo := repositories.NewIdentityReviewCampaignRepository(db)
	identityReviewCampaignSvc := services.NewIdentityReviewCampaignService(
		identityReviewCampaignRepo,
		sharedModule.ActivityLogService,
		reviewTypeSvc,
		reviewCampaignStatusSvc,
		organizationModule.UserSvc,
	)
	identityReviewCampaignController := controllers.NewIdentityReviewCampaignController(identityReviewCampaignSvc)

	privilegedAccountRepo := repositories.NewPrivilegedAccountRepository(db)
	privilegedAccountSvc := services.NewPrivilegedAccountService(
		privilegedAccountRepo,
		sharedModule.ActivityLogService,
	)
	privilegedAccountController := controllers.NewPrivilegedAccountController(privilegedAccountSvc)

	accountTypeRepo := repositories.NewAccountTypeRepository(db)
	accountTypeSvc := services.NewAccountTypeService(accountTypeRepo)
	accountTypeController := controllers.NewAccountTypeController(accountTypeSvc)

	identityReviewUserRepo := repositories.NewIdentityReviewUserRepository(db)
	identityReviewUserSvc := services.NewIdentityReviewUserService(
		identityReviewUserRepo,
		sharedModule.ActivityLogService,
		settingsModule.OrganizationUnitSvc,
		privilegedAccountSvc,
		accountTypeSvc,
	)
	identityReviewUserController := controllers.NewIdentityReviewUserController(identityReviewUserSvc)

	return &IdentityAccessReviewModule{
		ReviewTypeController:             reviewTypeController,
		ReviewCampaignStatusController:   reviewCampaignStatusController,
		IdentityReviewCampaignController: identityReviewCampaignController,
		IdentityReviewUserController:     identityReviewUserController,
		PrivilegedAccountController:      privilegedAccountController,
		AccountTypeController:            accountTypeController,
	}
}
