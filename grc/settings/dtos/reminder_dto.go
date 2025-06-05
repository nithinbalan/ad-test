package dtos

import "time"

type ReminderDTO struct {
	ID       uint   `json:"id"`
	TitleEn  string `json:"title_en"`
	TitleAr  string `json:"title_ar"`
	ModuleID uint   `json:"module_id"`
	Type     string `json:"type"`
	Status   string `json:"status"`
}

type ReminderListQueryParamsDTO struct {
	Sort    string
	SortBy  string
	Keyword string
	Status  string
	Filters map[string][]string
}

type GroupedReminderDTO struct {
	ID            uint                    `json:"id"`
	TitleEn       string                  `json:"title_en"`
	TitleAr       string                  `json:"title_ar"`
	Type          string                  `json:"type"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     time.Time               `json:"updated_at"`
	ReminderRules []ReminderRuleDetailDTO `json:"reminder_rules"`
}
