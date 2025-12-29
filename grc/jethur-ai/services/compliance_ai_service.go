package services

import (
	"fmt"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type ComplianceAIServiceInterface interface {
	GetComplianceAutofill(req dtos.ComplianceAutofillRequest) (map[string]interface{}, error)
}

type ComplianceAIService struct {
	systemPreferenceSvc services.SystemPreferenceServiceInterface
}

func NewComplianceAIService(systemPreferenceSvc services.SystemPreferenceServiceInterface) ComplianceAIServiceInterface {
	return &ComplianceAIService{
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (cais *ComplianceAIService) GetComplianceAutofill(req dtos.ComplianceAutofillRequest) (map[string]interface{}, error) {
	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"title": req.Title,
		"code":  req.Code,
	}

	return utility.CallAIAPI("/api/v1/compliance/autofill", payload)
}
