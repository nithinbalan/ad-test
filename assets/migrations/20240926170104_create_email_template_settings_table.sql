-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_template_settings (
  id INT PRIMARY KEY IDENTITY(1,1),
  header_image_url NVARCHAR(MAX) NOT NULL,
  support_text NVARCHAR(MAX) NOT NULL,
  closing_text NVARCHAR(MAX) NOT NULL,
  footer_logo_url NVARCHAR(MAX) NOT NULL,
  footer_address NVARCHAR(MAX) NOT NULL,
  footer_email NVARCHAR(MAX) NOT NULL,
  created_by_id INT NOT NULL,
  created_at DATE NOT NULL,
  updated_by_id INT,
  updated_at DATE
  CONSTRAINT FK_email_template_settings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_email_template_settings_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.email_template_settings;
-- +goose StatementEnd
