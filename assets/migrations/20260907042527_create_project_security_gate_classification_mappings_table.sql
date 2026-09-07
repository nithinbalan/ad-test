-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_classification_mappings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    security_gate_id INT NOT NULL,
    classification_id INT NOT NULL,

    CONSTRAINT UQ_project_security_gate_classification_mappings
        UNIQUE (security_gate_id, classification_id),

    CONSTRAINT FK_project_security_gate_classification_mappings_gate
        FOREIGN KEY (security_gate_id)
        REFERENCES project_security_gates(id),

    CONSTRAINT FK_project_security_gate_classification_mappings_classification
        FOREIGN KEY (classification_id)
        REFERENCES project_security_gate_classifications(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_classification_mappings;
-- +goose StatementEnd