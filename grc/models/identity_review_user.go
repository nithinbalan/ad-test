package models

import "time"

type IdentityReviewUser struct {
	ID                       uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IdentityReviewCampaignID uint       `json:"identity_review_campaign_id" gorm:"index"`
	Username                 string     `json:"username" gorm:"type:nvarchar(255);not null"`
	UserIdentifier           *string    `json:"user_identifier" gorm:"type:nvarchar(255);default:null"`
	LastLogin                *string    `json:"last_login" gorm:"type:nvarchar(255);default:null"`
	DepartmentID             uint       `json:"department_id" gorm:"not null"`
	PrivilegeID              *uint      `json:"privilege_id" gorm:"default:null"`
	AccountTypeID            *uint      `json:"account_type_id" gorm:"default:null"`
	AccountDecisionID        *uint      `json:"account_decision_id" gorm:"default:null"`
	CurrentRole              *string    `json:"current_role" gorm:"type:nvarchar(max);default:null"`
	ReviewComment            *string    `json:"review_comment" gorm:"type:nvarchar(max);default:null"`
	ActionOwnerID            *uint      `json:"action_owner_id" gorm:"default:null"`
	ActionPlanStatusID       *uint      `json:"action_plan_status_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              *uint      `json:"updated_by_id" gorm:"default:null"`

	IdentityReviewCampaign *IdentityReviewCampaign `json:"identity_review_campaign" gorm:"foreignKey:IdentityReviewCampaignID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Department             *OrganizationUnit       `json:"department" gorm:"foreignKey:DepartmentID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Privilege              *PrivilegedAccount      `json:"privilege" gorm:"foreignKey:PrivilegeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	AccountType            *AccountType            `json:"account_type" gorm:"foreignKey:AccountTypeID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	AccountDecision        *AccountDecision        `json:"account_decision" gorm:"foreignKey:AccountDecisionID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ActionOwner            *User                   `json:"action_owner" gorm:"foreignKey:ActionOwnerID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ActionPlanStatus       *ActionPlanStatus       `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser          *User                   `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	UpdatedByUser          *User                   `json:"updated_by" gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

func (IdentityReviewUser) TableName() string {
	return "identity_review_users"
}
