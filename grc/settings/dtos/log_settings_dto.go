package dtos

type LogSettingsFormDTO struct {
	AuditLogRetentionDays  int `json:"audit_log_retention_days" binding:"required"`
	EmailLogRetentionDays  int `json:"email_log_retention_days" binding:"required"`
	SystemLogRetentionDays int `json:"system_log_retention_days" binding:"required"`
}
