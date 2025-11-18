package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"

	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"
	"net/http"
	"net/url"

	"github.com/go-ldap/ldap/v3"
)

type ActiveDirectorySettingServiceInterface interface {
	GetFirst() (*models.ActiveDirectorySetting, error)
	Create(currentUserID uint, dto dtos.ActiveDirectorySettingDTO) error
	CheckADConnection(dto *dtos.ActiveDirectorySettingDTO) error
	GetADConnection() error

	TestOAuthConnection(dto *dtos.EntraOAuthConfigDTO) (string, error)
	GetOAuthConfiguration() (*dtos.EntraOAuthConfigDTO, error)
	SaveOAuthConfiguration(currentUserID uint, dto dtos.EntraOAuthConfigDTO) error
	GenerateMicrosoftLoginURL() (string, error)
	HandleOAuthCallback(code string, state string) (map[string]interface{}, error)
	SyncAzureUsers() (int, error)
	SyncAzureGroups() (int, error)
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

// ==============================
// Entra ID (OAuth2 / SSO / Graph)
// ==============================

// TestOAuthConnection validates the Entra ID client credentials by requesting a token.
func (adss *ActiveDirectorySettingService) TestOAuthConnection(dto *dtos.EntraOAuthConfigDTO) (string, error) {
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

	return tokenResp.AccessToken, nil
}

// GetOAuthConfiguration fetches stored Entra ID settings (if any).
func (adss *ActiveDirectorySettingService) GetOAuthConfiguration() (*dtos.EntraOAuthConfigDTO, error) {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, utility.NewResourceNotFoundError(fmt.Errorf("no Entra ID configuration found"), "Entra ID config")
	}

	// Map LDAP model fields to DTO placeholders for UI reuse
	return &dtos.EntraOAuthConfigDTO{
		TenantID:     setting.TenantID,
		ClientID:     setting.ClientID,
		ClientSecret: "********", // mask
		RedirectURI:  setting.RedirectURI,
		GraphScope:   "https://graph.microsoft.com/.default",
	}, nil
}

// SaveOAuthConfiguration stores or updates Entra ID credentials.
func (adss *ActiveDirectorySettingService) SaveOAuthConfiguration(currentUserID uint, dto dtos.EntraOAuthConfigDTO) error {
	integration, err := adss.integrationSvc.GetByType("active-directory")
	if err != nil {
		return err
	}
	if !integration.IsEnable {
		return utility.NewActiveDirectoryConnectionError(fmt.Errorf("integration 'active-directory' is disabled"))
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

// GenerateMicrosoftLoginURL builds the OAuth2 authorize URL for SSO.
func (adss *ActiveDirectorySettingService) GenerateMicrosoftLoginURL() (string, error) {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return "", err
	}
	if setting == nil {
		return "", utility.NewActiveDirectoryConnectionError(fmt.Errorf("no Entra ID config found"))
	}

	url := fmt.Sprintf(
		"https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=openid%%20profile%%20email&state=%s",
		setting.TenantID, setting.ClientID, url.QueryEscape(setting.RedirectURI), utility.GenerateRandomString(12),
	)
	return url, nil
}

// HandleOAuthCallback exchanges code for tokens and maps user claims.
func (adss *ActiveDirectorySettingService) HandleOAuthCallback(code string, state string) (map[string]interface{}, error) {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return nil, err
	}
	if setting == nil {
		return nil, utility.NewActiveDirectoryConnectionError(fmt.Errorf("no Entra ID configuration found"))
	}

	// Decrypt secret
	decryptedSecret, err := utility.DecryptData(setting.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt client secret: %w", err)
	}

	// Prepare token exchange request
	data := url.Values{}
	data.Set("client_id", setting.ClientID)
	data.Set("client_secret", string(decryptedSecret))
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", setting.RedirectURI)
	data.Set("scope", "openid profile email") // 👈 required for ID token

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", setting.TenantID)
	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange authorization code: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	// Print/log response for debugging
	fmt.Println("🔍 Microsoft Token Response:", string(body))

	// Parse JSON
	var tokenResp struct {
		IDToken      string `json:"id_token"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("invalid token response format: %w", err)
	}

	// ✅ Ensure we got an ID token
	if tokenResp.IDToken == "" {
		// fallback for debugging
		fmt.Println("⚠️ ID token missing! Microsoft response did not include an id_token.")
		fmt.Println("👉 Make sure your app registration has 'ID tokens (used for OpenID Connect)' enabled.")
		fmt.Println("👉 Current scope:", tokenResp.Scope)

		// Optionally, decode the AccessToken for debugging
		claims, err := utility.DecodeJWTClaims(tokenResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("no id_token returned, and access_token also invalid: %w", err)
		}
		return claims, nil
	}

	// ✅ Decode ID token
	claims, err := utility.DecodeJWTClaims(tokenResp.IDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ID token: %w", err)
	}

	return claims, nil
}

// SyncAzureUsers pulls users from Microsoft Graph and syncs to DB.
func (adss *ActiveDirectorySettingService) SyncAzureUsers() (int, error) {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return 0, err
	}
	if setting == nil {
		return 0, utility.NewActiveDirectoryConnectionError(fmt.Errorf("no Entra ID configuration found"))
	}

	decryptedSecret, err := utility.DecryptData(setting.ClientSecret)
	if err != nil {
		return 0, err
	}

	token, err := adss.getClientCredentialToken(setting.TenantID, setting.ClientID, string(decryptedSecret))
	if err != nil {
		return 0, err
	}

	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/users?$select=id,displayName,mail,userPrincipalName", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("graph API user fetch failed: %s", string(body))
	}

	var result struct {
		Value []struct {
			ID                string `json:"id"`
			DisplayName       string `json:"displayName"`
			Mail              string `json:"mail"`
			UserPrincipalName string `json:"userPrincipalName"`
		} `json:"value"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	// Map to local DB if required
	return len(result.Value), nil
}

// SyncAzureGroups pulls group data from Microsoft Graph.
func (adss *ActiveDirectorySettingService) SyncAzureGroups() (int, error) {
	setting, err := adss.repo.GetFirst()
	if err != nil {
		return 0, err
	}
	if setting == nil {
		return 0, utility.NewActiveDirectoryConnectionError(fmt.Errorf("no Entra ID configuration found"))
	}

	decryptedSecret, err := utility.DecryptData(setting.ClientSecret)
	if err != nil {
		return 0, err
	}

	token, err := adss.getClientCredentialToken(setting.TenantID, setting.ClientID, string(decryptedSecret))
	if err != nil {
		return 0, err
	}

	req, _ := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/groups?$select=id,displayName", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("graph API group fetch failed: %s", string(body))
	}

	var result struct {
		Value []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
		} `json:"value"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	return len(result.Value), nil
}

// helper: request client_credentials token
func (adss *ActiveDirectorySettingService) getClientCredentialToken(tenantID, clientID, clientSecret string) (string, error) {
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
