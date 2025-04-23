package dtos

type StrategyProfileResponseDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileFocusAreas  []StrategyProfileFocusAreaDTO  `json:"strategy_profile_focus_areas"`
	StrategyProfileGoals       []StrategyProfileGoalDTO       `json:"strategy_profile_goals"`
	StrategyProfileObjectives  []StrategyProfileObjectiveDTO  `json:"strategy_profile_objectives"`
	StrategyProfileInitiatives []StrategyProfileInitiativeDTO `json:"strategy_profile_initiatives"`
	StrategyProfileProjects    []StrategyProfileProjectDTO    `json:"strategy_profile_projects"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileFocusAreaDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileGoals       []StrategyProfileGoalDTO       `json:"strategy_profile_goals"`
	StrategyProfileObjectives  []StrategyProfileObjectiveDTO  `json:"strategy_profile_objectives"`
	StrategyProfileInitiatives []StrategyProfileInitiativeDTO `json:"strategy_profile_initiatives"`
	StrategyProfileProjects    []StrategyProfileProjectDTO    `json:"strategy_profile_projects"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileGoalDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileObjectives  []StrategyProfileObjectiveDTO  `json:"strategy_profile_objectives"`
	StrategyProfileInitiatives []StrategyProfileInitiativeDTO `json:"strategy_profile_initiatives"`
	StrategyProfileProjects    []StrategyProfileProjectDTO    `json:"strategy_profile_projects"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileObjectiveDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileInitiatives []StrategyProfileInitiativeDTO `json:"strategy_profile_initiatives"`
	StrategyProfileProjects    []StrategyProfileProjectDTO    `json:"strategy_profile_projects"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileInitiativeDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileProjects    []StrategyProfileProjectDTO    `json:"strategy_profile_projects"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileProjectDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileMilestones  []StrategyProfileMilestoneDTO  `json:"strategy_profile_milestones"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileMilestoneDTO struct {
	ID                         uint                           `json:"id"`
	Title                      string                         `json:"title"`
	StrategyProfileActionPlans []StrategyProfileActionPlanDTO `json:"strategy_profile_action_plans"`
}

type StrategyProfileActionPlanDTO struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}
