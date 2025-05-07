package dtos

type ADLoginRequestDTO struct {
	// Connection Details
	Server         string `json:"server" binding:"required"`           // e.g., ad1.example.local
	Port           int    `json:"port" binding:"required"`             // e.g., 389
	UseSSL         bool   `json:"use_ssl"`                             // true = LDAPS
	BaseDN         string `json:"base_dn" binding:"required"`          // e.g., DC=example,DC=local
	BindDN         string `json:"bind_dn" binding:"required"`          // e.g., CN=bind-user,...
	BindPassword   string `json:"bind_password" binding:"required"`    // service account password
	UserSearchAttr string `json:"user_search_attr" binding:"required"` // e.g., sAMAccountName

	// User Credentials
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
