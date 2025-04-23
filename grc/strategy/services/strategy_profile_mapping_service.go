package services

import (
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	StrategyProfileMappingRepositories "jethur/grc/strategy/repositories"
)

type StrategyProfileMappingServiceInterface interface {
	GetStrategyProfile(profileId uint) (*strategyDTOs.StrategyProfileResponseDTO, error)
}

// Service and implementations
type StrategyProfileMappingService struct {
	repo                      StrategyProfileMappingRepositories.StrategyProfileMappingRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	userSettingSvc            organizationServices.UserSettingServiceInterface
	organizationUnitSvc       settingServices.OrganizationUnitServiceInterface
	strategyProfileSettingSvc StrategyProfileSettingServiceInterface
}

func NewStrategyProfileMappingService(
	repo StrategyProfileMappingRepositories.StrategyProfileMappingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	strategyProfileSettingSvc StrategyProfileSettingServiceInterface,
) StrategyProfileMappingServiceInterface {
	return &StrategyProfileMappingService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		userSettingSvc:            userSettingSvc,
		organizationUnitSvc:       organizationUnitSvc,
		strategyProfileSettingSvc: strategyProfileSettingSvc,
	}
}

func (spms *StrategyProfileMappingService) GetStrategyProfile(profileID uint) (*strategyDTOs.StrategyProfileResponseDTO, error) {
	settings, err := spms.strategyProfileSettingSvc.GetEnabledSettings()
	if err != nil {
		return nil, err
	}

	enabled := buildEnabledSettingsMap(settings)

	profile, err := spms.repo.GetStrategyProfileWithRelations(profileID, enabled)
	if err != nil {
		return nil, err
	}

	response := &strategyDTOs.StrategyProfileResponseDTO{
		ID:    profile.ID,
		Title: profile.Title,
	}

	if enabled["focus-area"] {
		for _, focus := range profile.StrategyProfileFocusAreas {
			response.StrategyProfileFocusAreas = append(response.StrategyProfileFocusAreas, buildFocusDTO(focus, enabled))
		}
	} else if enabled["goal"] {
		for _, goal := range profile.StrategyProfileGoals {
			response.StrategyProfileGoals = append(response.StrategyProfileGoals, buildGoalDTO(goal, enabled))
		}

	} else if enabled["objective"] {
		for _, objective := range profile.StrategyProfileObjectives {
			response.StrategyProfileObjectives = append(response.StrategyProfileObjectives, buildObjectiveDTO(objective, enabled))
		}
	} else if enabled["initiative"] {
		for _, initiative := range profile.StrategyProfileInitiatives {
			response.StrategyProfileInitiatives = append(response.StrategyProfileInitiatives, buildInitiativeDTO(initiative, enabled))
		}

	} else if enabled["project"] {
		for _, project := range profile.StrategyProfileProjects {
			response.StrategyProfileProjects = append(response.StrategyProfileProjects, buildProjectDTO(project, enabled))
		}

	} else if enabled["milestone"] {
		for _, milestone := range profile.StrategyProfileMilestones {
			response.StrategyProfileMilestones = append(response.StrategyProfileMilestones, buildMilestoneDTO(milestone, enabled))
		}

	} else if enabled["action-plan"] {
		for _, actionPlan := range profile.StrategyProfileActionPlans {
			response.StrategyProfileActionPlans = append(response.StrategyProfileActionPlans, buildActionPlanDTO(actionPlan))
		}
	}
	return response, nil
}

func buildEnabledSettingsMap(settings []strategyDTOs.StrategyProfileSettingListDTO) map[string]bool {
	enabled := make(map[string]bool)
	for _, s := range settings {
		if s.Type != nil {
			enabled[*s.Type] = s.IsEnable
		}
	}
	return enabled
}

func buildFocusDTO(focus models.StrategyProfileFocusArea, enabled map[string]bool) strategyDTOs.StrategyProfileFocusAreaDTO {
	focusDTO := strategyDTOs.StrategyProfileFocusAreaDTO{
		ID:    focus.ID,
		Title: focus.Title,
	}
	if enabled["goal"] {
		for _, goal := range focus.StrategyProfileGoals {
			focusDTO.StrategyProfileGoals = append(focusDTO.StrategyProfileGoals, buildGoalDTO(goal, enabled))
		}
	} else if enabled["objective"] {
		for _, objective := range focus.StrategyProfileObjectives {
			objectiveDTO := buildObjectiveDTO(objective, enabled)
			focusDTO.StrategyProfileObjectives = append(focusDTO.StrategyProfileObjectives, objectiveDTO)
		}
	} else if enabled["initiative"] {
		for _, initiative := range focus.StrategyProfileInitiatives {
			initiativeDTO := buildInitiativeDTO(initiative, enabled)
			focusDTO.StrategyProfileInitiatives = append(focusDTO.StrategyProfileInitiatives, initiativeDTO)
		}
	} else if enabled["project"] {
		for _, project := range focus.StrategyProfileProjects {
			projectDTO := buildProjectDTO(project, enabled)
			focusDTO.StrategyProfileProjects = append(focusDTO.StrategyProfileProjects, projectDTO)
		}
	} else if enabled["milestone"] {
		for _, milestone := range focus.StrategyProfileMilestones {
			milestoneDTO := buildMilestoneDTO(milestone, enabled)
			focusDTO.StrategyProfileMilestones = append(focusDTO.StrategyProfileMilestones, milestoneDTO)
		}
	} else if enabled["action-plan"] {
		for _, actionPlan := range focus.StrategyProfileActionPlans {
			actionPlanDTO := buildActionPlanDTO(actionPlan)
			focusDTO.StrategyProfileActionPlans = append(focusDTO.StrategyProfileActionPlans, actionPlanDTO)
		}
	}

	return focusDTO
}

func buildGoalDTO(goal models.StrategyProfileGoal, enabled map[string]bool) strategyDTOs.StrategyProfileGoalDTO {
	goalDTO := strategyDTOs.StrategyProfileGoalDTO{
		ID:    goal.ID,
		Title: goal.Title,
	}

	if enabled["objective"] {
		for _, objective := range goal.StrategyProfileObjectives {
			goalDTO.StrategyProfileObjectives = append(goalDTO.StrategyProfileObjectives, buildObjectiveDTO(objective, enabled))
		}
	} else if enabled["initiative"] {
		for _, initiative := range goal.StrategyProfileInitiatives {
			initiativeDTO := buildInitiativeDTO(initiative, enabled)
			goalDTO.StrategyProfileInitiatives = append(goalDTO.StrategyProfileInitiatives, initiativeDTO)
		}
	} else if enabled["project"] {
		for _, project := range goal.StrategyProfileProjects {
			projectDTO := buildProjectDTO(project, enabled)
			goalDTO.StrategyProfileProjects = append(goalDTO.StrategyProfileProjects, projectDTO)
		}
	} else if enabled["milestone"] {
		for _, milestone := range goal.StrategyProfileMilestones {
			milestoneDTO := buildMilestoneDTO(milestone, enabled)
			goalDTO.StrategyProfileMilestones = append(goalDTO.StrategyProfileMilestones, milestoneDTO)
		}
	} else if enabled["action-plan"] {
		for _, actionPlan := range goal.StrategyProfileActionPlans {
			actionPlanDTO := buildActionPlanDTO(actionPlan)
			goalDTO.StrategyProfileActionPlans = append(goalDTO.StrategyProfileActionPlans, actionPlanDTO)
		}
	}

	return goalDTO
}

func buildObjectiveDTO(objective models.StrategyProfileObjective, enabled map[string]bool) strategyDTOs.StrategyProfileObjectiveDTO {
	objectiveDTO := strategyDTOs.StrategyProfileObjectiveDTO{
		ID:    objective.ID,
		Title: objective.Title,
	}
	if enabled["initiative"] {
		for _, initiative := range objective.StrategyProfileInitiatives {
			objectiveDTO.StrategyProfileInitiatives = append(objectiveDTO.StrategyProfileInitiatives, buildInitiativeDTO(initiative, enabled))
		}
	} else if enabled["project"] {
		for _, project := range objective.StrategyProfileProjects {
			objectiveDTO.StrategyProfileProjects = append(objectiveDTO.StrategyProfileProjects, buildProjectDTO(project, enabled))
		}
	} else if enabled["milestone"] {
		for _, milestone := range objective.StrategyProfileMilestones {
			milestoneDTO := buildMilestoneDTO(milestone, enabled)
			objectiveDTO.StrategyProfileMilestones = append(objectiveDTO.StrategyProfileMilestones, milestoneDTO)
		}
	} else if enabled["action-plan"] {
		for _, actionPlan := range objective.StrategyProfileActionPlans {
			actionPlanDTO := buildActionPlanDTO(actionPlan)
			objectiveDTO.StrategyProfileActionPlans = append(objectiveDTO.StrategyProfileActionPlans, actionPlanDTO)
		}
	}
	return objectiveDTO
}

func buildInitiativeDTO(initiative models.StrategyProfileInitiative, enabled map[string]bool) strategyDTOs.StrategyProfileInitiativeDTO {
	initiativeDTO := strategyDTOs.StrategyProfileInitiativeDTO{
		ID:    initiative.ID,
		Title: initiative.Title,
	}

	if enabled["project"] {
		for _, project := range initiative.StrategyProfileProjects {
			initiativeDTO.StrategyProfileProjects = append(initiativeDTO.StrategyProfileProjects, buildProjectDTO(project, enabled))
		}
	} else if enabled["milestone"] {
		for _, milestone := range initiative.StrategyProfileMilestones {
			milestoneDTO := buildMilestoneDTO(milestone, enabled)
			initiativeDTO.StrategyProfileMilestones = append(initiativeDTO.StrategyProfileMilestones, milestoneDTO)
		}
	} else if enabled["action-plan"] {
		for _, actionPlan := range initiative.StrategyProfileActionPlans {
			actionPlanDTO := buildActionPlanDTO(actionPlan)
			initiativeDTO.StrategyProfileActionPlans = append(initiativeDTO.StrategyProfileActionPlans, actionPlanDTO)
		}
	}
	return initiativeDTO
}

func buildProjectDTO(project models.StrategyProfileProject, enabled map[string]bool) strategyDTOs.StrategyProfileProjectDTO {
	projectDTO := strategyDTOs.StrategyProfileProjectDTO{
		ID:    project.ID,
		Title: project.Title,
	}

	if enabled["milestone"] {
		for _, milestone := range project.StrategyProfileMilestones {
			projectDTO.StrategyProfileMilestones = append(projectDTO.StrategyProfileMilestones, buildMilestoneDTO(milestone, enabled))
		}
	} else if enabled["action-plan"] {
		for _, actionPlan := range project.StrategyProfileActionPlans {
			actionPlanDTO := buildActionPlanDTO(actionPlan)
			projectDTO.StrategyProfileActionPlans = append(projectDTO.StrategyProfileActionPlans, actionPlanDTO)
		}
	}
	return projectDTO
}

func buildMilestoneDTO(milestone models.StrategyProfileMilestone, enabled map[string]bool) strategyDTOs.StrategyProfileMilestoneDTO {
	milestoneDTO := strategyDTOs.StrategyProfileMilestoneDTO{
		ID:    milestone.ID,
		Title: milestone.Title,
	}

	if enabled["action-plan"] {
		for _, action := range milestone.StrategyProfileActionPlans {
			milestoneDTO.StrategyProfileActionPlans = append(milestoneDTO.StrategyProfileActionPlans, buildActionPlanDTO(action))
		}
	}
	return milestoneDTO
}

func buildActionPlanDTO(action models.StrategyProfileActionPlan) strategyDTOs.StrategyProfileActionPlanDTO {
	return strategyDTOs.StrategyProfileActionPlanDTO{
		ID:    action.ID,
		Title: action.Title,
	}
}
