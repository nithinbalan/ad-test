-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_exception_reports (
    id INT IDENTITY(1,1) PRIMARY KEY,
    title NVARCHAR(255) NOT NULL,
    security_gate_ids NVARCHAR(MAX) NOT NULL,
    report_data NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT CK_project_security_gate_exception_reports_report_data_json
        CHECK (ISJSON(report_data) = 1),

    CONSTRAINT FK_project_security_gate_exception_reports_users_created
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_exception_reports;
-- +goose StatementEnd
