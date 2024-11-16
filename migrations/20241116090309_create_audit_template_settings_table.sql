-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_template_settings 
(
    id INT PRIMARY KEY IDENTITY(1,1),
    logo_image_path VARCHAR(500) DEFAULT NULL,
    commencement_letter_header_bg_image_path VARCHAR(500) DEFAULT NULL, 
    draft_report_cover_page_bg_image_path VARCHAR(500) DEFAULT NULL, 
    final_report_cover_page_bg_image_path VARCHAR(500) DEFAULT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT,
    updated_at DATETIME,

    CONSTRAINT FK_audit_template_settings_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_template_settings_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_template_settings;
-- +goose StatementEnd
