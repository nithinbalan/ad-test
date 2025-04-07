-- +goose Up
-- +goose StatementBegin
CREATE TABLE governance_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    governance_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'executive-summary', 'policy-acknowledgement', 'meetings', 'actions', 'conclusion'
    )),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL,
    CONSTRAINT FK_governance_report_sections_governance_reports FOREIGN KEY (governance_report_id) REFERENCES governance_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE governance_report_sections;
-- +goose StatementEnd
