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
	ValidateTenableAssets(input []map[string]interface{}, currentUserID uint) ([]map[string]interface{}, error)
	ImportAssets(input []map[string]interface{}, currentUserID uint) ([]uint, error)
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
		var statusTitle string
		if fmt.Sprintf("%v", obj["status"]) == "0" {
			statusTitle = "Active"
		} else {
			statusTitle = "Disposed"
		}

		statusObj, err := s.ass.GetByTitleEn(statusTitle)
		if err != nil {
			return nil, fmt.Errorf("failed to get asset status '%s': %v", statusTitle, err)
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
			AssetStatus:        nonNil(obj["asset_status"]),
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
			Documents:          nil,
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

func (s *TenableService) ValidateTenableAssets(input []map[string]interface{}, currentUserID uint) ([]map[string]interface{}, error) {
	var validated []map[string]interface{}

	params := &dtos.AssetRegisterListQueryParamsDTO{
		Export:   true,
		SortBy:   "asset_register.id",
		Sort:     "DESC",
		Page:     1,
		PageSize: 1,
	}

	existingAssets, err := s.assetRegisterSvc.GetAll(params)
	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool)
	for _, asset := range existingAssets.Items {
		existingMap[strings.ToLower(asset.TitleEn)] = true
	}

	for _, record := range input {
		rowErrors := s.validateAssetRecord(record, existingMap)
		processed := make(map[string]interface{})

		jsonData, _ := json.Marshal(record)
		json.Unmarshal(jsonData, &processed)

		if len(rowErrors) > 0 {
			processed["errors"] = rowErrors
		} else {
			processed["errors"] = nil
		}

		validated = append(validated, processed)
	}

	return validated, nil
}

func (s *TenableService) validateAssetRecord(record map[string]interface{}, existingMap map[string]bool) map[string]interface{} {
	rowErrors := make(map[string]interface{})

	titleEn := utility.GetStringValue(record["title_en"])
	titleAr := utility.GetStringValue(record["title_ar"])
	categoryID := utility.GetUintValue(record["asset_category_id"])
	subCategoryID := utility.GetUintValue(record["asset_sub_category_id"])
	statusID := utility.GetUintValue(record["asset_status_id"])
	responsibleUserID := utility.GetUintValue(record["responsible_user_id"])

	if titleEn == "" {
		rowErrors["title_en"] = "Field is required"
	}
	if titleAr == "" {
		rowErrors["title_ar"] = "Field is required"
	}

	if _, exists := existingMap[strings.ToLower(titleEn)]; exists {
		rowErrors["title_en"] = "Asset already exists in the system"
	}

	if categoryID == 0 {
		rowErrors["asset_category_id"] = "Field is required"
	} else if category, err := s.assetCategorySvc.GetByID(categoryID); err != nil {
		rowErrors["asset_category_id"] = "Asset Category not found"
	} else {
		record["asset_category"] = category.TitleEn
	}

	if subCategoryID == 0 {
		rowErrors["asset_sub_category_id"] = "Field is required"
	} else if subCategory, err := s.assetSubCategorySvc.GetByID(categoryID); err != nil {
		rowErrors["asset_sub_category_id"] = "Asset Sub Category not found"
	} else {
		record["asset_sub_category"] = subCategory.TitleEn
	}

	if statusID == 0 {
		rowErrors["asset_status_id"] = "Field is required"
	} else if status, err := s.ass.GetByID(statusID); err != nil {
		rowErrors["asset_status_id"] = "Asset Status not found"
	} else {
		record["asset_status"] = status.TitleEn
	}

	if responsibleUserID == 0 {
		rowErrors["responsible_user_id"] = "Field is required"
	} else if user, err := s.userSvc.GetByID(responsibleUserID); err != nil {
		rowErrors["responsible_user_id"] = "Responsible User not found"
	} else {
		if _, err := s.userSvc.GetUserPrimaryUnitID(responsibleUserID); err != nil {
			rowErrors["responsible_user_id"] = "Organization unit not set for responsible user"
		} else {
			record["responsible_user"] = user.FirstNameEn + " " + user.LastNameEn
		}
	}

	for _, dateField := range []string{"acquisition_date", "assignment_date", "warranty_expiry_date"} {
		if val, ok := record[dateField]; ok && val != nil {
			if strVal := utility.GetStringValue(val); strVal != "" {
				if _, err := utility.ParseDate(strVal); err != nil {
					rowErrors[dateField] = "Invalid date format (expected YYYY-MM-DD)"
				}
			}
		}
	}

	return rowErrors
}

func (s *TenableService) ImportAssets(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validated, err := s.ValidateTenableAssets(input, currentUserID)
	if err != nil {
		return nil, err
	}

	// Check if any record has validation errors
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
