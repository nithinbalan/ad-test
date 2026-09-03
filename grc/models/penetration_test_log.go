package models

import "time"

type PenetrationTest struct {
	ID                              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReferenceNumber                 string     `json:"reference_number" gorm:"type:varchar(45);not null;unique"`
	SystemApplicationName           string     `json:"system_application_name" gorm:"type:nvarchar(255);not null" binding:"required"`
	PenetrationTestTypeID           uint       `json:"penetration_test_type_id" gorm:"not null" binding:"required"`
	PenetrationTestPlanID           uint       `json:"penetration_test_plan_id" gorm:"not null" binding:"required"`
	Year                            int        `json:"year" gorm:"not null" binding:"required"`
	Period                          string     `json:"period" gorm:"type:varchar(2);not null" binding:"required"`
	PenetrationTestClassificationID uint       `json:"penetration_test_classification_id" gorm:"not null" binding:"required"`
	PenetrationTestExecutorID       uint       `json:"penetration_test_executor_id" gorm:"not null" binding:"required"`
	CSResponsibleEmployeeID         uint       `json:"cs_responsible_employee_id" gorm:"not null" binding:"required"`
	Scope                           string     `json:"scope" gorm:"type:nvarchar(max);not null" binding:"required"`
	CreatedBYID                     uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID                     *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt                       *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	PenetrationTestType           *PenetrationTestType                   `json:"penetration_test_type,omitempty" gorm:"foreignKey:PenetrationTestTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestPlan           *PenetrationTestPlan                   `json:"penetration_test_plan,omitempty" gorm:"foreignKey:PenetrationTestPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestClassification *PenetrationTestClassification         `json:"penetration_test_classification,omitempty" gorm:"foreignKey:PenetrationTestClassificationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PenetrationTestExecutor       *PenetrationTestExecutor               `json:"penetration_test_executor,omitempty" gorm:"foreignKey:PenetrationTestExecutorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CSResponsibleEmployee         *User                                  `json:"cs_responsible_employee,omitempty" gorm:"foreignKey:CSResponsibleEmployeeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	EnvironmentPreparation        *PenetrationTestEnvironmentPreparation `json:"environment_preparation,omitempty" gorm:"foreignKey:PenetrationTestID"`
	Execution                     *PenetrationTestExecution              `json:"execution,omitempty" gorm:"foreignKey:PenetrationTestID"`
	ManagementSharing             *PenetrationTestManagementSharing      `json:"management_sharing,omitempty" gorm:"foreignKey:PenetrationTestID"`
	CreatedByUser                 *User                                  `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser                 *User                                  `json:"updated_by_user,omitempty" gorm:"foreignKey:UpdatedBYID"`
}

func (PenetrationTest) TableName() string {
	return "penetration_tests"
}
