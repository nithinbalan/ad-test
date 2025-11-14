package services

import (
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"math"
	"net/http"
	"strings"
	"time"
)

type TenableServiceInterface interface {
	GetAllAssets(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error)
	ValidateTenableAssets(input []dtos.TenableAssetImportDTO, currentUserID uint) ([]map[string]interface{}, error)
	ImportAssets(input []dtos.TenableAssetImportDTO, currentUserID uint) ([]uint, error)
}

type TenableService struct {
	repo                repositories.TenableRepositoryInterface
	ass                 AssetStatusServiceInterface
	assetCategorySvc    AssetCategoryServiceInterface
	userSvc             organizationServices.UserServiceInterface
	assetRegisterSvc    AssetRegisterServiceInterface
	assetSubCategorySvc AssetSubCategoryServiceInterface
}

func NewTenableService(
	repo repositories.TenableRepositoryInterface,
	ass AssetStatusServiceInterface,
	assetCategorySvc AssetCategoryServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	assetSubCategorySvc AssetSubCategoryServiceInterface,
) TenableServiceInterface {
	return &TenableService{
		repo:                repo,
		ass:                 ass,
		assetCategorySvc:    assetCategorySvc,
		userSvc:             userSvc,
		assetRegisterSvc:    assetRegisterSvc,
		assetSubCategorySvc: assetSubCategorySvc,
	}
}

const (
	BaseURL   = "https://sc.tenalab.online"
	AccessKey = "803fe2c0529f448983c9168df30fd801"
	SecretKey = "17511bc49b254f2283e20de69f7cce6d"
)

func (s *TenableService) GetAllAssets(params *dtos.TenableAssetListQueryParamsDTO) (*dtos.TenableAssetListResponseDTO, error) {

	client := &http.Client{Timeout: 20 * time.Second}
	url := fmt.Sprintf("%s/rest/asset?fields=uuid,name,description,status,createdTime,modifiedTime", BaseURL)

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

	strPtr := func(s string) *string { return &s }

	nonNil := func(val interface{}) *string {
		if val == nil || fmt.Sprintf("%v", val) == "<nil>" {
			empty := ""
			return &empty
		}
		s := fmt.Sprintf("%v", val)
		return &s
	}

	var items []dtos.TenableAssetDTO
	for _, item := range allAssets[start:end] {
		obj := item.(map[string]interface{})

		// Determine status
		var statusType string
		if fmt.Sprintf("%v", obj["status"]) == "0" {
			statusType = "active"
		} else {
			statusType = "disposed"
		}

		statusObj, err := s.ass.GetByType(statusType)
		if err != nil {
			return nil, fmt.Errorf("failed to get asset status '%s': %v", statusType, err)
		}

		items = append(items, dtos.TenableAssetDTO{
			TenableUUID:        fmt.Sprintf("%v", obj["uuid"]),
			TitleEn:            fmt.Sprintf("%v", obj["name"]),
			TitleAr:            fmt.Sprintf("%v", obj["name"]),
			AssetCategoryID:    0,
			AssetCategory:      nonNil(obj["asset_category"]),
			AssetSubCategoryID: 0,
			AssetSubCategory:   nonNil(obj["asset_sub_category"]),
			AssetStatusID:      statusObj.ID,
			AssetStatus:        &statusObj.TitleEn,
			ResponsibleUserID:  0,
			ResponsibleUser:    nonNil(obj["responsible_user"]),
			Location:           nonNil(obj["location"]),
			Description:        nonNil(obj["description"]),
			AcquisitionDate:    nonNil(obj["acquisition_date"]),
			AssignmentDate:     nonNil(obj["assignment_date"]),
			WarrantyExpiryDate: nonNil(obj["warranty_expiry_date"]),
			Manufacturer:       strPtr("Tenable Security Center"),
			Model:              nonNil(obj["model"]),
			SerialNumber:       nonNil(obj["serial_number"]),
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

func (s *TenableService) ValidateTenableAssets(input []dtos.TenableAssetImportDTO, currentUserID uint) ([]map[string]interface{}, error) {
	var validatedData []map[string]interface{}

	categories, err := s.assetCategorySvc.GetAll(&dtos.AssetCategoryListQueryParamsDTO{Status: "active"})
	if err != nil {
		return nil, fmt.Errorf("error fetching asset categories: %v", err)
	}
	categoryTitles := make(map[string]uint)
	for _, cat := range categories {
		categoryTitles[strings.ToLower(strings.TrimSpace(cat.TitleEn))] = uint(cat.ID)
	}

	subCategories, err := s.assetSubCategorySvc.GetAll(&dtos.AssetSubCategoryListQueryParamsDTO{Status: "active"})
	if err != nil {
		return nil, fmt.Errorf("error fetching asset subcategories: %v", err)
	}
	subCategoryTitles := make(map[string]uint)
	for _, sub := range subCategories {
		subCategoryTitles[strings.ToLower(strings.TrimSpace(sub.TitleEn))] = uint(sub.ID)
	}

	statuses, err := s.ass.GetAll("", "", "", "")
	if err != nil {
		return nil, fmt.Errorf("error fetching asset statuses: %v", err)
	}
	statusTitles := make(map[string]uint)
	for _, st := range statuses {
		statusTitles[strings.ToLower(strings.TrimSpace(st.TitleEn))] = st.ID
	}

	users, err := s.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.ToLower(strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN))
		userFullNames[fullName] = uint(user.ID)
	}

	for _, record := range input {
		rowErrors := make(map[string]string)
		processed := make(map[string]interface{})

		fields := []struct {
			key   string
			value interface{}
		}{
			{"tenable_uuid", record.TenableUUID},
			{"title_en", record.TitleEn},
			{"title_ar", record.TitleAr},
			{"asset_category", record.AssetCategory},
			{"asset_sub_category", record.AssetSubCategory},
			{"asset_status", record.AssetStatus},
			{"responsible_user", record.ResponsibleUser},
			{"location", record.Location},
			{"description", record.Description},
			{"acquisition_date", record.AcquisitionDate},
			{"assignment_date", record.AssignmentDate},
			{"warranty_expiry_date", record.WarrantyExpiryDate},
			{"manufacturer", record.Manufacturer},
			{"model", record.Model},
			{"serial_number", record.SerialNumber},
		}

		for _, f := range fields {
			processed[f.key] = f.value
		}

		if record.TitleEn == "" {
			rowErrors["title_en"] = "Title (EN) is required"
		}
		if record.TitleAr == "" {
			rowErrors["title_ar"] = "Title (AR) is required"
		}

		if strings.TrimSpace(record.AssetCategory) == "" {
			rowErrors["asset_category"] = "Asset category is required"
		} else if id, ok := categoryTitles[strings.ToLower(strings.TrimSpace(record.AssetCategory))]; ok {
			processed["asset_category_id"] = id
		} else {
			rowErrors["asset_category"] = fmt.Sprintf("Invalid asset category: %s", record.AssetCategory)
		}

		if strings.TrimSpace(record.AssetSubCategory) == "" {
			rowErrors["asset_sub_category"] = "Asset subcategory is required"
		} else if id, ok := subCategoryTitles[strings.ToLower(strings.TrimSpace(record.AssetSubCategory))]; ok {
			processed["asset_sub_category_id"] = id
		} else {
			rowErrors["asset_sub_category"] = fmt.Sprintf("Invalid asset subcategory: %s", record.AssetSubCategory)
		}

		if strings.TrimSpace(record.AssetStatus) == "" {
			rowErrors["asset_status"] = "Asset status is required"
		} else if id, ok := statusTitles[strings.ToLower(strings.TrimSpace(record.AssetStatus))]; ok {
			processed["asset_status_id"] = id
		} else {
			rowErrors["asset_status"] = fmt.Sprintf("Invalid asset status: %s", record.AssetStatus)
		}

		if strings.TrimSpace(record.ResponsibleUser) == "" {
			rowErrors["responsible_user"] = "Responsible user is required"
		} else if id, ok := userFullNames[strings.ToLower(strings.TrimSpace(record.ResponsibleUser))]; ok {
			if _, err := s.userSvc.GetUserPrimaryUnitID(id); err != nil {
				rowErrors["responsible_user"] = "Organization unit not set for responsible user"
			} else {
				processed["responsible_user_id"] = id
			}
		} else {
			rowErrors["responsible_user"] = fmt.Sprintf("Invalid responsible user: %s", record.ResponsibleUser)
		}

		if len(rowErrors) > 0 {
			processed["errors"] = rowErrors
		} else {
			processed["errors"] = nil
		}

		validatedData = append(validatedData, processed)
	}

	return validatedData, nil
}

func (s *TenableService) ImportAssets(input []dtos.TenableAssetImportDTO, currentUserID uint) ([]uint, error) {
	validated, err := s.ValidateTenableAssets(input, currentUserID)
	if err != nil {
		return nil, err
	}

	for _, r := range validated {
		if errs, ok := r["errors"].(map[string]interface{}); ok && len(errs) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint

	for _, record := range validated {
		dto := &dtos.AssetRegisterFormDTO{
			TitleEn:            utility.GetStringValue(record["title_en"]),
			TitleAr:            utility.GetStringValue(record["title_ar"]),
			AssetCategoryID:    utility.GetUintValue(record["asset_category_id"]),
			AssetSubCategoryID: utility.GetUintValue(record["asset_sub_category_id"]),
			AssetStatusID:      utility.GetUintValue(record["asset_status_id"]),
			ResponsibleUserID:  utility.GetUintValue(record["responsible_user_id"]),
			Location:           utility.GetOptionalStringValue(record["location"]),
			Description:        utility.GetOptionalStringValue(record["description"]),
			AcquisitionDate:    utility.GetOptionalStringValue(record["acquisition_date"]),
			AssignmentDate:     utility.GetOptionalStringValue(record["assignment_date"]),
			WarrantyExpiryDate: utility.GetOptionalStringValue(record["warranty_expiry_date"]),
			Manufacturer:       utility.GetOptionalStringValue(record["manufacturer"]),
			Model:              utility.GetOptionalStringValue(record["model"]),
			SerialNumber:       utility.GetOptionalStringValue(record["serial_number"]),
			Documents:          []*commonDTOs.DocumentFormDTO{},
		}

		id, err := s.assetRegisterSvc.Create(dto, currentUserID)
		if err != nil {
			return nil, fmt.Errorf("failed to create asset for record: %+v, error: %v", record, err)
		}

		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}
