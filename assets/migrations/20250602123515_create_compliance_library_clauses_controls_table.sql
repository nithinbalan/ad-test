-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_library_clauses_controls (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_library_clause_id INT NOT NULL,
    control_library_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    FOREIGN KEY (compliance_library_clause_id) REFERENCES compliance_library_clauses(id),
    FOREIGN KEY (control_library_id) REFERENCES control_library(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_library_clauses_controls;
-- +goose StatementEnd


