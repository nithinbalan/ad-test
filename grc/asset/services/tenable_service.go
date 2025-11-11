package services

import (
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/utility"
	"math"
	"net/http"
	"time"
)

type TenableServiceInterface interface {
	GetAllAssets(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error)
}

type TenableService struct {
	repo repositories.TenableRepositoryInterface
}

func NewTenableService(repo repositories.TenableRepositoryInterface) TenableServiceInterface {
	return &TenableService{repo: repo}
}

const (
	BaseURL   = "https://sc.tenalab.online"
	AccessKey = "803fe2c0529f448983c9168df30fd801"
	SecretKey = "17511bc49b254f2283e20de69f7cce6d"
)

func (s *TenableService) GetAllAssets(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error) {

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/asset?fields=id,name,description,status,createdTime,modifiedTime", BaseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-apikey", fmt.Sprintf("accessKey=%s; secretKey=%s", AccessKey, SecretKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	response, ok := data["response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format from Tenable SC")
	}

	usable, _ := response["usable"].([]interface{})
	manageable, _ := response["manageable"].([]interface{})
	allAssets := append(usable, manageable...)

	total := len(allAssets)
	start := (params.Page - 1) * params.PageSize
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	if start >= total {
		return &dtos.TenableAssetListResponseDTO{
			Items:      []dtos.TenableAssetDTO{},
			Pagination: commonDTOs.PaginationInfoDTO{TotalCount: int64(total), CurrentPage: params.Page, PageSize: params.PageSize, TotalPages: totalPages},
		}, nil
	}

	end := start + params.PageSize
	if end > total {
		end = total
	}

	var items []dtos.TenableAssetDTO
	for _, item := range allAssets[start:end] {
		obj := item.(map[string]interface{})

		status := "DISABLED"
		if fmt.Sprintf("%v", obj["status"]) == "0" {
			status = "ACTIVE"
		}

		items = append(items, dtos.TenableAssetDTO{
			ReferenceNumber: fmt.Sprintf("%v", obj["id"]),
			TitleEn:         fmt.Sprintf("%v", obj["name"]),
			Description:     fmt.Sprintf("%v", obj["description"]),
			StatusName:      status,
			Manufacturer:    "Tenable Security Center",
			CreatedAt:       utility.ParseUnixTimestamp(obj["createdTime"]),
			UpdatedAt:       utility.ParseUnixTimestamp(obj["modifiedTime"]),
		})
	}

	return &dtos.TenableAssetListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			TotalCount:  int64(total),
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
		},
	}, nil
}
