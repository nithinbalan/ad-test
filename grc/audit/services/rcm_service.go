package services

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"time"
)

type RCMServiceInterface interface {
	GetAll(params *auditDTOs.RCMListQueryParamsDTO, auditID uint) (*auditDTOs.RCMListResponseDTO, error)
}

type RCMService struct {
	repo           repositories.RCMRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditSvc       AuditServiceInterface
}

func NewRCMService(
	repo repositories.RCMRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface) RCMServiceInterface {
	return &RCMService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditSvc:       auditSvc,
	}
}

func (svc *RCMService) GetAll(params *auditDTOs.RCMListQueryParamsDTO, auditID uint) (*auditDTOs.RCMListResponseDTO, error) {
	audit, err := svc.auditSvc.GetByID(auditID)
	if err != nil {
		return nil, err
	}

	rawData, totalCount, err := svc.repo.GetAll(params, audit.AuditPlanID)
	if err != nil {
		return nil, err
	}

	// Use a map to group controls by their risk ID
	rcmResponseMap := make(map[uint]*auditDTOs.RCMListDTO)
	for _, row := range rawData {
		riskID := uint(row["id"].(int64))

		if _, exists := rcmResponseMap[riskID]; !exists {

			var createdAt time.Time
			if row["created_at"] != nil {
				createdAt, _ = row["created_at"].(time.Time)
			}

			var updatedAt *time.Time
			if row["updated_at"] != nil {
				if t, ok := row["updated_at"].(time.Time); ok {
					updatedAt = &t
				}
			}

			// Initialize a new RCMListDTO for this risk
			rcmResponseMap[riskID] = &auditDTOs.RCMListDTO{
				ID:                       riskID,
				AuditPlanAuditableItemID: uint(row["audit_plan_auditable_item_id"].(int64)),
				RiskRegisterID:           uint(row["risk_register_id"].(int64)),
				RiskTitle:                convertToStringOrEmpty(row["risk_title"]),
				CreatedAt:                createdAt,
				UpdatedAt:                updatedAt,
				Controls:                 []*auditDTOs.RCMControlDTO{},
			}
		}

		// If control details are available, add them to the respective risk
		if row["control_id"] != nil {
			control := &auditDTOs.RCMControlDTO{
				ID:                                   uint(row["control_id"].(int64)),
				ControlID:                            uint(row["control_id"].(int64)),
				ControlTitleEn:                       convertToStringOrEmpty(row["control_title_en"]),
				ControlTitleAr:                       convertToStringOrEmpty(row["control_title_ar"]),
				ControlImplementationStatusTitleEn:   convertToStringPointer(row["control_implementation_status_title_en"]),
				ControlImplementationStatusTitleAr:   convertToStringPointer(row["control_implementation_status_title_ar"]),
				ControlImplementationStatusColorCode: convertToStringPointer(row["control_implementation_status_color_code"]),
				ControlEfficiencyStatusTitleEn:       convertToStringPointer(row["control_efficiency_status_title_en"]),
				ControlEfficiencyStatusTitleAr:       convertToStringPointer(row["control_efficiency_status_title_ar"]),
				ControlEfficiencyStatusColorCode:     convertToStringPointer(row["control_efficiency_status_color_code"]),
			}

			// Append the control to the respective risk's control list
			rcmResponseMap[riskID].Controls = append(rcmResponseMap[riskID].Controls, control)
		}
	}

	// Prepare the response DTO
	rcmResponse := &auditDTOs.RCMListResponseDTO{
		Items: make([]*auditDTOs.RCMListDTO, 0, len(rcmResponseMap)),
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  totalCount,
		},
	}

	// Convert map to slice for the final response
	for _, item := range rcmResponseMap {
		rcmResponse.Items = append(rcmResponse.Items, item)
	}

	return rcmResponse, nil
}

// Helper function to safely convert an interface{} to a string
func convertToStringOrEmpty(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func convertToStringPointer(value interface{}) *string {
	if value == nil {
		return nil
	}
	strVal := fmt.Sprintf("%v", value)
	if strVal == "" {
		return nil
	}
	return &strVal
}
