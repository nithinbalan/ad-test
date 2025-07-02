-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'executive-summary', 'compliance-status-by-department', 'compliance-status-by-standard',
        'action-plan-implementation', 'conclusion-recommendations'
    )),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL,
    CONSTRAINT FK_compliance_report_sections_compliance_reports FOREIGN KEY (compliance_report_id) REFERENCES compliance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_report_sections;
-- +goose StatementEnd