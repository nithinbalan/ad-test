package models

import (
	"time"
)

type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Action      Action    `json:"action" gorm:"type:varchar(50);not null"`
	TableName   string    `json:"table_name" gorm:"type:varchar(100);not null"`
	RecordID    uint      `json:"record_id"`
	ChangedData string    `json:"changed_data" gorm:"type:nvarchar(max)"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null;default:GETDATE()"`

	User *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

type Action string

const (
	LoginAction  Action = "login"
	LogoutAction Action = "logout"

	CreateAction Action = "create"
	ViewAction   Action = "view"
	UpdateAction Action = "update"
	DeleteAction Action = "delete"

	ActivateAction   Action = "activate"
	DeactivateAction Action = "deactivate"
	LockAction       Action = "lock"
	UnlockAction     Action = "unlock"

	AccessAction          Action = "access"
	EditPermissionsAction Action = "edit_permissions"
	ChangeRoleAction      Action = "change_role"

	ApproveAction      Action = "approve"
	RejectAction       Action = "reject"
	AssignTaskAction   Action = "assign_task"
	CompleteTaskAction Action = "complete_task"
	EscalateAction     Action = "escalate"
	SubmitAction       Action = "submit"
	ScheduleAction     Action = "schedule"

	ImportAction        Action = "import"
	ExportAction        Action = "export"
	DownloadAction      Action = "download"
	UploadAction        Action = "upload"
	RemoveFileAction    Action = "remove_file"
	RestoreAction       Action = "restore"
	ArchiveAction       Action = "archive"
	UnarchiveAction     Action = "unarchive"
	ChangeSettingAction Action = "change_setting"

	GenerateReportAction   Action = "generate_report"
	SendNotificationAction Action = "send_notification"

	CommentAction    Action = "comment"
	UpdateNoteAction Action = "update_note"
)
