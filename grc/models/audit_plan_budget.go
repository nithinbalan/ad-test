package models

import (
	"time"
)

type AuditPlanBudget struct {
	ID                uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID       uint       `json:"audit_plan_id" gorm:"not null"`
	ExpenseCategoryID uint       `json:"expense_category_id" gorm:"not null"`
	Title             string     `json:"title" gorm:"type:nvarchar(255);not null"`
	EstimatedCost     float64    `json:"estimated_cost" gorm:"type:decimal(18,2);not null"`
	Comments          *string    `json:"comments" gorm:"type:nvarchar(MAX);default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID       uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditPlan       *AuditPlan       `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ExpenseCategory *ExpenseCategory `gorm:"foreignKey:ExpenseCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser   *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser   *User            `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by AuditPlanBudget to `audit_plan_budgets`
func (AuditPlanBudget) TableName() string {
	return "audit_plan_budgets"
}
