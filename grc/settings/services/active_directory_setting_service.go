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
	// 🔐 Encrypt bind password before saving
	dataToEncrypt := []byte(dto.BindPassword)
	encryptedBindPassword, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return err
	}
	encryptedUserPassword, err := utility.EncryptData([]byte(dto.Password))
	if err != nil {
		return fmt.Errorf("failed to encrypt user password: %w", err)
	}
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return fmt.Errorf("failed to find integration: %w", err)
	}
	// ✅ Check if integration is enabled
	if !integration.IsEnable {
		return fmt.Errorf("integration 'active-directory' is disabled")
	}
	// Check if an entry already exists
	activeDirectorySetting, err := adss.repo.GetFirst()
	if err != nil {
		return err
	}

	if activeDirectorySetting != nil {
		// Update activeDirectorySetting
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
		activeDirectorySetting := models.ActiveDirectorySetting{
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
		if err := adss.repo.Create(&activeDirectorySetting); err != nil {
			return err
		}
	}

	// 2. 🔓 Decrypt password for use
	decryptedBindPassword, err := utility.DecryptData(activeDirectorySetting.BindPassword)
	if err != nil {
		return fmt.Errorf("failed to decrypt bind password: %w", err)
	}
	bindPassword := string(decryptedBindPassword)
	// 3. ✅ Connect to LDAP
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

	// 4. ✅ Bind with saved config
	if err := l.Bind(activeDirectorySetting.BindDN, bindPassword); err != nil {
		return fmt.Errorf("bind user failed: %w", err)
	}
	// searchReq := ldap.NewSearchRequest(
	// 	activeDirectorySetting.BaseDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", activeDirectorySetting.UserSearchAttr, dto.Username),
	// 	[]string{"dn", "cn", "displayName", "mail"},
	// 	nil,
	// )

	// 5. 🔍 Search for user from DTO
	const userFilter = "(&(objectCategory=person)(objectClass=user)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

	searchReq := ldap.NewSearchRequest(
		activeDirectorySetting.BaseDN, // 👈 pulled from request body
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

	// 6. ✅ Bind with user credentials (login check)
	if err := l.Bind(userDN, dto.Password); err != nil {
		return fmt.Errorf("user auth failed: %w", err)
	}

	// 7. 🖨️ Print success and user info
	fmt.Printf("Total Users Found: %d\n", len(sr.Entries))
	for _, entry := range sr.Entries {
		fmt.Println("----")
		fmt.Println("CN:", entry.GetAttributeValue("cn"))
		fmt.Println("Display Name:", entry.GetAttributeValue("displayName"))
		fmt.Println("Email:", entry.GetAttributeValue("mail"))
	}

	return nil
}
