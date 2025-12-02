package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"jethur/grc/models"

	organizationRepositories "jethur/grc/organization/repositories"

	// "jethur/grc/settings/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"net/http"
	"net/url"

	"github.com/go-ldap/ldap/v3"
)

type ActiveDirectorySettingServiceInterface interface {
	GetFirst() (*models.ActiveDirectorySetting, error)
	Create(currentUserID uint, dto settingsDTOs.ActiveDirectorySettingDTO) error
	CheckADConnection(dto *settingsDTOs.ActiveDirectorySettingDTO) error
	GetADConnection() error
	Delete(currentUserID uint) error

	TestOAuthConnection(dto *settingsDTOs.EntraOAuthConfigDTO) (string, error)
	GetOAuthConfiguration() (*models.ActiveDirectorySetting, error)
	SaveOAuthConfiguration(currentUserID uint, dto settingsDTOs.EntraOAuthConfigDTO) error
	DeleteOAuthConfiguration(currentUserID uint) error

	GetAzureUserGroups(userID string) ([]map[string]interface{}, error)

	GetClientCredentialToken(tenantID, clientID, clientSecret string) (string, error)
}

type ActiveDirectorySettingService struct {
	repo             repositories.ActiveDirectorySettingRepositoryInterface
	userRepo         organizationRepositories.UserRepositoryInterface
	organizationRepo repositories.OrganizationLevelRepositoryInterface
	integrationSvc   IntegrationServiceInterface
}

func NewActiveDirectorySettingService(repo repositories.ActiveDirectorySettingRepositoryInterface, userRepo organizationRepositories.UserRepositoryInterface, organizationRepo repositories.OrganizationLevelRepositoryInterface, integrationSvc IntegrationServiceInterface) ActiveDirectorySettingServiceInterface {
	return &ActiveDirectorySettingService{
		repo:             repo,
		userRepo:         userRepo,
		integrationSvc:   integrationSvc,
		organizationRepo: organizationRepo,
	}
}

func (adss *ActiveDirectorySettingService) GetFirst() (*models.ActiveDirectorySetting, error) {
	return adss.repo.GetFirst()
}

func (adss *ActiveDirectorySettingService) Create(currentUserID uint, dto settingsDTOs.ActiveDirectorySettingDTO) error {
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

	// emptyCredential := &settingsDTOs.IntegrationCredentialFormDTO{}
	// _, err = adss.integrationSvc.EnableIntegration(integration.ID, currentUserID, emptyCredential)
	// if err != nil {
	// 	return fmt.Errorf("failed to enable integration: %w", err)
	// }
	_, err = adss.integrationSvc.EnableADIntegration(integration, currentUserID)
	if err != nil {
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

func (adss *ActiveDirectorySettingService) CheckADConnection(dto *settingsDTOs.ActiveDirectorySettingDTO) error {
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
	// if err := l.Bind(dto.BindDN, dto.BindPassword); err != nil {
	// 	return utility.NewActiveDirectoryConnectionError(
	// 		fmt.Errorf("LDAP bind failed: %w", err),
	// 	)
	// }
	// Bind and perform a **light-weight health-check** instead of pulling all users.
	// We query the LDAP *Root DSE* (base DN = ""), which always returns a single
	// entry and never hits the AD “size-limit-exceeded” cap.

	if err := l.Bind(dto.BindDN, dto.BindPassword); err != nil {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("LDAP bind failed: %w", err),
		)
	}

	// Root DSE query – proves the bind, BaseDN visibility, and overall connectivity.
	searchReq := ldap.NewSearchRequest(
		"",                   // empty DN → Root DSE
		ldap.ScopeBaseObject, // just that one entry
		ldap.NeverDerefAliases,
		0, 0, false, // no size/time limits needed
		"(objectClass=*)",                // standard catch-all filter
		[]string{"defaultNamingContext"}, // return a harmless attribute
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil || len(sr.Entries) == 0 {
		return utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("LDAP connectivity check failed: %w", err),
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
	dto := settingsDTOs.ActiveDirectorySettingDTO{
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

func (adss *ActiveDirectorySettingService) Delete(currentUserID uint) error {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return err
	}

	if setting == nil {
		return fmt.Errorf("active directory setting not found")
	} else {
		adss.repo.Delete(setting)
	}

	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return err
	}

	_, err = adss.integrationSvc.DisableADIntegration(integration, currentUserID)
	if err != nil {
		return err
	}
	return nil
}

// ==============================
// Entra ID (OAuth2 / SSO / Graph)
// ==============================

// TestOAuthConnection validates the Entra ID client credentials by requesting a token.
func (adss *ActiveDirectorySettingService) TestOAuthConnection(dto *settingsDTOs.EntraOAuthConfigDTO) (string, error) {
	form := url.Values{}
	form.Set("client_id", dto.ClientID)
	form.Set("client_secret", dto.ClientSecret)
	form.Set("scope", "https://graph.microsoft.com/.default")
	form.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", dto.TenantID)
	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		return "", utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("failed to connect to Entra ID token endpoint: %w", err),
		)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("token request failed: %s", string(body)),
		)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", utility.NewActiveDirectoryConnectionError(
			fmt.Errorf("invalid token response: %w", err),
		)
	}

	if err := adss.validateRedirectURI(dto, tokenResp.AccessToken); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// GetOAuthConfiguration fetches stored Entra ID settings (if any).
func (adss *ActiveDirectorySettingService) GetOAuthConfiguration() (*models.ActiveDirectorySetting, error) {
	// setting, err := adss.repo.GetFirst()
	// if err != nil {
	// 	return nil, err
	// }

	// // Map LDAP model fields to DTO placeholders for UI reuse
	// return &settingsDTOs.EntraOAuthConfigDTO{
	// 	TenantID:     setting.TenantID,
	// 	ClientID:     setting.ClientID,
	// 	ClientSecret: "********", // mask
	// 	RedirectURI:  setting.RedirectURI,
	// 	GraphScope:   "https://graph.microsoft.com/.default",
	// }, nil
	res, err := adss.repo.GetFirst()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SaveOAuthConfiguration stores or updates Entra ID credentials.
func (adss *ActiveDirectorySettingService) SaveOAuthConfiguration(currentUserID uint, dto settingsDTOs.EntraOAuthConfigDTO) error {
	integration, err := adss.integrationSvc.GetByType("azure-ad")
	if err != nil {
		return err
	}

	_, err = adss.integrationSvc.EnableADIntegration(integration, currentUserID)
	if err != nil {
		return fmt.Errorf("failed to enable integration: %w", err)
	}

	encryptedSecret, err := utility.EncryptData([]byte(dto.ClientSecret))
	if err != nil {
		return err
	}

	setting, err := adss.repo.GetFirst()
	if err != nil {
		return err
	}

	if setting != nil {
		setting.TenantID = dto.TenantID
		setting.ClientID = dto.ClientID
		setting.ClientSecret = encryptedSecret
		setting.RedirectURI = dto.RedirectURI
		setting.UpdatedBYID = &currentUserID
		return adss.repo.Update(setting)
	}

	newSetting := models.ActiveDirectorySetting{
		IntegrationID: integration.ID,
		TenantID:      dto.TenantID,
		ClientID:      dto.ClientID,
		ClientSecret:  encryptedSecret,
		RedirectURI:   dto.RedirectURI,
		CreatedBYID:   currentUserID,
	}
	return adss.repo.Create(&newSetting)
}

func (adss *ActiveDirectorySettingService) DeleteOAuthConfiguration(currentUserID uint) error {

	setting, err := adss.repo.GetFirst()
	if err != nil {
		return err
	}
	if setting == nil {
		return fmt.Errorf("no OAuth/Active Directory setting found")
	}

	integration, err := adss.integrationSvc.GetByType("azure-ad")
	if err != nil {
		return err
	}

	_, err = adss.integrationSvc.DisableADIntegration(integration, currentUserID)
	if err != nil {
		return err
	}

	return adss.repo.Delete(setting)
}

// helper: request client_credentials token
func (adss *ActiveDirectorySettingService) GetClientCredentialToken(tenantID, clientID, clientSecret string) (string, error) {
	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("scope", "https://graph.microsoft.com/.default")
	form.Set("grant_type", "client_credentials")

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		return "", fmt.Errorf("failed token request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token request failed: %s", string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}
	return tokenResp.AccessToken, nil
}

func (adss *ActiveDirectorySettingService) GetAzureUserGroups(userID string) ([]map[string]interface{}, error) {

	// 1. Load Azure AD config
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, fmt.Errorf("Azure AD configuration not found")
	}

	// 2. Decrypt client secret
	decryptedSecret, err := utility.DecryptData(setting.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 3. Get access token
	token, err := adss.GetClientCredentialToken(setting.TenantID, setting.ClientID, string(decryptedSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 4. Build Graph URL
	url := fmt.Sprintf(
		"https://graph.microsoft.com/v1.0/users/%s/memberOf?$select=id,displayName,description,groupTypes",
		userID,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Graph API error: %s", string(body))
	}

	// 5. Parse response
	var result struct {
		Value []struct {
			ID          string   `json:"id"`
			Name        string   `json:"displayName"`
			Description string   `json:"description"`
			GroupTypes  []string `json:"groupTypes"`
		} `json:"value"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// 6. Format response
	final := []map[string]interface{}{}
	for _, g := range result.Value {

		// filter: only group objects
		if g.ID == "" || g.Name == "" {
			continue
		}

		final = append(final, map[string]interface{}{
			"id":          g.ID,
			"name":        g.Name,
			"description": g.Description,
			"group_types": g.GroupTypes,
		})
	}

	return final, nil
}

func (adss *ActiveDirectorySettingService) validateRedirectURI(dto *settingsDTOs.EntraOAuthConfigDTO, accessToken string) error {
	if _, err := url.ParseRequestURI(dto.RedirectURI); err != nil {
		return fmt.Errorf("invalid redirect URI format: %w", err)
	}

	graphURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/applications?$filter=appId%%20eq%%20'%s'", dto.ClientID)

	req, _ := http.NewRequest("GET", graphURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	graphResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to contact Microsoft Graph: %w", err)
	}
	defer graphResp.Body.Close()

	graphBody, _ := io.ReadAll(graphResp.Body)

	var appResp struct {
		Value []struct {
			Web struct {
				RedirectUris []string `json:"redirectUris"`
			} `json:"web"`
		} `json:"value"`
	}
	json.Unmarshal(graphBody, &appResp)

	for _, r := range appResp.Value[0].Web.RedirectUris {
		if r == dto.RedirectURI {
			return nil // success
		}
	}

	return fmt.Errorf("redirect URI not registered in Azure App: %s", dto.RedirectURI)
}
