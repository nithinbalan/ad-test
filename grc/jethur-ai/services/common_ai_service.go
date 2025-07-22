package services

import (
	"jethur/grc/jethur-ai/dtos"
	"jethur/utility"
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
	return utility.CallAIAPI("/api/v1/description-rewriter", req)
}

func (s *CommonAIService) TranslateLanguage(req dtos.LanguageTranslationRequest) (map[string]interface{}, error) {
	return utility.CallAIAPI("/api/v1/language-translation", req)
}
