package repositories

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type TenableRepositoryInterface interface {
	FetchAllAssets() (map[string]interface{}, error)
}

type TenableRepository struct{}

func NewTenableRepository() TenableRepositoryInterface {
	return &TenableRepository{}
}

const (
	BaseURL   = "https://sc.tenalab.online"
	AccessKey = "803fe2c0529f448983c9168df30fd801"
	SecretKey = "17511bc49b254f2283e20de69f7cce6d"
)

func (repo *TenableRepository) FetchAllAssets() (map[string]interface{}, error) {
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

	return data, nil
}
