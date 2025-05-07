package services

import (
	"crypto/tls"
	"fmt"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"

	"github.com/go-ldap/ldap/v3"
)

type ActiveDirectorySettingServiceInterface interface {
	Authenticate(dto dtos.ADLoginRequestDTO) error
}

type ActiveDirectorySettingService struct {
	repo repositories.ActiveDirectorySettingRepositoryInterface
}

func NewActiveDirectorySettingService(repo repositories.ActiveDirectorySettingRepositoryInterface) ActiveDirectorySettingServiceInterface {
	return &ActiveDirectorySettingService{
		repo: repo,
	}
}

func (adss *ActiveDirectorySettingService) Authenticate(dto dtos.ADLoginRequestDTO) error {
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

	if err := l.Bind(dto.BindDN, dto.BindPassword); err != nil {
		return fmt.Errorf("bind user failed: %w", err)
	}

	// searchReq := ldap.NewSearchRequest(
	// 	dto.BaseDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", dto.UserSearchAttr, dto.Username),
	// 	[]string{"dn", "cn", "displayName", "mail", "memberOf"},
	// 	nil,
	// )
	// Filter for active user accounts
	const userFilter = "(&(objectCategory=person)(objectClass=user)(!(userAccountControl:1.2.840.113556.1.4.803:=2)))"

	searchReq := ldap.NewSearchRequest(
		dto.BaseDN, // 👈 pulled from request body
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"cn", "displayName", "mail"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil || len(sr.Entries) == 0 {
		return fmt.Errorf("user not found: %w", err)
	}

	entry := sr.Entries[0]
	userDN := entry.DN

	if err := l.Bind(userDN, dto.Password); err != nil {
		return fmt.Errorf("user auth failed: %w", err)
	}

	fmt.Printf("Total Users Found: %d\n", len(sr.Entries))
	for _, entry := range sr.Entries {
		fmt.Println("----")
		fmt.Println("CN:", entry.GetAttributeValue("cn"))
		fmt.Println("Display Name:", entry.GetAttributeValue("displayName"))
		fmt.Println("Email:", entry.GetAttributeValue("mail"))
	}

	return nil
}
