package services

import (
	"fmt"
	"time"

	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
)

type ComplianceAssessmentRequirementResidualRiskAnalysisServiceInterface interface {
	// Residual Risk Analysis methods
	CreateResidualRiskAnalysis(responsibleUnitID uint, currentUserID uint, riskAnalysisFormDTO []*complianceDTOs.ComplianceAssessmentRequirementResidualRiskFormDTO) error
	DeleteResidualRiskAnalysis(responsibleUnitID uint, currentUserID uint) error
	CalculateResidualRiskScore(responsibleUnitID uint) error

	GetAllByResponsibleUnitID(responsibleUnitID uint, params *complianceDTOs.ComplianceAssessmentResidualControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentResidualControlListResponseDTO, error)
	CreateControl(responsibleUnitID uint, currentUserID uint, residualControlFormDTO *complianceDTOs.ComplianceAssessmentResidualControlFormDTO) error
	DeleteControl(responsibleUnitID uint, currentUserID uint) error
	UpdateResidualControlByImpactType(responsibleUnitID uint, impactTypeID uint, tx ...utility.DatabaseTransaction) error
	CreateResidualControlNotApplicable(responsibleUnitID uint, impactTypeID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceAssessmentRequirementResidualRiskAnalysisService struct {
	repo           complianceRepositories.ComplianceAssessmentRequirementResidualRiskAnalysisRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewComplianceAssessmentRequirementResidualRiskAnalysisService(
	repo complianceRepositories.ComplianceAssessmentRequirementResidualRiskAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) ComplianceAssessmentRequirementResidualRiskAnalysisServiceInterface {
	return &ComplianceAssessmentRequirementResidualRiskAnalysisService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

// getResidualRiskRatingByScore determines the appropriate risk rating based on the calculated score
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) getResidualRiskRatingByScore(score float64) (*uint, error) {
	// Get all compliance assessment ratings
	ratings, err := carras.repo.GetComplianceAssessmentRatings()
	if err != nil {
		return nil, err
	}

	// Find the appropriate rating based on score range
	for _, rating := range ratings {
		if rating.ScoreFrom != nil && rating.ScoreTo != nil {
			if score >= *rating.ScoreFrom && score <= *rating.ScoreTo {
				return &rating.ID, nil
			}
		}
	}

	return nil, nil // No rating found
}

// CreateResidualRiskAnalysis creates residual risk analysis for a compliance assessment requirement responsible unit
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) CreateResidualRiskAnalysis(responsibleUnitID uint, currentUserID uint, riskAnalysisFormDTO []*complianceDTOs.ComplianceAssessmentRequirementResidualRiskFormDTO) error {
	tx := carras.repo.BeginTransaction()

	// Delete existing risk analysis for this responsible unit
	if err := carras.repo.DeleteResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID, tx); err != nil {
		tx.Rollback()
		return err
	}

	var riskAnalysisRecords []*models.ComplianceAssessmentRequirementResidualRiskAnalysis
	var totalScore float64

	// Create new risk analysis records
	for _, analysis := range riskAnalysisFormDTO {
		if analysis.Score == nil {
			return fmt.Errorf("score is required")
		}
		record := &models.ComplianceAssessmentRequirementResidualRiskAnalysis{
			ComplianceAssessmentRequirementResponsibleUnitID: responsibleUnitID,
			ComplianceAssessmentImpactTypeID:                 analysis.ComplianceAssessmentImpactTypeID,
			Score:                                            *analysis.Score,
			Comment:                                          analysis.Comment,
			CreatedBYID:                                      currentUserID,
			CreatedAt:                                        time.Now(),
		}
		riskAnalysisRecords = append(riskAnalysisRecords, record)
		totalScore += float64(*analysis.Score)
	}

	// Create residual risk analysis records
	if err := carras.repo.CreateResidualRiskAnalysis(riskAnalysisRecords, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Calculate and update residual risk score
	if len(riskAnalysisFormDTO) > 0 {
		averageScore := totalScore / float64(len(riskAnalysisFormDTO))
		ratingID, err := carras.getResidualRiskRatingByScore(averageScore)
		if err != nil {
			tx.Rollback()
			return err
		}
		if err := carras.repo.UpdateResidualRiskScore(responsibleUnitID, averageScore, ratingID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// DeleteResidualRiskAnalysis deletes all residual risk analysis for a responsible unit
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) DeleteResidualRiskAnalysis(responsibleUnitID uint, currentUserID uint) error {
	tx := carras.repo.BeginTransaction()

	// Delete all risk analysis for this responsible unit
	if err := carras.repo.DeleteResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Reset residual risk score and rating
	if err := carras.repo.UpdateResidualRiskScore(responsibleUnitID, 0, nil, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// CalculateResidualRiskScore recalculates the residual risk score for a responsible unit
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) CalculateResidualRiskScore(responsibleUnitID uint) error {
	tx := carras.repo.BeginTransaction()

	// Get all risk analysis for this responsible unit
	riskAnalysisList, err := carras.repo.GetResidualRiskAnalysisByResponsibleUnitID(responsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(riskAnalysisList) == 0 {
		// No risk analysis, reset score
		if err := carras.repo.UpdateResidualRiskScore(responsibleUnitID, 0, nil, tx); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	var totalScore float64
	for _, analysis := range riskAnalysisList {
		totalScore += float64(analysis.Score)
	}

	// Calculate average residual risk score
	averageScore := totalScore / float64(len(riskAnalysisList))

	// Get appropriate risk rating based on score
	ratingID, err := carras.getResidualRiskRatingByScore(averageScore)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update responsible unit with residual risk score and rating
	if err := carras.repo.UpdateResidualRiskScore(responsibleUnitID, averageScore, ratingID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// GetAllByResponsibleUnitID retrieves all residual controls for a responsible unit with pagination
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) GetAllByResponsibleUnitID(responsibleUnitID uint, params *complianceDTOs.ComplianceAssessmentResidualControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentResidualControlListResponseDTO, error) {
	return carras.repo.GetAllByResponsibleUnitID(responsibleUnitID, params)
}

// Create creates residual control mappings for a responsible unit
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) CreateControl(responsibleUnitID uint, currentUserID uint, residualControlFormDTO *complianceDTOs.ComplianceAssessmentResidualControlFormDTO) error {
	tx := carras.repo.BeginTransaction()

	// Create new residual control record
	record := &models.ComplianceAssessmentResidualControl{
		ComplianceAssessmentRequirementResponsibleUnitID: responsibleUnitID,
		ComplianceAssessmentImpactTypeID:                 residualControlFormDTO.ComplianceAssessmentImpactTypeID,
		ControlID:                                        residualControlFormDTO.ControlID,
		LinkControl:                                      true,
		CreatedBYID:                                      currentUserID,
		CreatedAt:                                        time.Now(),
	}

	var residualControlRecords []*models.ComplianceAssessmentResidualControl
	residualControlRecords = append(residualControlRecords, record)

	// Create residual control records
	if err := carras.repo.Create(residualControlRecords, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		err := carras.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceAssessmentResidualControl{}.TableName(),
			responsibleUnitID,
			residualControlRecords,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":            models.CreateAction,
				"tableName":         models.ComplianceAssessmentResidualControl{}.TableName(),
				"responsibleUnitID": responsibleUnitID,
				"userID":            currentUserID,
				"error":             err.Error(),
			}
			// Log error (utility.LogError would be used here)
			fmt.Printf("Failed to log activity: %v\n", fields)
		}
	}()

	return nil
}

// Delete deletes all residual control mappings for a responsible unit
func (carras *ComplianceAssessmentRequirementResidualRiskAnalysisService) DeleteControl(responsibleUnitID uint, currentUserID uint) error {
	tx := carras.repo.BeginTransaction()

	// Delete all residual controls for this responsible unit
	if err := carras.repo.DeleteByResponsibleUnitID(responsibleUnitID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		err := carras.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.ComplianceAssessmentResidualControl{}.TableName(),
			responsibleUnitID,
			nil,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":            models.DeleteAction,
				"tableName":         models.ComplianceAssessmentResidualControl{}.TableName(),
				"responsibleUnitID": responsibleUnitID,
				"userID":            currentUserID,
				"error":             err.Error(),
			}
			// Log error (utility.LogError would be used here)
			fmt.Printf("Failed to log activity: %v\n", fields)
		}
	}()

	return nil
}

// UpdateResidualControlByImpactType updates IsNotApplicable flag for residual control by impact type
func (carcs *ComplianceAssessmentRequirementResidualRiskAnalysisService) UpdateResidualControlByImpactType(responsibleUnitID uint, impactTypeID uint, tx ...utility.DatabaseTransaction) error {
	return carcs.repo.UpdateResidualControlByImpactType(responsibleUnitID, impactTypeID, tx...)
}

// CreateResidualControlNotApplicable creates a new residual control record with is_not_applicable=true
func (carcs *ComplianceAssessmentRequirementResidualRiskAnalysisService) CreateResidualControlNotApplicable(responsibleUnitID uint, impactTypeID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	// Create new residual control record with is_not_applicable=true
	record := &models.ComplianceAssessmentResidualControl{
		ComplianceAssessmentRequirementResponsibleUnitID: responsibleUnitID,
		ComplianceAssessmentImpactTypeID:                 impactTypeID,
		ControlID:                                        nil, // No control linked
		LinkControl:                                      false,
		IsNotApplicable:                                  true, // Set to true for action plan creation
		CreatedBYID:                                      currentUserID,
		CreatedAt:                                        time.Now(),
	}

	var residualControlRecords []*models.ComplianceAssessmentResidualControl
	residualControlRecords = append(residualControlRecords, record)

	// Create residual control record
	return carcs.repo.Create(residualControlRecords, tx...)
}
