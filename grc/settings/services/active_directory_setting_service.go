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
	if err := adss.TestADConnection(dto); err != nil {
		return fmt.Errorf("AD connection test failed: %w", err)
	}

	// Step 2: Encrypt passwords
	encryptedBindPassword, err := utility.EncryptData([]byte(dto.BindPassword))
	if err != nil {
		return fmt.Errorf("failed to encrypt bind password: %w", err)
	}
	encryptedUserPassword, err := utility.EncryptData([]byte(dto.Password))
	if err != nil {
		return fmt.Errorf("failed to encrypt user password: %w", err)
	}

	// Step 3: Check integration setting
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return fmt.Errorf("failed to find integration: %w", err)
	}
	if !integration.IsEnable {
		return fmt.Errorf("integration 'active-directory' is disabled")
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
			return fmt.Errorf("failed to update AD config: %w", err)
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

func (adss *ActiveDirectorySettingService) TestADConnection(dto dtos.ActiveDirectorySettingDTO) error {
	// 🔐 Decrypting here not needed, since dto contains raw passwords

	address := fmt.Sprintf("%s:%d", dto.Server, dto.Port)

	var l *ldap.Conn
	var err error
	if dto.UseSSL {
		l, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		l, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer l.Close()

	// Bind with BindDN and BindPassword
	if err := l.Bind(dto.BindDN, dto.BindPassword); err != nil {
		return fmt.Errorf("bind user failed: %w", err)
	}

	// 🔍 Search user to validate
	const userFilter = "(&(objectCategory=person)(objectClass=user)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

	searchReq := ldap.NewSearchRequest(
		dto.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"cn", "displayName", "mail"},
		nil,
	)
	sr, err := l.Search(searchReq)
	if err != nil || len(sr.Entries) == 0 {
		return fmt.Errorf("user not found: %w", err)
	}

	userDN := sr.Entries[0].DN

	// Bind with actual username and password (user login check)
	if err := l.Bind(userDN, dto.Password); err != nil {
		return fmt.Errorf("user auth failed: %w", err)
	}

	return nil
}

func (adss *ActiveDirectorySettingService) GetADConnection() error {
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return fmt.Errorf("failed to find integration: %w", err)
	}
	if !integration.IsEnable {
		return fmt.Errorf("integration 'active-directory' is disabled")
	}
	// 1. Fetch saved settings
	activeDirectorySetting, err := adss.repo.GetFirst()
	if err != nil {
		return fmt.Errorf("failed to fetch saved AD settings: %w", err)
	}
	if activeDirectorySetting == nil {
		return fmt.Errorf("no AD settings found")
	}

	// 2. Decrypt saved bind password
	decryptedBindPassword, err := utility.DecryptData(activeDirectorySetting.BindPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt bind password: %w", err)
	}
	bindPassword := string(decryptedBindPassword)

	// 3. Connect to LDAP server
	address := fmt.Sprintf("%s:%d", activeDirectorySetting.Server, activeDirectorySetting.Port)

	var l *ldap.Conn
	if activeDirectorySetting.UseSSL {
		l, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		l, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer l.Close()

	// 4. Bind with saved BindDN and decrypted password
	if err := l.Bind(activeDirectorySetting.BindDN, bindPassword); err != nil {
		return fmt.Errorf("bind failed: %w", err)
	}
	const userFilter = "(&(objectCategory=person)(objectClass=user)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

	// 5. (Optional) Run a very basic search to ensure base DN is correct
	searchReq := ldap.NewSearchRequest(
		activeDirectorySetting.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"cn", "displayName", "mail"},
		nil,
	)
	_, err = l.Search(searchReq)
	if err != nil {
		return fmt.Errorf("basic search failed: %w", err)
	}

	return nil
}
