package services

import (
	"fmt"
	"jethur/grc/jethur-ai/dtos"
	"jethur/utility"
	"os"
)

type CommonAIServiceInterface interface {
	RewriteDescription(dtos.DescriptionRewriteRequest) (map[string]interface{}, error)
	TranslateLanguage(dtos.LanguageTranslationRequest) (map[string]interface{}, error)
}

type CommonAIService struct{}

func NewCommonAIService() CommonAIServiceInterface {
	return &CommonAIService{}
}

func (s *CommonAIService) RewriteDescription(req dtos.DescriptionRewriteRequest) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("%s/api/v1/description-rewriter", os.Getenv("AI_SERVICE_BASE_URL"))
	return utility.CallAIAPI(apiURL, req)
}

func (s *CommonAIService) TranslateLanguage(req dtos.LanguageTranslationRequest) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("%s/api/v1/language-translation", os.Getenv("AI_SERVICE_BASE_URL"))
	return utility.CallAIAPI(apiURL, req)
}
