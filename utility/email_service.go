package utility

import (
	"bytes"
	"html/template"
	"jethur/config"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendTemplatedEmail(to []string, subject, templateRelativePath string, templateData interface{}, attachments []string) error {

	rootPath := config.GetEnv("ROOT_PATH", "")
	emailFrom := config.GetEnv("EMAIL_FROM", "no-reply@example.com")
	emailHost := config.GetEnv("EMAIL_HOST", "")
	emailPort, _ := strconv.Atoi(config.GetEnv("EMAIL_PORT", "25"))
	emailUser := config.GetEnv("EMAIL_USERNAME", "")
	emailPass := config.GetEnv("EMAIL_PASSWORD", "")

	baseTemplatePath := filepath.Join(rootPath, "grc", "common", "templates", "email", "base.html")
	headerTemplatePath := filepath.Join(rootPath, "grc", "common", "templates", "email", "header.html")
	footerTemplatePath := filepath.Join(rootPath, "grc", "common", "templates", "email", "footer.html")
	templateRelativePath = filepath.Join(rootPath, "grc", templateRelativePath)

	template, err := template.ParseFiles(baseTemplatePath, headerTemplatePath, footerTemplatePath, templateRelativePath)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err = template.Execute(&body, templateData); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	for _, attachment := range attachments {
		if attachment != "" {
			m.Attach(attachment)
		}
	}

	d := gomail.NewDialer(emailHost, emailPort, emailUser, emailPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
