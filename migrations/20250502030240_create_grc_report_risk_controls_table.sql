-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_risk_controls (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    control_code NVARCHAR(100) NOT NULL,
    control_category_title_en NVARCHAR(255) NOT NULL,
    control_category_title_ar NVARCHAR(255) NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,

    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    responsible_unit_title_en NVARCHAR(255) NOT NULL,
    responsible_unit_title_ar NVARCHAR(255) NOT NULL,

    control_implementation_status_title_en NVARCHAR(255) NOT NULL,
    control_implementation_status_title_ar NVARCHAR(255) NOT NULL,
    control_implementation_status_color_code NVARCHAR(50) NOT NULL,

    control_efficiency_status_title_en NVARCHAR(255) NOT NULL,
    control_efficiency_status_title_ar NVARCHAR(255) NOT NULL,
    control_efficiency_status_color_code NVARCHAR(50) NOT NULL,

    test_frequency_title_en NVARCHAR(255) NULL,
    test_frequency_title_ar NVARCHAR(255) NULL,
    next_test_date DATE NULL,
    test_responsible_user_name_en NVARCHAR(255) NULL,
    test_responsible_user_name_ar NVARCHAR(255) NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_grc_report_risk_controls_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_risk_controls_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_risk_controls;
-- +goose StatementEnd
