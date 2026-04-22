package models

type RiskRegisterUpdatePlanStatus struct {
	ID        int    `json:"id" db:"id"`
	TitleEn   string `json:"title_en" db:"title_en"`
	TitleAr   string `json:"title_ar" db:"title_ar"`
	ColorCode string `json:"color_code" db:"color_code"`
	Type      string `json:"type" db:"type" binding:"required,max=20,oneof=not_started in_progress overdue completed"`
	Status    string `json:"status" db:"status" binding:"required,max=10,oneof=active inactive"`
}

func (RiskRegisterUpdatePlanStatus) TableName() string {
	return "risk_register_update_plan_statuses"
}