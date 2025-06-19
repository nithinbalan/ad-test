package utility

import (
	"bytes"
	"fmt"
	"html/template"
	"jethur/assets"
	"jethur/config"
	"jethur/grc/common/dtos"

	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	TLS      bool
}

func SendTemplatedEmail(to []string, subject, templateRelativePath string, templateData interface{}, emailTemplateSettingsData *dtos.EmailTemplateAndSocialDTO, attachments []string, cc ...[]string) error {
	emailConfig := getEmailConfig()

	templates, err := parseEmailTemplates(templateRelativePath, emailTemplateSettingsData)
	if err != nil {
		fmt.Printf("Failed to parse templates from path '%s': %v\n", templateRelativePath, err)
		return err
	} else {
		fmt.Println("Templates parsed successfully.")
	}

	var body bytes.Buffer
	if err := templates.Execute(&body, templateData); err != nil {
		fmt.Printf("Failed to execute template with data %v: %v\n", templateData, err)
		return err
	} else {
		fmt.Println("Template executed successfully, email body generated.")
	}

	message := createEmailMessage(to, subject, body.String(), attachments, emailConfig.From)

	// Check if CC is provided and not empty, then set CC in the header
	if len(cc) > 0 && len(cc[0]) > 0 {
		message.SetHeader("Cc", cc[0]...)
	}

	d := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.Username, emailConfig.Password)
	if !emailConfig.TLS {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		d.SSL = false
	}
	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func getEmailConfig() EmailConfig {
	return EmailConfig{
		Host:     config.GetEnv("EMAIL_HOST", ""),
		Port:     config.GetIntEnv("EMAIL_PORT", 25),
		Username: config.GetEnv("EMAIL_USERNAME", ""),
		Password: config.GetEnv("EMAIL_PASSWORD", ""),
		From:     config.GetEnv("EMAIL_FROM", "no-reply@jethur.com"),
		TLS:      config.GetBoolEnv("EMAIL_TLS", true),
	}
}

func parseEmailTemplates(templateRelativePath string, emailTemplateSettingsData *dtos.EmailTemplateAndSocialDTO) (*template.Template, error) {
	baseTemplatePath := "common/templates/email/base.html"
	headerTemplatePath := "common/templates/email/header.html"
	footerTemplatePath := "common/templates/email/footer.html"

	baseTemplate, err := assets.FS.ReadFile(baseTemplatePath)
	if err != nil {
		return nil, err
	}
	headerTemplate, err := assets.FS.ReadFile(headerTemplatePath)
	if err != nil {
		return nil, err
	}
	footerTemplate, err := assets.FS.ReadFile(footerTemplatePath)
	if err != nil {
		return nil, err
	}
	templateFile, err := assets.FS.ReadFile(templateRelativePath)
	if err != nil {
		return nil, err
	}

	templates, err := template.New("email").Funcs(template.FuncMap{
		"headerImageURL": func() string { return emailTemplateSettingsData.HeaderImageURL },
		"footerAddress":  func() string { return emailTemplateSettingsData.FooterAddress },
		"footerEmail":    func() string { return emailTemplateSettingsData.FooterEmail },
		"footerLogoURL":  func() string { return emailTemplateSettingsData.FooterLogoURL },
		"supportText":    func() string { return emailTemplateSettingsData.SupportText },
		"closingText":    func() string { return emailTemplateSettingsData.ClosingText },
		"facebook":       func() string { return *emailTemplateSettingsData.Facebook },
		"linkedin":       func() string { return *emailTemplateSettingsData.Linkedin },
		"twitter":        func() string { return *emailTemplateSettingsData.X },
	}).Parse(string(baseTemplate))
	if err != nil {
		return nil, err
	}
	_, err = templates.Parse(string(headerTemplate))
	if err != nil {
		return nil, err
	}
	_, err = templates.Parse(string(footerTemplate))
	if err != nil {
		return nil, err
	}
	_, err = templates.Parse(string(templateFile))
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func createEmailMessage(to []string, subject, body string, attachments []string, from string) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	for _, attachment := range attachments {
		if attachment != "" {
			message.Attach(attachment)
		}
	}

	return message
}
