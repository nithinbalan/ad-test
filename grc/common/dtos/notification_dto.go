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
	CC                   []string    `json:"cc"`
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

type SocialMediaLinksDTO struct {
	Linkedin string `json:"linkedin"`
	X        string `json:"x"`
	Facebook string `json:"facebook"`
}

type EmailTemplateAndSocialDTO struct {
	HeaderImageURL string  `json:"header_image_url"`
	SupportText    string  `json:"support_text"`
	ClosingText    string  `json:"closing_text"`
	FooterLogoURL  string  `json:"footer_logo_url"`
	FooterAddress  string  `json:"footer_address"`
	FooterEmail    string  `json:"footer_email"`
	Linkedin       *string `json:"linkedin"`
	X              *string `json:"x"`
	Facebook       *string `json:"facebook"`
}
