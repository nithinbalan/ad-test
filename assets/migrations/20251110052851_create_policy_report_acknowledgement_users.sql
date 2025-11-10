-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_acknowledgement_users (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    user_name_en NVARCHAR(255) NULL,
    user_name_ar NVARCHAR(255) NULL,
    organization_unit_title_en NVARCHAR(255) NULL,
    organization_unit_title_ar NVARCHAR(255) NULL,
    acknowledged BIT NOT NULL DEFAULT 0,
    acknowledged_on VARCHAR(50) NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_acknowledgement_users_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_acknowledgement_users_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_acknowledgement_users;
-- +goose StatementEnd
