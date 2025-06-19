-- +goose Up
-- +goose StatementBegin
CREATE TABLE report_template_settings 
(
    id INT PRIMARY KEY IDENTITY(1,1),
    logo_image_path VARCHAR(500) DEFAULT NULL,
    cover_page_bg_image_path VARCHAR(500) DEFAULT NULL, 
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT,
    updated_at DATETIME,

    CONSTRAINT FK_report_template_settings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_report_template_settings_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS report_template_settings;
-- +goose StatementEnd
