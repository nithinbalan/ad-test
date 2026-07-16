package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"net/http"
	"strings"
	"time"
)

type QualysServiceInterface interface {
	GetKBDetailsByQID(qid string) (*dtos.QualysKBDetailsDTO, error)
}

type QualysService struct {
	repo                  repositories.QualysRepositoryInterface
	qualysSettingsService settingsServices.QualysSettingsServiceInterface
}

func NewQualysService(
	repo repositories.QualysRepositoryInterface,
	qualysSettingsService settingsServices.QualysSettingsServiceInterface,
) QualysServiceInterface {
	return &QualysService{repo: repo, qualysSettingsService: qualysSettingsService}
}

func (s *QualysService) GetKBDetailsByQID(qid string) (*dtos.QualysKBDetailsDTO, error) {
	settings, err := s.qualysSettingsService.GetQualysSettings()
	if err != nil {
		return nil, err
	}
	if settings == nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("no Qualys settings found"))
	}

	kbDetails, err := s.fetchKBDetails(qid, settings)
	if err != nil {
		return nil, err
	}

	detections, err := s.fetchDetectionsByQID(qid, settings)
	if err != nil {
		return nil, err
	}
	kbDetails.Detections = detections

	return kbDetails, nil
}

func (s *QualysService) fetchKBDetails(qid string, settings *models.QualysSettings) (*dtos.QualysKBDetailsDTO, error) {
	client := &http.Client{Timeout: 300 * time.Second}

	url := fmt.Sprintf("%s/api/4.0/fo/knowledge_base/vuln/", strings.TrimRight(settings.PlatformURL, "/"))
	payload := strings.NewReader(fmt.Sprintf("action=list&ids=%s&details=All", qid))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("unable to create request: %w", err))
	}

	req.SetBasicAuth(settings.Username, settings.Password)
	req.Header.Set("X-Requested-With", "GRC-KB")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("connection failed: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("failed to read response: %w", err))
	}

	if resp.StatusCode != 200 {
		var xmlResp struct {
			Response struct {
				Code string `xml:"CODE"`
				Text string `xml:"TEXT"`
			} `xml:"RESPONSE"`
		}
		if err := xml.Unmarshal(body, &xmlResp); err == nil && strings.TrimSpace(xmlResp.Response.Text) != "" {
			return nil, utility.NewQualysConnectionError(fmt.Errorf("%s", strings.TrimSpace(xmlResp.Response.Text)))
		}
		return nil, utility.NewQualysConnectionError(fmt.Errorf("Qualys KB request failed (HTTP %d)", resp.StatusCode))
	}

	var parsed dtos.QualysKBResponse
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("failed to parse Qualys KB response: %w", err))
	}

	if len(parsed.Response.Vulns) == 0 {
		return nil, utility.NewResourceNotFoundError(fmt.Errorf("no KB entry found for QID %s", qid), qid)
	}

	vuln := parsed.Response.Vulns[0]
	cveIDs := make([]string, len(vuln.CVEList))
	for i, cve := range vuln.CVEList {
		cveIDs[i] = cve.ID
	}

	return &dtos.QualysKBDetailsDTO{
		QID:         vuln.QID,
		Title:       vuln.Title,
		Diagnosis:   vuln.Diagnosis,
		Category:    vuln.Category,
		Consequence: vuln.Consequence,
		Solution:    vuln.Solution,
		CVSS:        vuln.CVSS3,
		CVEList:     cveIDs,
	}, nil
}

func (s *QualysService) fetchDetectionsByQID(qid string, settings *models.QualysSettings) ([]dtos.QualysDetectionDTO, error) {
	client := &http.Client{Timeout: 300 * time.Second}

	url := fmt.Sprintf("%s/api/4.0/fo/asset/host/vm/detection/", strings.TrimRight(settings.PlatformURL, "/"))
	payload := strings.NewReader(fmt.Sprintf("action=list&qids=%s&show_qds=1&show_tags=1&output_format=XML", qid))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("unable to create request: %w", err))
	}

	req.SetBasicAuth(settings.Username, settings.Password)
	req.Header.Set("X-Requested-With", "GRC-Detections")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("connection failed: %w", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("failed to read response: %w", err))
	}

	if resp.StatusCode != 200 {
		var xmlResp struct {
			Response struct {
				Code string `xml:"CODE"`
				Text string `xml:"TEXT"`
			} `xml:"RESPONSE"`
		}
		if err := xml.Unmarshal(body, &xmlResp); err == nil && strings.TrimSpace(xmlResp.Response.Text) != "" {
			return nil, utility.NewQualysConnectionError(fmt.Errorf("%s", strings.TrimSpace(xmlResp.Response.Text)))
		}
		return nil, utility.NewQualysConnectionError(fmt.Errorf("Qualys detection request failed (HTTP %d)", resp.StatusCode))
	}

	var parsed dtos.QualysDetectionResponse
	if err := xml.Unmarshal(body, &parsed); err != nil {
		return nil, utility.NewQualysConnectionError(fmt.Errorf("failed to parse Qualys detection response: %w", err))
	}

	var result []dtos.QualysDetectionDTO
	for _, host := range parsed.Response.HostList.Hosts {
		for _, d := range host.Detections {
			result = append(result, dtos.QualysDetectionDTO{
				HostID:           host.ID,
				IP:               host.IP,
				IPv6:             host.IPv6,
				DNS:              host.DNS,
				QID:              d.QID,
				Severity:         d.Severity,
				Status:           d.Status,
				QDS:              strings.TrimSpace(d.QDS.Value),
				CriticalityScore: d.CriticalityScore,
				Port:             d.Port,
				Protocol:         d.Protocol,
				LastFound:        d.LastFound,
				Tags:             d.Tags,
			})
		}
	}

	return result, nil
}
