-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_risks (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    risk_title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    risk_category_title_en NVARCHAR(255) NULL,
    risk_category_title_ar NVARCHAR(255) NULL,
    risk_status_title_en NVARCHAR(100) NULL,
    risk_status_title_ar NVARCHAR(100) NULL,
    risk_status_color_code VARCHAR(9) NULL,
    risk_likelihood_title_en NVARCHAR(100) NULL,
    risk_likelihood_title_ar NVARCHAR(100) NULL,
    risk_impact_title_en NVARCHAR(100) NULL,
    risk_impact_title_ar NVARCHAR(100) NULL,
    risk_rating_title_en NVARCHAR(100) NULL,
    risk_rating_title_ar NVARCHAR(100) NULL,
    risk_rating_color_code VARCHAR(9) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_risks_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_risks_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_risks;
-- +goose StatementEnd
