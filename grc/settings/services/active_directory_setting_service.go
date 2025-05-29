package services

import (
	"crypto/tls"
	"fmt"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/go-ldap/ldap/v3"
)

type ActiveDirectorySettingServiceInterface interface {
	GetFirst() (*models.ActiveDirectorySetting, error)
	Create(currentUserID uint, dto dtos.ActiveDirectorySettingDTO) error
	CheckADConnection(dto *dtos.ActiveDirectorySettingDTO) error
	GetADConnection() error
}

type ActiveDirectorySettingService struct {
	repo           repositories.ActiveDirectorySettingRepositoryInterface
	integrationSvc IntegrationServiceInterface
}

func NewActiveDirectorySettingService(repo repositories.ActiveDirectorySettingRepositoryInterface, integrationSvc IntegrationServiceInterface) ActiveDirectorySettingServiceInterface {
	return &ActiveDirectorySettingService{
		repo:           repo,
		integrationSvc: integrationSvc,
	}
}

func (adss *ActiveDirectorySettingService) GetFirst() (*models.ActiveDirectorySetting, error) {
	return adss.repo.GetFirst()
}

func (adss *ActiveDirectorySettingService) Create(currentUserID uint, dto dtos.ActiveDirectorySettingDTO) error {
	// Step 1: Check connection first
	if err := adss.CheckADConnection(&dto); err != nil {
		return err
	}

	// Step 2: Encrypt passwords
	encryptedBindPassword, err := utility.EncryptData([]byte(dto.BindPassword))
	if err != nil {
		return err
	}
	encryptedUserPassword, err := utility.EncryptData([]byte(dto.Password))
	if err != nil {
		return err
	}

	// Step 3: Check integration setting
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return err
	}
	if !integration.IsEnable {
		return err
	}

	// Step 4: Save or Update setting
	activeDirectorySetting, err := adss.repo.GetFirst()
	if err != nil {
		return err
	}

	if activeDirectorySetting != nil {
		// Update
		activeDirectorySetting.IntegrationID = integration.ID
		activeDirectorySetting.Server = dto.Server
		activeDirectorySetting.Port = dto.Port
		activeDirectorySetting.UseSSL = dto.UseSSL
		activeDirectorySetting.BaseDN = dto.BaseDN
		activeDirectorySetting.BindDN = dto.BindDN
		activeDirectorySetting.BindPassword = encryptedBindPassword
		activeDirectorySetting.UserSearchAttr = dto.UserSearchAttr
		activeDirectorySetting.Username = dto.Username
		activeDirectorySetting.Password = encryptedUserPassword
		activeDirectorySetting.UpdatedBYID = &currentUserID

		if err := adss.repo.Update(activeDirectorySetting); err != nil {
			return err
		}
	} else {
		// Create new
		newSetting := models.ActiveDirectorySetting{
			IntegrationID:  integration.ID,
			Server:         dto.Server,
			Port:           dto.Port,
			UseSSL:         dto.UseSSL,
			BaseDN:         dto.BaseDN,
			BindDN:         dto.BindDN,
			BindPassword:   encryptedBindPassword,
			UserSearchAttr: dto.UserSearchAttr,
			Username:       dto.Username,
			Password:       encryptedUserPassword,
			CreatedBYID:    currentUserID,
		}
		if err := adss.repo.Create(&newSetting); err != nil {
			return err
		}
	}

	return nil
}

func (adss *ActiveDirectorySettingService) CheckADConnection(dto *dtos.ActiveDirectorySettingDTO) error {
	address := fmt.Sprintf("%s:%d", dto.Server, dto.Port)

	var l *ldap.Conn
	var err error
	if dto.UseSSL {
		l, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		l, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to connect to LDAP server: %w", err),
		)
	}
	defer l.Close()

	// Bind with BindDN and BindPassword
	if err := l.Bind(dto.BindDN, dto.BindPassword); err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("LDAP bind failed: %w", err),
		)
	}
	// Search user to validate
	const userFilter = "(&(objectCategory=person)(objectClass=user)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

	searchReq := ldap.NewSearchRequest(
		dto.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"givenName", "sn", "mail", "distinguishedName", "title"},
		nil,
	)
	sr, err := l.Search(searchReq)
	if err != nil || len(sr.Entries) == 0 {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("LDAP search failed or no users found: %w", err),
		)
	}

	return nil
}

func (adss *ActiveDirectorySettingService) GetADConnection() error {
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to find integration: %w", err),
		)
	}
	if !integration.IsEnable {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("integration 'active-directory' is disabled"),
		)
	}

	activeDirectorySetting, err := adss.repo.GetFirst()
	if err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to fetch saved AD settings: %w", err),
		)
	}
	if activeDirectorySetting == nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("no AD settings found"),
		)
	}

	decryptedBindPassword, err := utility.DecryptData(activeDirectorySetting.BindPassword)
	if err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to decrypt bind password: %w", err),
		)
	}
	decryptedUserPassword, err := utility.DecryptData(activeDirectorySetting.Password)
	if err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to decrypt user password: %w", err),
		)
	}

	// Prepare DTO manually
	dto := dtos.ActiveDirectorySettingDTO{
		Server:       activeDirectorySetting.Server,
		Port:         activeDirectorySetting.Port,
		UseSSL:       activeDirectorySetting.UseSSL,
		BaseDN:       activeDirectorySetting.BaseDN,
		BindDN:       activeDirectorySetting.BindDN,
		BindPassword: string(decryptedBindPassword),
		Password:     string(decryptedUserPassword),
	}

	// This will already return ActiveDirectoryConnectionError if fails
	return adss.CheckADConnection(&dto)
}
