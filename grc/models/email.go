package models

import "time"

type Email struct {
	ID                 uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ModuleID           uint       `json:"module_id" gorm:"default:null"`
	SubModuleID        uint       `json:"sub_module_id" gorm:"default:null"`
	Subject            string     `json:"subject" gorm:"type:nvarchar(500);not null"`
	Body               string     `json:"body" gorm:"type:nvarchar(max);not null"`
	EmailStatusID      uint       `json:"email_status_id" gorm:"not null"`
	IsPermanentFailed  bool       `json:"is_permanent_failed" gorm:"not null;default:false"`
	AttemptCount       int        `json:"attempt_count" gorm:"not null;default:0"`
	MaxAttempts        int        `json:"max_attempts" gorm:"not null;default:5"`
	NextRetryAt        *time.Time `json:"next_retry_at" gorm:"type:datetime;default:null"`
	LastAttemptAt      *time.Time `json:"last_attempt_at" gorm:"type:datetime;default:null"`
	LockedAt           *time.Time `json:"locked_at" gorm:"type:datetime;default:null"`
	LockedBy           *string    `json:"locked_by" gorm:"type:nvarchar(64);default:null"`
	SmtpRemoteResponse *string    `json:"smtp_remote_response" gorm:"type:nvarchar(max);default:null"`
	ErrorMessage       *string    `json:"error_message" gorm:"type:nvarchar(max);default:null"`
	SentAt             *time.Time `json:"sent_at" gorm:"type:datetime;default:null"`
	CreatedByID        *uint      `json:"created_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	Recipients  []*EmailRecipient  `json:"recipients" gorm:"foreignKey:EmailID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Attachments []*EmailAttachment `json:"attachments" gorm:"foreignKey:EmailID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedBy   *User              `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Module      *Module            `json:"module" gorm:"foreignKey:ModuleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SubModule   *Module            `json:"sub_module" gorm:"foreignKey:SubModuleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	EmailStatus *EmailStatus       `json:"email_status" gorm:"foreignKey:EmailStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (Email) TableName() string {
	return "emails"
}
