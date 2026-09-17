-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_report_sections (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    report_id INT NOT NULL,
    [type] VARCHAR(100) NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_psg_report_sections_reports
        FOREIGN KEY (report_id) REFERENCES project_security_gate_reports (id)
        ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_report_sections;
-- +goose StatementEnd
