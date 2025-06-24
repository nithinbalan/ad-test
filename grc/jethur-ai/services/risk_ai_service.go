package services

import (
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/jethur-ai/dtos"
	"jethur/utility"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RiskAIServiceInterface interface {
	GetRiskAutofill(req dtos.RiskAutofillRequest) (map[string]interface{}, error)
	GetRiskRCA(req dtos.RiskRCARequest) ([]dtos.RCAItem, error)
	GetRiskAnalysis(req dtos.RiskAIRequest) (map[string]interface{}, error)
	GetRiskResponseStrategyAnalysis(req dtos.RiskAIRequest) (map[string]interface{}, error)
	GetRiskControlAnalysis(riskId uint) (map[string]interface{}, error)
	GetRiskKriAnalysis(riskId uint) ([]map[string]interface{}, error)
	GetRiskResidualControlAnalysis(riskId uint) (map[string]interface{}, error)
}

type RiskAIService struct{}

func NewRiskAIService() RiskAIServiceInterface {
	return &RiskAIService{}
}

func (s *RiskAIService) GetRiskAutofill(req dtos.RiskAutofillRequest) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"risk_title":      req.RiskTitle,
		"risk_categories": req.RiskCategories,
		"risk_sources":    req.RiskSources,
	}

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/risk-autofill", baseURL)

	return utility.CallAIAPI(apiURL, payload)
}

func (s *RiskAIService) GetRiskRCA(req dtos.RiskRCARequest) ([]dtos.RCAItem, error) {
	form := make(url.Values)
	form.Set("risk_title", req.RiskTitle)
	form.Set("consequences", req.Consequences)
	form.Set("category", req.Category)
	form.Set("source", req.Source)

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}
	apiURL := fmt.Sprintf("%s/api/v1/risk/rca", baseURL)

	reqBody := strings.NewReader(form.Encode())
	httpReq, err := http.NewRequest("POST", apiURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI API returned %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result []dtos.RCAItem
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI response: %w", err)
	}

	return result, nil
}

func (s *RiskAIService) GetRiskAnalysis(req dtos.RiskAIRequest) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"risk_id":        req.RiskId,
		"register_title": req.RegisterTitle,
		"description":    req.Description,
		"category":       req.Category,
		"source":         req.Source,
	}

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/risk-analysis", baseURL)

	return utility.CallAIAPI(apiURL, payload)
}

func (s *RiskAIService) GetRiskResponseStrategyAnalysis(req dtos.RiskAIRequest) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"risk_id":        req.RiskId,
		"register_title": req.RegisterTitle,
		"description":    req.Description,
		"category":       req.Category,
		"source":         req.Source,
	}

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/risk-response-strategy-analysis", baseURL)
	return utility.CallAIAPI(apiURL, payload)
}

func (s *RiskAIService) GetRiskControlAnalysis(riskID uint) (map[string]interface{}, error) {

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/%d/control-summary", baseURL, riskID)
	fmt.Println(apiURL)

	return utility.CallAIAPIObjectResponse(apiURL)
}

func (s *RiskAIService) GetRiskKriAnalysis(riskID uint) ([]map[string]interface{}, error) {

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/%d/risk-kri", baseURL, riskID)
	fmt.Println(apiURL)

	return utility.CallAIAPIArrayResponse(apiURL)
}

func (s *RiskAIService) GetRiskResidualControlAnalysis(riskID uint) (map[string]interface{}, error) {

	baseURL := os.Getenv("AI_SERVICE_BASE_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("AI_SERVICE_BASE_URL is not set")
	}

	apiURL := fmt.Sprintf("%s/api/v1/risk/%d/risk-residual-control-analysis", baseURL, riskID)
	fmt.Println(apiURL)

	return utility.CallAIAPIObjectResponse(apiURL)
}
