-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_report_risk_summary (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_report_id INT NOT NULL,
    total_risk_count INT NOT NULL DEFAULT 0,
    open_risk_count INT NOT NULL DEFAULT 0,
    closed_risk_count INT NOT NULL DEFAULT 0,
    closure_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_organization_unit_report_risk_summary_reports 
        FOREIGN KEY (organization_unit_report_id) REFERENCES organization_unit_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_organization_unit_report_risk_summary_users_1 
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_report_risk_summary;
-- +goose StatementEnd
