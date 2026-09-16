-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_id INT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('executive-summary', 'scope', 'objectives', 'overview', 'details', 'findings', 'retest-results', 'management-sharing', 'closure')),

    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,

    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,

    order_num INT NOT NULL,

    CONSTRAINT FK_penetration_test_report_sections_penetration_test_reports FOREIGN KEY (penetration_test_report_id) REFERENCES penetration_test_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_sections;
-- +goose StatementEnd
