package dtos

type IncidentEscalationMatrixListDTO struct {
	ID uint `json:"id"`

	SeverityID             uint   `json:"severity_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelType      string `json:"severity_level_type"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelWeight    int    `json:"severity_level_weight"`

	NotifyWithin            *int    `json:"notify_within"`
	NotifyWithinUnitID      *uint   `json:"notify_within_unit_id"`
	NotifyWithinUnitTitleEn *string `json:"notify_within_unit_title_en"`
	NotifyWithinUnitTitleAr *string `json:"notify_within_unit_title_ar"`

	NotifyToIncidentOwner bool `json:"notify_to_incident_owner"`
	NotifyToResponseTeam  bool `json:"notify_to_response_team"`

	NotifyRoleID      *uint   `json:"notify_role_id"`
	NotifyRoleTitleEn *string `json:"notify_role_title_en"`
	NotifyRoleTitleAr *string `json:"notify_role_title_ar"`

	EscalateAfterIncidentOwner bool `json:"escalate_after_incident_owner"`
	EscalateAfterResponseTeam  bool `json:"escalate_after_response_team"`

	EscalateRoleID      *uint   `json:"escalate_role_id"`
	EscalateRoleTitleEn *string `json:"escalate_role_title_en"`
	EscalateRoleTitleAr *string `json:"escalate_role_title_ar"`

	Status string `json:"status"`
}

type IncidentEscalationMatrixFormDTO struct {
	SeverityID                 uint  `json:"severity_id" binding:"required"`
	NotifyWithin               *int  `json:"notify_within"`
	NotifyWithinUnitID         *uint `json:"notify_within_unit_id"`
	NotifyToIncidentOwner      bool  `json:"notify_to_incident_owner"`
	NotifyToResponseTeam       bool  `json:"notify_to_response_team"`
	NotifyRoleID               *uint `json:"notify_role_id"`
	EscalateAfterIncidentOwner bool  `json:"escalate_after_incident_owner"`
	EscalateAfterResponseTeam  bool  `json:"escalate_after_response_team"`
	EscalateRoleID             *uint `json:"escalate_role_id"`
}

type IncidentEscalationMatrixUpdateDTO struct {
	Settings []*IncidentEscalationMatrixFormDTO `json:"settings" binding:"required"`
}

type IncidentResolutionMatrixListDTO struct {
	ID uint `json:"id"`

	SeverityID             uint   `json:"severity_id"`
	SeverityLevelTitleEn   string `json:"severity_level_title_en"`
	SeverityLevelTitleAr   string `json:"severity_level_title_ar"`
	SeverityLevelType      string `json:"severity_level_type"`
	SeverityLevelColorCode string `json:"severity_level_color_code"`
	SeverityLevelWeight    int    `json:"severity_level_weight"`

	CloseWithin            *int    `json:"close_within"`
	CloseWithinUnitID      *uint   `json:"close_within_unit_id"`
	CloseWithinUnitTitleEn *string `json:"close_within_unit_title_en"`
	CloseWithinUnitTitleAr *string `json:"close_within_unit_title_ar"`

	CloseWithinIncidentOwner bool `json:"close_within_incident_owner"`
	CloseWithinResponseTeam  bool `json:"close_within_response_team"`

	CloseWithinRoleID      *uint   `json:"close_within_role_id"`
	CloseWithinRoleTitleEn *string `json:"close_within_role_title_en"`
	CloseWithinRoleTitleAr *string `json:"close_within_role_title_ar"`

	Status string `json:"status"`
}

type IncidentResolutionMatrixFormDTO struct {
	SeverityID               uint  `json:"severity_id" binding:"required"`
	CloseWithin              *int  `json:"close_within"`
	CloseWithinUnitID        *uint `json:"close_within_unit_id"`
	CloseWithinIncidentOwner bool  `json:"close_within_incident_owner"`
	CloseWithinResponseTeam  bool  `json:"close_within_response_team"`
	CloseWithinRoleID        *uint `json:"close_within_role_id"`
}

type IncidentResolutionMatrixUpdateDTO struct {
	Settings []*IncidentResolutionMatrixFormDTO `json:"settings" binding:"required,dive"`
}
