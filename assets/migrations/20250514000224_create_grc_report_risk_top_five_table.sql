-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_risk_top_five (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    risk_reference_number NVARCHAR(50) NOT NULL,
    risk_category_title_en NVARCHAR(255) NOT NULL,
    risk_category_title_ar NVARCHAR(255) NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    responsible_user_name_en NVARCHAR(255) NOT NULL,
    responsible_user_name_ar NVARCHAR(255) NOT NULL,
    responsible_unit_title_en NVARCHAR(255) NOT NULL,
    responsible_unit_title_ar NVARCHAR(255) NOT NULL,
    is_open BIT NOT NULL DEFAULT 1,

    inherent_risk_rating_title_en NVARCHAR(255) NULL,
    inherent_risk_rating_title_ar NVARCHAR(255) NULL,
    inherent_risk_rating_color_code NVARCHAR(7) NULL,
    inherent_risk_rating_score INT NULL,

    inherent_likelihood_title_en NVARCHAR(255) NULL,
    inherent_likelihood_title_ar NVARCHAR(255) NULL,
    inherent_impact_title_en NVARCHAR(255) NULL,
    inherent_impact_title_ar NVARCHAR(255) NULL,

    residual_risk_rating_title_en NVARCHAR(255) NULL,
    residual_risk_rating_title_ar NVARCHAR(255) NULL,
    residual_risk_rating_color_code NVARCHAR(7) NULL,
    residual_risk_rating_score INT NULL,

    residual_likelihood_title_en NVARCHAR(255) NULL,
    residual_likelihood_title_ar NVARCHAR(255) NULL,
    residual_impact_title_en NVARCHAR(255) NULL,
    residual_impact_title_ar NVARCHAR(255) NULL,

    existing_controls NVARCHAR(MAX) NULL,

    risk_response_strategy_title_en NVARCHAR(255) NULL,
    risk_response_strategy_title_ar NVARCHAR(255) NULL,

    approval_note NVARCHAR(MAX) NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_grc_report_risk_top_five_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_risk_top_five_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_risk_top_five;
-- +goose StatementEnd
