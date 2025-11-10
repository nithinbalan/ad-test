-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_acknowledgement_organization_units (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    organization_unit_title_en NVARCHAR(255) NULL,
    organization_unit_title_ar NVARCHAR(255) NULL,
    total_users INT NOT NULL DEFAULT 0,
    acknowledged_users INT NOT NULL DEFAULT 0,
    not_acknowledged_users INT NOT NULL DEFAULT 0,
    acknowledged_percent NUMERIC(5,2) NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_acknowledgement_organization_units_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_acknowledgement_organization_units_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_acknowledgement_organization_units;
-- +goose StatementEnd
