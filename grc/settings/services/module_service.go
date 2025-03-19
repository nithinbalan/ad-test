package services

import (
	"errors"
	"fmt"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ModuleServiceInterface interface {
	GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error)
	GetModuleIDByIdentifier(identifier string) (uint, error)
	GetByID(id uint) (*models.Module, error)
	GetFullDetailsByID(id uint) (*dtos.ModuleDetailDTO, error)
	EnableModule(id uint) (*models.Module, error)
	DisableModule(id uint) (*models.Module, error)
	ReorderModule(modules []dtos.ModuleReorderDTO) error
	ReorderSubModule(modules []dtos.ModuleReorderDTO) error
	UpdateModule(id uint, moduleFormDTO dtos.ModuleFormDTO) error
	UpdateSubModule(id uint, subModuleFormDTO dtos.SubModuleFormDTO) error
}

type ModuleService struct {
	repo repositories.ModuleRepositoryInterface
}

func NewModuleService(repo repositories.ModuleRepositoryInterface) ModuleServiceInterface {
	return &ModuleService{
		repo: repo,
	}
}

func (ms *ModuleService) GetModuleIDByIdentifier(identifier string) (uint, error) {
	return ms.repo.GetModuleIDByIdentifier(identifier)
}

func (ms *ModuleService) GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error) {
	return ms.repo.GetAll(params)
}

func (ms *ModuleService) GetByID(id uint) (*models.Module, error) {
	module, err := ms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return module, nil
}

func (ms *ModuleService) GetByIdentifier(identifier string) (*models.Module, error) {
	module, err := ms.repo.GetByIdentifier(identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, identifier)
		}
		return nil, err
	}
	return module, nil
}

func (ms ModuleService) GetFullDetailsByID(moduleId uint) (*dtos.ModuleDetailDTO, error) {
	// Initialize your DTO
	var moduleDTO dtos.ModuleDetailDTO

	module, err := ms.repo.GetFullDetailsByID(moduleId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, moduleId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(module, &moduleDTO)
	if err != nil {
		return nil, err
	}

	return &moduleDTO, nil
}

func (ms *ModuleService) EnableModule(id uint) (*models.Module, error) {
	existingSystemPreference, err := ms.repo.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	// Extract allowed modules from the license key
	claims, err := utility.DecryptJWTToken(*existingSystemPreference.LicenseKey)
	if err != nil {
		return nil, err
	}

	// Extract allowed modules (can be a string or array)
	rawModules, ok := claims["enabled_modules"]
	if !ok {
		return nil, fmt.Errorf("invalid or missing enabled modules in license key")
	}

	allowedModules, ok := rawModules.([]string)
	if !ok {
		return nil, utility.NewValidationError(fmt.Errorf("invalid enabled_modules format"), "invalid", "enabled modules")
	}

	// var allowedModules []string

	// if str, ok := rawModules.(string); ok {
	// 	if err := json.Unmarshal([]byte(str), &allowedModules); err != nil {
	// 		fmt.Println("1")
	// 		return nil, utility.NewValidationError(fmt.Errorf("failed to parse enabled_modules"), "invalid", "enabled modules")
	// 	}
	// } else {
	// 	fmt.Println("2")
	// 	return nil, utility.NewValidationError(fmt.Errorf("invalid enabled_modules format"), "invalid", "enabled modules")
	// }

	moduleDetail, err := ms.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}
	fmt.Println(allowedModules)
	// Check if the module identifier is allowed in the license
	isAllowed := false
	for _, allowedModule := range allowedModules {
		if allowedModule == moduleDetail.Identifier {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return nil, utility.NewValidationError(fmt.Errorf("invalid module"), "module_not_allowed", moduleDetail.TitleEn)
	}

	// Enable the module
	module, err := ms.GetByID(id)
	if err != nil {
		return nil, err
	}
	module.IsEnable = true
	_, err = ms.repo.Update(module.ID, module)
	if err != nil {
		return nil, err
	}

	// Enable submodules if it's a parent module
	if module.ModuleID == nil {
		if len(moduleDetail.SubModules) > 0 {

			// enabling all sub modules
			for _, subModule := range moduleDetail.SubModules {
				existingSubModule, err := ms.GetByID(subModule.ID)

				if err != nil {
					return nil, err
				}
				existingSubModule.IsEnable = true
				_, err = ms.repo.Update(existingSubModule.ID, existingSubModule)
				if err != nil {
					return nil, err
				}
			}
		}

	} else {
		// enabling main module
		existingModule, err := ms.GetByID(*module.ModuleID)

		if err != nil {
			return nil, err
		}
		existingModule.IsEnable = true
		_, err = ms.repo.Update(existingModule.ID, existingModule)
		if err != nil {
			return nil, err
		}

	}
	return module, nil
}

func (ms *ModuleService) DisableModule(id uint) (*models.Module, error) {
	moduleDetail, err := ms.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	module, err := ms.GetByID(id)
	if err != nil {
		return nil, err
	}

	module.IsEnable = false
	_, err = ms.repo.Update(module.ID, module)
	if err != nil {
		return nil, err
	}

	if module.ModuleID == nil {
		if len(moduleDetail.SubModules) > 0 {
			// disabling all sub modules
			for _, subModule := range moduleDetail.SubModules {
				existingSubModule, err := ms.GetByID(subModule.ID)

				if err != nil {
					return nil, err
				}
				existingSubModule.IsEnable = false
				_, err = ms.repo.Update(module.ID, existingSubModule)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		mainModuleDetail, err := ms.GetFullDetailsByID(*module.ModuleID)
		if err != nil {
			return nil, err
		}
		// disabling main module
		allDisabled := true
		for _, subModule := range mainModuleDetail.SubModules {
			if subModule.ID != module.ID && subModule.IsEnable {
				allDisabled = false
				break
			}
		}

		// Disable parent module if all sub-modules are disabled
		if allDisabled {
			existingModule, err := ms.GetByID(*module.ModuleID)

			if err != nil {
				return nil, err
			}
			existingModule.IsEnable = false
			_, err = ms.repo.Update(existingModule.ID, existingModule)
			if err != nil {
				return nil, err
			}
		}
	}
	return module, nil
}

func (ms *ModuleService) ReorderModule(modules []dtos.ModuleReorderDTO) error {
	for _, module := range modules {
		existingModule, err := ms.GetByIdentifier(module.Identifier)

		if err != nil {
			return nil
		}
		existingModule.Order = module.Order
		_, err = ms.repo.Update(existingModule.ID, existingModule)
		if err != nil {
			return nil
		}

	}

	return nil
}

func (ms *ModuleService) ReorderSubModule(modules []dtos.ModuleReorderDTO) error {
	for _, module := range modules {
		existingModule, err := ms.GetByIdentifier(module.Identifier)

		if err != nil {
			return nil
		}
		existingModule.Order = module.Order
		_, err = ms.repo.Update(existingModule.ID, existingModule)
		if err != nil {
			return nil
		}

	}

	return nil
}

func (ms *ModuleService) UpdateModule(id uint, moduleFormDTO dtos.ModuleFormDTO) error {
	existingModule, err := ms.GetByID(id)

	if err != nil {
		return nil
	}
	existingModule.TitleEn = moduleFormDTO.TitleEn
	existingModule.TitleAr = moduleFormDTO.TitleAr
	existingModule.Icon = moduleFormDTO.Icon
	_, err = ms.repo.Update(existingModule.ID, existingModule)
	if err != nil {
		return nil
	}
	return nil
}

func (ms *ModuleService) UpdateSubModule(id uint, subModuleFormDTO dtos.SubModuleFormDTO) error {
	existingModule, err := ms.GetByID(id)

	if err != nil {
		return nil
	}
	existingModule.TitleEn = subModuleFormDTO.TitleEn
	existingModule.TitleAr = subModuleFormDTO.TitleAr
	_, err = ms.repo.Update(existingModule.ID, existingModule)
	if err != nil {
		return nil
	}
	return nil
}
