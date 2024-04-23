-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_clauses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    clause_no VARCHAR(45) NOT NULL,
    title NVARCHAR(MAX) NOT NULL,
    CONSTRAINT FK_compliance_register_clauses_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_clauses;
-- +goose StatementEnd
