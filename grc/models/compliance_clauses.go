package models

import "time"

// ComplianceClause represents a clause in the compliance system.
type ComplianceClause struct {
	ID                   uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ComplianceRegisterID uint       `json:"compliance_register_id" binding:"required" gorm:"not null"`
	ParentClauseID       *uint      `json:"compliance_clause_id" gorm:"index;default:null;column:compliance_clause_id"`
	ClauseNo             string     `json:"clause_no" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	TitleEn              string     `json:"title_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	TitleAr              *string    `json:"title_ar" gorm:"type:nvarchar(max);default:null"`
	DescriptionEn        *string    `json:"description_en" gorm:"type:nvarchar(max);default:null"`
	DescriptionAr        *string    `json:"description_ar" gorm:"type:nvarchar(max);default:null"`
	ScopeEn              *string    `json:"scope_en" gorm:"type:nvarchar(max);default:null"`
	ScopeAr              *string    `json:"scope_ar" gorm:"type:nvarchar(max);default:null"`
	ObjectiveEn          *string    `json:"objective_en" gorm:"type:nvarchar(max);default:null"`
	ObjectiveAr          *string    `json:"objective_ar" gorm:"type:nvarchar(max);default:null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	ComplianceRegister     *ComplianceRegister      `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ParentClause           *ComplianceClause        `gorm:"foreignKey:ParentClauseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // Self-referencing relationship
	ChildClauses           []*ComplianceClause      `gorm:"foreignKey:ParentClauseID"`                                               // Recursive relationship
	CreatedByUser          *User                    `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser          *User                    `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRequirements []*ComplianceRequirement `gorm:"foreignKey:ParantClauseID"`
}

// TableName overrides the table name used by ComplianceClause to `compliance_clauses`
func (ComplianceClause) TableName() string {
	return "compliance_clauses"
}
