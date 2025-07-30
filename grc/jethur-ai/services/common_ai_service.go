package services

import (
	"fmt"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
)

type CommonAIServiceInterface interface {
	RewriteDescription(dtos.DescriptionRewriteRequest) (map[string]interface{}, error)
	TranslateLanguage(dtos.LanguageTranslationRequest) (map[string]interface{}, error)
}

type CommonAIService struct {
	systemPreferenceSvc services.SystemPreferenceServiceInterface
}

func NewCommonAIService(systemPreferenceSvc services.SystemPreferenceServiceInterface) CommonAIServiceInterface {
	return &CommonAIService{
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (s *CommonAIService) RewriteDescription(req dtos.DescriptionRewriteRequest) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	return utility.CallAIAPI("/api/v1/description-rewriter", req)
}

func (s *CommonAIService) TranslateLanguage(req dtos.LanguageTranslationRequest) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	return utility.CallAIAPI("/api/v1/language-translation", req)
}
