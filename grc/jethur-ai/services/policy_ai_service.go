package services

import (
	"fmt"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type PolicyAIServiceInterface interface {
	GetPolicyAutofill(req dtos.PolicyAutofillRequest) (map[string]interface{}, error)
}


type PolicyAIService struct {
	systemPreferenceSvc services.SystemPreferenceServiceInterface
}

func NewPolicyAIService(systemPreferenceSvc services.SystemPreferenceServiceInterface) PolicyAIServiceInterface {
	return &PolicyAIService{
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (pais *PolicyAIService) GetPolicyAutofill(req dtos.PolicyAutofillRequest) (map[string]interface{}, error) {
	systemPreference, err := pais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"title": req.Title,
	}

	return utility.CallAIAPI("/api/v1/policy/autofill", payload)
}