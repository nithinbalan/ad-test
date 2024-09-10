package dtos

type DBNotificationDTO struct {
	UserID               uint    `json:"user_id"`
	ModuleIdentifier     string  `json:"module_identifier"`
	NotificationType     string  `json:"notification_type"`
	ReferenceNumber      string  `json:"reference_number"`
	HeadingEn            string  `json:"heading_en"`
	HeadingAr            string  `json:"heading_ar"`
	TitleEn              string  `json:"title_en"`
	TitleAr              string  `json:"title_ar"`
	Entity               *string `json:"entity"`
	Action               *string `json:"action"`
	URLParam             string  `json:"url_param"`
	QueryParam           string  `json:"query_param"`
	IsSystemNotification bool    `json:"is_system_notification"`
	CreatedByID          *uint   `json:"created_by_id"`
}

type EmailNotificationDTO struct {
	EmailTo              []string    `json:"email_to"`
	Subject              string      `json:"subject"`
	TemplateRelativePath string      `json:"template_relative_path"`
	EmailData            interface{} `json:"email_data"`
}

type NotificationType string

const (
	AssignmentType     NotificationType = "assignment"
	SecurityType       NotificationType = "security"
	ActionType         NotificationType = "action"
	DeletionType       NotificationType = "deletion"
	CompletionType     NotificationType = "completion"
	ReminderType       NotificationType = "reminder"
	WarningType        NotificationType = "warning"
	InformationType    NotificationType = "information"
	AcknowledgmentType NotificationType = "acknowledgment"
	EscalationType     NotificationType = "escalation"
	DeadlineType       NotificationType = "deadline"
	ApprovalType       NotificationType = "approval"
	ConfirmationType   NotificationType = "confirmation"
	ErrorType          NotificationType = "error"
)
