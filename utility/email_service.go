package utility

import (
	"bytes"
	"html/template"
	"jethur/assets"
	"jethur/config"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func SendTemplatedEmail(to []string, subject, templateRelativePath string, templateData interface{}, attachments []string) error {

	emailConfig := getEmailConfig()

	templates, err := parseEmailTemplates(templateRelativePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := templates.Execute(&body, templateData); err != nil {
		return err
	}

	message := createEmailMessage(to, subject, body.String(), attachments, emailConfig.From)

	d := gomail.NewDialer(emailConfig.Host, emailConfig.Port, emailConfig.Username, emailConfig.Password)
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
	}
}

func parseEmailTemplates(templateRelativePath string) (*template.Template, error) {
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

	templates, err := template.New("email").Parse(string(baseTemplate))
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
