-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    type VARCHAR(100) NOT NULL CHECK (type IN (
        'main-executive-summary',
        'governance-summary', 
        'executive-summary-governance', 
        'policy-approved-overview', 
        'governance-tasks', 
        'conclusion-governance',
        'main-conclusion'
    )),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL,
    CONSTRAINT FK_grc_report_sections_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_sections;
-- +goose StatementEnd
