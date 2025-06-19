-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_library_clauses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_library_version_id INT NOT NULL,
    clause_number NVARCHAR(100) NOT NULL,
    clause_title_en NVARCHAR(255) NOT NULL,
    clause_title_ar NVARCHAR(255) NULL,
    clause_description_en NVARCHAR(MAX) NULL,
    clause_description_ar NVARCHAR(MAX) NULL,
    compliance_library_clause_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_clc_library_versions FOREIGN KEY (compliance_library_version_id) REFERENCES compliance_library_versions(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_clc_self_reference FOREIGN KEY (compliance_library_clause_id) REFERENCES compliance_library_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_library_clauses;
-- +goose StatementEnd