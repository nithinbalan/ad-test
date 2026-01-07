package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/jethur-ai/dtos"
	"jethur/grc/settings/services"
	"jethur/utility"
	"mime/multipart"
	"net/http"
)

type ComplianceAIServiceInterface interface {
	GetComplianceAutofill(req dtos.ComplianceAutofillRequest) (map[string]interface{}, error)
	ImportComplianceRegister(file *multipart.FileHeader) (map[string]interface{}, error)
	GetComplianceSummary(complianceID uint) (map[string]interface{}, error)
	GetComplianceClauseGuidance(clauseID uint) (map[string]interface{}, error)
	GetComplianceStrategicObjectivesJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceSuggestKPIJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceSuggestProcessJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceSuggestRisksJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceSuggestControlsJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceSuggestAssociatesJob(complianceID uint) (map[string]interface{}, error)
	GetComplianceJobStatus(jobID string) (map[string]interface{}, error)

	GetComplianceClauseStrategicObjectivesJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestKPIJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestProcessJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestRisksJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestControlsJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestIncidentsJob(clauseID uint) (map[string]interface{}, error)
	GetComplianceClauseSuggestAssociatesJob(clauseID uint) (map[string]interface{}, error)
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

func (cais *ComplianceAIService) ImportComplianceRegister(
	fileHeader *multipart.FileHeader,
) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	writer.Close()

	apiURL, err := utility.GetAIEndpoint("/api/v1/compliance/compliance-register-import")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("AI API error: %s", string(b))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (cais *ComplianceAIService) GetComplianceSummary(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/compliance-summary", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseGuidance(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/clause-guidance", clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceStrategicObjectivesJob(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/strategic-objectives", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceSuggestKPIJob(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/suggest-kpi", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceSuggestProcessJob(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/suggest-process", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceSuggestRisksJob(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/suggest-risks", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceSuggestControlsJob(complianceID uint) (map[string]interface{}, error) {
	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/suggest-controls", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceSuggestAssociatesJob(complianceID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/%d/start-job/suggest-associates", complianceID)

	return utility.CallAIAPIObjectResponse(path)
}

func (ajs *ComplianceAIService) GetComplianceJobStatus(jobID string) (map[string]interface{}, error) {

	systemPreference, err := ajs.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/job-status/%s", jobID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseStrategicObjectivesJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}

	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf(
		"/api/v1/compliance/clause/%d/start-job/strategic-objectives",
		clauseID,
	)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestKPIJob(clauseID uint) (map[string]interface{}, error) {

    systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
    if err != nil {
        return nil, err
    }

    if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
        return nil, fmt.Errorf("AI features are disabled in system preferences")
    }

    path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-kpi",clauseID)

    return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestProcessJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-process",clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestRisksJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-risks",clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestControlsJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-controls",clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestIncidentsJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-incidents",clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

func (cais *ComplianceAIService) GetComplianceClauseSuggestAssociatesJob(clauseID uint) (map[string]interface{}, error) {

	systemPreference, err := cais.systemPreferenceSvc.GetSystemPreference()
	if err != nil {
		return nil, err
	}
	if systemPreference.IsAIEnable == nil || !*systemPreference.IsAIEnable {
		return nil, fmt.Errorf("AI features are disabled in system preferences")
	}

	path := fmt.Sprintf("/api/v1/compliance/clause/%d/start-job/suggest-associates",clauseID)

	return utility.CallAIAPIObjectResponse(path)
}

