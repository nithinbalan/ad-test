package services

import (
	"jethur/grc/common/repositories"
	settingService "jethur/grc/settings/services"
)

type ActionPlanDelegateeUserserviceInterface interface {
	AssignDelegateeIfPermitted(responsibleUserID uint, moduleIdentifier string) (*uint, error)
}

type ActionPlanDelegateeUserservice struct {
	repo      repositories.ActionPlanDelegateeUserRepositoryInterface
	moduleSvc settingService.ModuleServiceInterface
}

func NewActionPlanDelegateeUserservice(
	repo repositories.ActionPlanDelegateeUserRepositoryInterface,
	moduleSvc settingService.ModuleServiceInterface,
) ActionPlanDelegateeUserserviceInterface {
	return &ActionPlanDelegateeUserservice{
		repo:      repo,
		moduleSvc: moduleSvc,
	}
}

func (taps *ActionPlanDelegateeUserservice) AssignDelegateeIfPermitted(responsibleUserID uint, moduleIdentifier string) (*uint, error) {
	// Step 1: Get Responsible User
	responsibleUser, err := taps.repo.GetUserByID(responsibleUserID)
	if err != nil {
		return nil, err
	}

	// Step 2: Get Module ID
	moduleID, err := taps.moduleSvc.GetModuleIDByIdentifier(moduleIdentifier)
	if err != nil {
		return nil, err
	}
	if responsibleUser.DelegateeUserID != nil {
		// Step 3: Check permission
		hasPermission, err := taps.repo.HasViewPermission(*responsibleUser.DelegateeUserID, moduleID)
		if err != nil {
			return nil, err
		}
		if hasPermission {
			return responsibleUser.DelegateeUserID, nil
		}
	}

	return nil, nil
}
