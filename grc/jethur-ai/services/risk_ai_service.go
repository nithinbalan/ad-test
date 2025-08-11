package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"net/url"
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
	GetRiskResidualRiskAnalysis(riskId uint) (map[string]interface{}, error)
	ChatWithRisk(ctx context.Context, req dtos.RiskAIChatRequest, currentUserID uint) (<-chan string, error)
	GetRiskChatsByContext(contextID string, offset, limit int) ([]map[string]interface{}, error)
	GetDuplicateRisks() ([]map[string]interface{}, error)
	GetRiskClosureAutofill(riskID uint) (map[string]interface{}, error)
}

type RiskAIService struct {
	systemPreferenceSvc services.SystemPreferenceServiceInterface
}

func NewRiskAIService(systemPreferenceSvc services.SystemPreferenceServiceInterface) RiskAIServiceInterface {
	return &RiskAIService{
		systemPreferenceSvc: systemPreferenceSvc,
	}
}

func (s *RiskAIService) GetRiskAutofill(req dtos.RiskAutofillRequest) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"risk_title":      req.RiskTitle,
		"risk_categories": req.RiskCategories,
		"risk_sources":    req.RiskSources,
	}

	return utility.CallAIAPI("/api/v1/risk/risk-autofill", payload)
}

func (s *RiskAIService) GetRiskRCA(req dtos.RiskRCARequest) ([]dtos.RCAItem, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	form := make(url.Values)
	form.Set("risk_title", req.RiskTitle)
	form.Set("consequences", req.Consequences)
	form.Set("category", req.Category)
	form.Set("source", req.Source)

	apiURL, err := utility.GetAIEndpoint("/api/v1/risk/rca")
	if err != nil {
		return nil, err
	}

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
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"risk_id":        req.RiskId,
		"register_title": req.RegisterTitle,
		"description":    req.Description,
		"category":       req.Category,
		"source":         req.Source,
	}

	return utility.CallAIAPI("/api/v1/risk/risk-analysis", payload)
}

func (s *RiskAIService) GetRiskResponseStrategyAnalysis(req dtos.RiskAIRequest) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"risk_id":        req.RiskId,
		"register_title": req.RegisterTitle,
		"description":    req.Description,
		"category":       req.Category,
		"source":         req.Source,
	}

	return utility.CallAIAPI("/api/v1/risk/risk-response-strategy-analysis", payload)
}

func (s *RiskAIService) GetRiskControlAnalysis(riskID uint) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/%d/control-summary", riskID)
	return utility.CallAIAPIObjectResponse(path)
}

func (s *RiskAIService) GetRiskKriAnalysis(riskID uint) ([]map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/%d/risk-kri", riskID)
	return utility.CallAIAPIArrayResponse(path)
}

func (s *RiskAIService) GetRiskResidualControlAnalysis(riskID uint) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/%d/risk-residual-control-analysis", riskID)
	return utility.CallAIAPIObjectResponse(path)
}

func (s *RiskAIService) GetRiskResidualRiskAnalysis(riskID uint) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/%d/residual-risk-analysis", riskID)
	return utility.CallAIAPIObjectResponse(path)
}

func (s *RiskAIService) ChatWithRisk(ctx context.Context, req dtos.RiskAIChatRequest, currentUserID uint) (<-chan string, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	payload := map[string]interface{}{
		"content":         req.Content,
		"role":            req.Role,
		"current_user_id": fmt.Sprintf("%d", currentUserID),
		"context_id":      req.ContextID,
	}

	return utility.StreamAIAPI(ctx, "/api/v1/risk/risk-chat", payload)
}

func (s *RiskAIService) GetRiskChatsByContext(contextID string, offset, limit int) ([]map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/get/%s?offset=%d&limit=%d", contextID, offset, limit)
	return utility.CallAIAPIArrayResponse(path)
}

func (s *RiskAIService) GetDuplicateRisks() ([]map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := "/api/v1/risk/get/get-duplicate-risks"
	return utility.CallAIAPIArrayResponse(path)
}

func (s *RiskAIService) GetRiskClosureAutofill(riskID uint) (map[string]interface{}, error) {
	systemPreference, err := s.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/risk/%d/risk-closure-auto-fill", riskID)
	return utility.CallAIAPIObjectResponse(path)
}
