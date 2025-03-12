-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN (
        'executive-summary', 'context', 'identification', 'assessment', 
        'mitigation', 'kri', 'communication', 'recommendation', 'documents'
    )),
    title_en VARCHAR(255) NOT NULL,
    title_ar VARCHAR(255) NOT NULL,
    description_en TEXT,
    description_ar TEXT,
    order_num INT NOT NULL,
    CONSTRAINT FK_risk_report_sections_risk_reports FOREIGN KEY (risk_report_id) REFERENCES risk_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_report_sections;
-- +goose StatementEnd
