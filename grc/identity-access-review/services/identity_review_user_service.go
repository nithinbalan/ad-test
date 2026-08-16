package services

import (
	"errors"
	"fmt"
	"io"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IdentityReviewUserServiceInterface interface {
	GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error)
	GetByID(id uint) (*models.IdentityReviewUser, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error)
	GetAccountDecisionHistory(id uint) ([]*dtos.IdentityReviewUserAccountDecisionHistoryDTO, error)
	Create(currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*dtos.IdentityReviewUserDetailDTO, error)
	Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*dtos.IdentityReviewUserDetailDTO, error)
	UpdateAccountDecision(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserAccountDecisionUpdateFormDTO, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error)
	UpdatePrivilege(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserPrivilegeUpdateFormDTO, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GenerateIdentityReviewUserTemplate(currentUserID uint) (*os.File, error)
	ValidateIdentityReviewUserTemplate(file io.Reader, identityReviewCampaignID uint, currentUserID uint) ([]map[string]interface{}, error)
	ValidateImportIdentityReviewUsers(input []map[string]interface{}, identityReviewCampaignID uint, currentUserID uint) ([]map[string]interface{}, error)
	ImportIdentityReviewUsers(input []map[string]interface{}, identityReviewCampaignID uint, currentUserID uint) ([]uint, error)
}

type IdentityReviewUserService struct {
	repo                      repositories.IdentityReviewUserRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	organizationUnitSvc       settingsServices.OrganizationUnitServiceInterface
	privilegedAccountSvc      PrivilegedAccountServiceInterface
	accountTypeSvc            AccountTypeServiceInterface
	accountDecisionSvc        AccountDecisionServiceInterface
	identityReviewCampaignSvc IdentityReviewCampaignServiceInterface
}

func NewIdentityReviewUserService(
	repo repositories.IdentityReviewUserRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	privilegedAccountSvc PrivilegedAccountServiceInterface,
	accountTypeSvc AccountTypeServiceInterface,
	accountDecisionSvc AccountDecisionServiceInterface,
	identityReviewCampaignSvc IdentityReviewCampaignServiceInterface,
) IdentityReviewUserServiceInterface {
	return &IdentityReviewUserService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		organizationUnitSvc:       organizationUnitSvc,
		privilegedAccountSvc:      privilegedAccountSvc,
		accountTypeSvc:            accountTypeSvc,
		accountDecisionSvc:        accountDecisionSvc,
		identityReviewCampaignSvc: identityReviewCampaignSvc,
	}
}

func (irus *IdentityReviewUserService) GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error) {
	return irus.repo.GetAll(params)
}

func (irus *IdentityReviewUserService) GetByID(id uint) (*models.IdentityReviewUser, error) {
	identityReviewUser, err := irus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return identityReviewUser, nil
}

func (irus *IdentityReviewUserService) GetAccountDecisionHistory(id uint) ([]*dtos.IdentityReviewUserAccountDecisionHistoryDTO, error) {
	if _, err := irus.GetByID(id); err != nil {
		return nil, err
	}

	return irus.repo.GetAccountDecisionHistory(id)
}

func (irus *IdentityReviewUserService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error) {
	var identityReviewUserDetailDTO dtos.IdentityReviewUserDetailDTO

	identityReviewUser, err := irus.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(identityReviewUser, &identityReviewUserDetailDTO); err != nil {
		return nil, err
	}

	identityReviewCampaign, err := irus.identityReviewCampaignSvc.GetFullDetailsByID(
		identityReviewUser.IdentityReviewCampaignID,
		tx...,
	)
	if err != nil {
		return nil, err
	}
	identityReviewUserDetailDTO.IdentityReviewCampaign = identityReviewCampaign

	department, err := irus.organizationUnitSvc.GetFullDetailsByID(identityReviewUser.DepartmentID)
	if err != nil {
		return nil, err
	}
	identityReviewUserDetailDTO.Department = department

	identityReviewUserDetailDTO.CreatedAt = identityReviewUser.CreatedAt
	identityReviewUserDetailDTO.UpdatedAt = identityReviewUser.UpdatedAt

	return &identityReviewUserDetailDTO, nil
}

func (irus *IdentityReviewUserService) Create(currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*dtos.IdentityReviewUserDetailDTO, error) {
	if err := irus.validateForm(formDTO); err != nil {
		return nil, err
	}

	tx := irus.repo.BeginTransaction()

	identityReviewUser := &models.IdentityReviewUser{
		IdentityReviewCampaignID: formDTO.IdentityReviewCampaignID,
		Username:                 formDTO.Username,
		UserIdentifier:           &formDTO.UserIdentifier,
		LastLogin:                &formDTO.LastLogin,
		CurrentRole:              &formDTO.CurrentRole,
		DepartmentID:             formDTO.DepartmentID,
		PrivilegeID:              formDTO.PrivilegeID,
		AccountTypeID:            formDTO.AccountTypeID,
		CreatedByID:              currentUserID,
	}

	id, err := irus.repo.Create(identityReviewUser, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := irus.identityReviewCampaignSvc.CalculateAndUpdateCompletionPercentage(formDTO.IdentityReviewCampaignID, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	createdUser, err := irus.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			createdUser,
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"user",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return createdUser, nil
}

func (irus *IdentityReviewUserService) Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserFormDTO) (*dtos.IdentityReviewUserDetailDTO, error) {
	existingUser, err := irus.GetByID(id)
	if err != nil {
		return nil, err
	}

	oldUser, err := irus.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	if err := irus.validateForm(formDTO); err != nil {
		return nil, err
	}
	if existingUser.IdentityReviewCampaignID != formDTO.IdentityReviewCampaignID {
		return nil, utility.NewValidationError(
			fmt.Errorf("identity review campaign cannot be changed"),
			"immutable",
			"identity_review_campaign_id",
		)
	}

	existingUser.Username = formDTO.Username
	existingUser.UserIdentifier = &formDTO.UserIdentifier
	existingUser.LastLogin = &formDTO.LastLogin
	existingUser.CurrentRole = &formDTO.CurrentRole
	existingUser.DepartmentID = formDTO.DepartmentID
	existingUser.PrivilegeID = formDTO.PrivilegeID
	existingUser.AccountTypeID = formDTO.AccountTypeID
	existingUser.UpdatedByID = &currentUserID

	if err := irus.repo.Update(existingUser); err != nil {
		return nil, err
	}

	updatedUser, err := irus.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			map[string]interface{}{"old": oldUser, "new": updatedUser},
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"user",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return updatedUser, nil
}

func (irus *IdentityReviewUserService) UpdateAccountDecision(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserAccountDecisionUpdateFormDTO, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error) {
	databaseTx := utility.DatabaseTransaction(nil)
	ownsTransaction := len(tx) == 0 || tx[0] == nil
	if ownsTransaction {
		databaseTx = irus.repo.BeginTransaction()
	} else {
		databaseTx = tx[0]
	}

	existingUser, err := irus.repo.GetFullDetailsByID(id, databaseTx)
	if err != nil {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	oldUser := *existingUser
	decisionChanged := existingUser.AccountDecisionID == nil || *existingUser.AccountDecisionID != formDTO.AccountDecisionID

	if existingUser.IdentityReviewCampaignID == 0 {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		return nil, utility.NewValidationError(
			fmt.Errorf("identity review user is not assigned to a campaign"),
			"required",
			"identity_review_campaign_id",
		)
	}

	if _, err := irus.accountDecisionSvc.GetByID(formDTO.AccountDecisionID); err != nil {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		return nil, utility.NewValidationError(
			fmt.Errorf("account decision does not exist"),
			"not_exist",
			"account_decision_id",
		)
	}

	now := time.Now()
	existingUser.AccountDecisionID = &formDTO.AccountDecisionID
	existingUser.UpdatedByID = &currentUserID
	existingUser.UpdatedAt = &now

	if err := irus.repo.UpdateAccountDecision(existingUser, databaseTx); err != nil {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		return nil, err
	}

	if decisionChanged {
		history := &models.IdentityReviewUserAccountDecisionHistory{
			IdentityReviewUserID: id,
			AccountDecisionID:    formDTO.AccountDecisionID,
			CreatedByID:          currentUserID,
		}
		if err := irus.repo.CreateAccountDecisionHistory(history, databaseTx); err != nil {
			if ownsTransaction {
				databaseTx.Rollback()
			}
			return nil, err
		}
	}

	if err := irus.identityReviewCampaignSvc.CalculateAndUpdateCompletionPercentage(existingUser.IdentityReviewCampaignID, databaseTx); err != nil {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		return nil, err
	}

	updatedUser, err := irus.GetFullDetailsByID(id, databaseTx)
	if err != nil {
		if ownsTransaction {
			databaseTx.Rollback()
		}
		return nil, err
	}

	if ownsTransaction {
		if err := databaseTx.Commit(); err != nil {
			return nil, err
		}
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			map[string]interface{}{"old": &oldUser, "new": updatedUser},
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"account_decision",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()

	return updatedUser, nil
}

func (irus *IdentityReviewUserService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedUsers []interface{}

	for _, id := range ids {
		identityReviewUser, err := irus.GetFullDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if identityReviewUser.IdentityReviewCampaignID == 0 {
			failedIDs = append(failedIDs, id)
			continue
		}

		tx := irus.repo.BeginTransaction()
		if err := irus.repo.DeleteAccountDecisionHistory(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := irus.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := irus.identityReviewCampaignSvc.CalculateAndUpdateCompletionPercentage(identityReviewUser.IdentityReviewCampaignID, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedUsers = append(deletedUsers, identityReviewUser)
	}

	if len(successfulIDs) > 0 {
		go func() {
			if err := irus.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewUser{}.TableName(),
				successfulIDs,
				deletedUsers,
				currentUserID,
				nil,
			); err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}

	return failedIDs, nil
}

func (irus *IdentityReviewUserService) validateForm(formDTO *dtos.IdentityReviewUserFormDTO) error {
	if _, err := irus.identityReviewCampaignSvc.GetByID(formDTO.IdentityReviewCampaignID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("identity review campaign does not exist"),
			"not_exist",
			"identity_review_campaign_id",
		)
	}
	if _, err := irus.organizationUnitSvc.GetByID(formDTO.DepartmentID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("department does not exist"),
			"not_exist",
			"department_id",
		)
	}

	if formDTO.PrivilegeID != nil {
		if _, err := irus.privilegedAccountSvc.GetByID(*formDTO.PrivilegeID); err != nil {
			return utility.NewValidationError(
				fmt.Errorf("privilege does not exist"),
				"not_exist",
				"privilege_id",
			)
		}
	}

	if formDTO.AccountTypeID != nil {
		if _, err := irus.accountTypeSvc.GetByID(*formDTO.AccountTypeID); err != nil {
			return utility.NewValidationError(
				fmt.Errorf("account type does not exist"),
				"not_exist",
				"account_type_id",
			)
		}
	}

	return nil
}

func (irus *IdentityReviewUserService) UpdatePrivilege(id uint, currentUserID uint, formDTO *dtos.IdentityReviewUserPrivilegeUpdateFormDTO, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewUserDetailDTO, error) {
	existingUser, err := irus.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	oldUser := *existingUser

	if formDTO.PrivilegeID != nil {
		if _, err := irus.privilegedAccountSvc.GetByID(*formDTO.PrivilegeID); err != nil {
			return nil, utility.NewValidationError(
				fmt.Errorf("privilege does not exist"),
				"not_exist",
				"privilege_id",
			)
		}
	}

	now := time.Now()

	existingUser.PrivilegeID = formDTO.PrivilegeID
	existingUser.UpdatedByID = &currentUserID
	existingUser.UpdatedAt = &now

	if err := irus.repo.UpdatePrivilege(existingUser, tx...); err != nil {
		return nil, err
	}

	updatedUser, err := irus.GetFullDetailsByID(id, tx...)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := irus.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewUser{}.TableName(),
			id,
			map[string]interface{}{
				"old": &oldUser,
				"new": updatedUser,
			},
			currentUserID,
			"identity_access_review",
			"identity_review_user",
			"privilege",
			nil,
		); err != nil {
			utility.LogError(
				"Failed to log activity",
				map[string]interface{}{
					"error": err.Error(),
				},
			)
		}
	}()

	return updatedUser, nil
}

func (irus *IdentityReviewUserService) GenerateIdentityReviewUserTemplate(currentUserID uint) (*os.File, error) {
	templateHeaders := []string{
		"Name",
		"Identifier",
		"LastLogin",
		"AccountType",
		"Privilege",
		"Department",
		"CurrentRole",
	}

	templateData := [][]interface{}{}

	accountTypes, err := irus.accountTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account types: %v", err)
	}
	accountTypeData := [][]interface{}{{"AccountType"}}
	for _, accountType := range accountTypes {
		accountTypeData = append(accountTypeData, []interface{}{accountType.TitleEn})
	}

	privileges, err := irus.privilegedAccountSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch privileges: %v", err)
	}
	privilegeData := [][]interface{}{{"Privilege"}}
	for _, privilege := range privileges {
		privilegeData = append(privilegeData, []interface{}{privilege.TitleEn})
	}

	departments, err := irus.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %v", err)
	}
	departmentData := [][]interface{}{{"Department"}}
	for _, department := range departments {
		departmentData = append(departmentData, []interface{}{department.TitleEn})
	}

	dropDownConfig := map[string]bool{
		"AccountType": true,
		"Privilege":   true,
		"Department":  true,
	}

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":    append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"AccountType": accountTypeData,
		"Privilege":   privilegeData,
		"Department":  departmentData,
	}, dropDownConfig, map[string]string{
		"LastLogin": "yyyy-mm-dd hh:mm",
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (irus *IdentityReviewUserService) ValidateIdentityReviewUserTemplate(file io.Reader, identityReviewCampaignID uint, currentUserID uint) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"Name",
		"Identifier",
		"LastLogin",
		"AccountType",
		"Privilege",
		"Department",
		"CurrentRole",
	}

	snakeHeaders := []string{
		"name",
		"identifier",
		"last_login",
		"account_type",
		"privilege",
		"department",
		"current_role",
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}
	if !utility.ValidateHeaders(headers, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Identity Review User Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, header := range snakeHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = ""
			}
		}
		input = append(input, record)
	}

	validatedData, err := irus.validateIdentityReviewUsers(input, identityReviewCampaignID, currentUserID)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (irus *IdentityReviewUserService) ValidateImportIdentityReviewUsers(input []map[string]interface{}, identityReviewCampaignID uint, currentUserID uint) ([]map[string]interface{}, error) {
	validatedData, err := irus.validateIdentityReviewUsers(input, identityReviewCampaignID, currentUserID)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (irus *IdentityReviewUserService) ImportIdentityReviewUsers(input []map[string]interface{}, identityReviewCampaignID uint, currentUserID uint) ([]uint, error) {
	validatedData, err := irus.validateIdentityReviewUsers(input, identityReviewCampaignID, currentUserID)
	if err != nil {
		return nil, err
	}

	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validatedData {
		accountTypeID := utility.GetOptionalUintValue(record["account_type_id"])
		privilegeID := utility.GetOptionalUintValue(record["privilege_id"])

		formDTO := &dtos.IdentityReviewUserFormDTO{
			IdentityReviewCampaignID: identityReviewCampaignID,
			Username:                 record["name"].(string),
			UserIdentifier:           strings.TrimSpace(record["identifier"].(string)),
			LastLogin:                strings.TrimSpace(record["last_login"].(string)),
			CurrentRole:              strings.TrimSpace(record["current_role"].(string)),
			DepartmentID:             *utility.GetOptionalUintValue(record["department_id"]),
			PrivilegeID:              privilegeID,
			AccountTypeID:            accountTypeID,
		}

		identityReviewUser, err := irus.Create(currentUserID, formDTO)
		if err != nil {
			return nil, fmt.Errorf("failed to create identity review user for record: %+v, error: %v", record, err)
		}
		createdIDs = append(createdIDs, identityReviewUser.ID)
	}

	return createdIDs, nil
}

func (irus *IdentityReviewUserService) validateIdentityReviewUsers(input []map[string]interface{}, identityReviewCampaignID uint, currentUserID uint) ([]map[string]interface{}, error) {
	if _, err := irus.identityReviewCampaignSvc.GetByID(identityReviewCampaignID); err != nil {
		return nil, utility.NewValidationError(
			fmt.Errorf("identity review campaign does not exist"),
			"not_exist",
			"identity_review_campaign_id",
		)
	}

	headers := []string{
		"name",
		"identifier",
		"last_login",
		"account_type",
		"privilege",
		"department",
		"current_role",
	}

	accountTypes, err := irus.accountTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching account types: %v", err)
	}
	accountTypeTitles := make(map[string]uint)
	for _, accountType := range accountTypes {
		accountTypeTitles[strings.TrimSpace(accountType.TitleEn)] = accountType.ID
	}

	privileges, err := irus.privilegedAccountSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching privileges: %v", err)
	}
	privilegeTitles := make(map[string]uint)
	for _, privilege := range privileges {
		privilegeTitles[strings.TrimSpace(privilege.TitleEn)] = privilege.ID
	}

	departments, err := irus.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching departments: %v", err)
	}
	departmentTitles := make(map[string]uint)
	for _, department := range departments {
		departmentTitles[strings.TrimSpace(department.TitleEn)] = department.ID
	}

	var validatedData []map[string]interface{}
	for _, record := range input {
		rowErrors := make(map[string]interface{})
		identityReviewUserRecord := make(map[string]interface{})

		for key, value := range record {
			identityReviewUserRecord[key] = value
		}

		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {
			case "name":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Name is required"
				} else {
					identityReviewUserRecord[header] = strings.TrimSpace(value)
				}

			case "identifier":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "User Identifier is required"
				} else {
					identityReviewUserRecord[header] = strings.TrimSpace(value)
				}

			case "last_login":
				trimmedValue := strings.TrimSpace(value)
				if !ok || trimmedValue == "" {
					identityReviewUserRecord[header] = ""
				} else if _, err := time.Parse("2006-01-02 15:04", trimmedValue); err != nil {
					rowErrors[header] = "Last Login must be a valid date and time in YYYY-MM-DD HH:mm format"
				} else {
					identityReviewUserRecord[header] = trimmedValue
				}

			case "account_type":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Account type is required"
				} else if id, exists := accountTypeTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid account type: %s", value)
				} else {
					identityReviewUserRecord[header] = strings.TrimSpace(value)
					identityReviewUserRecord["account_type_id"] = id
				}

			case "privilege":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Privilege is required"
				} else if id, exists := privilegeTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid privilege: %s", value)
				} else {
					identityReviewUserRecord[header] = strings.TrimSpace(value)
					identityReviewUserRecord["privilege_id"] = id
				}

			case "department":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Department is required"
				} else if id, exists := departmentTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid department: %s", value)
				} else {
					identityReviewUserRecord["department_id"] = id
				}

			case "current_role":
				if !ok {
					identityReviewUserRecord[header] = ""
				} else {
					identityReviewUserRecord[header] = strings.TrimSpace(value)
				}
			}
		}

		identityReviewUserRecord["errors"] = rowErrors
		validatedData = append(validatedData, identityReviewUserRecord)
	}

	return validatedData, nil
}
