-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_risk_kris (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    risk_register_title NVARCHAR(255) NOT NULL,
    kri_register_title NVARCHAR(255) NOT NULL,
    kri_register_description NVARCHAR(MAX) NULL,
    low_threshold INT NULL,
    medium_threshold INT NULL,
    current_value INT NULL,
    risk_status_title_en NVARCHAR(255) NOT NULL,
    risk_status_title_ar NVARCHAR(255) NOT NULL,
    risk_status_title_color_code NVARCHAR(50) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_grc_report_risk_kris_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_risk_kris_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_risk_kris;
-- +goose StatementEnd
