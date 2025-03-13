-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_identifications (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    reference_number NVARCHAR(50) NOT NULL,
    risk_category_title_en NVARCHAR(255) NOT NULL,
    risk_category_title_ar NVARCHAR(255) NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    risk_source_title_en NVARCHAR(255) NULL,
    risk_source_title_ar NVARCHAR(255) NULL,
    risk_reported_at DATE NOT NULL,
    responsible_unit_title_en NVARCHAR(255) NOT NULL,
    responsible_unit_title_ar NVARCHAR(255) NOT NULL,

    inherent_risk_rating_title_en NVARCHAR(255) NULL,
    inherent_risk_rating_title_ar NVARCHAR(255) NULL,
    inherent_risk_rating_color_code NVARCHAR(7) NULL,

    inherent_likelihood_title_en NVARCHAR(255) NULL,
    inherent_likelihood_title_ar NVARCHAR(255) NULL,

    inherent_impact_title_en NVARCHAR(255) NULL,
    inherent_impact_title_ar NVARCHAR(255) NULL,

    residual_risk_rating_title_en NVARCHAR(255) NULL,
    residual_risk_rating_title_ar NVARCHAR(255) NULL,
    residual_risk_rating_color_code NVARCHAR(7) NULL,

    residual_likelihood_title_en NVARCHAR(255) NULL,
    residual_likelihood_title_ar NVARCHAR(255) NULL,

    residual_impact_title_en NVARCHAR(255) NULL,
    residual_impact_title_ar NVARCHAR(255) NULL,

    existing_controls NVARCHAR(MAX) NULL,
    proposed_controls NVARCHAR(MAX) NULL,

    risk_status_title_en NVARCHAR(255) NOT NULL,
    risk_status_title_ar NVARCHAR(255) NOT NULL,
    risk_status_color_code NVARCHAR(7) NOT NULL,

    CONSTRAINT FK_risk_report_identifications_risk_reports FOREIGN KEY (risk_report_id) REFERENCES risk_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_identifications;
-- +goose StatementEnd
