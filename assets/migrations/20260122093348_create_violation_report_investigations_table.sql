-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_investigations
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    violation_report_id INT NOT NULL,
    steps_taken NVARCHAR(MAX) NULL,
    investigation_date DATETIME NOT NULL,
    suggested_recommendation NVARCHAR(MAX) NULL,
    finding_summary NVARCHAR(MAX) NULL,
    disciplinary_procedure NVARCHAR(MAX) NULL,
    is_disciplinary_regulatory_compliance BIT NULL,
    actual_recommendation NVARCHAR(MAX) NULL,
    note_and_observations NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_violation_report_investigations_violation_reports FOREIGN KEY (violation_report_id) REFERENCES violation_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_investigations_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_investigations_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_investigations;

-- +goose StatementEnd