-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_identifications (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    reference_number VARCHAR(50) NOT NULL,
    risk_category_title_en VARCHAR(255) NOT NULL,
    risk_category_title_ar VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    risk_source_title_en VARCHAR(255),
    risk_source_title_ar VARCHAR(255),
    risk_reported_at DATE NOT NULL,
    responsible_unit_title_en VARCHAR(255),
    responsible_unit_title_ar VARCHAR(255),

    inherent_risk_rating_title_en VARCHAR(255),
    inherent_risk_rating_title_ar VARCHAR(255),
    inherent_risk_rating_color_code VARCHAR(7),

    inherent_likelihood_title_en VARCHAR(255),
    inherent_likelihood_title_ar VARCHAR(255),

    inherent_impact_title_en VARCHAR(255),
    inherent_impact_title_ar VARCHAR(255),

    residual_risk_rating_title_en VARCHAR(255),
    residual_risk_rating_title_ar VARCHAR(255),
    residual_risk_rating_color_code VARCHAR(7),

    residual_likelihood_title_en VARCHAR(255),
    residual_likelihood_title_ar VARCHAR(255),

    residual_impact_title_en VARCHAR(255),
    residual_impact_title_ar VARCHAR(255),

    existing_controls TEXT,
    proposed_controls TEXT,

    risk_status_title_en VARCHAR(255),
    risk_status_title_ar VARCHAR(255),
    risk_status_color_code VARCHAR(7),

    CONSTRAINT FK_risk_report_identifications_risk_reports FOREIGN KEY (risk_report_id) REFERENCES risk_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_identifications;
-- +goose StatementEnd
