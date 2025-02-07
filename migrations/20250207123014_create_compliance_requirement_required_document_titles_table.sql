-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_required_document_titles (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_requirement_id INT NOT NULL,
    required_document_title_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_requirement_required_document_titles_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_required_document_titles_required_document FOREIGN KEY (required_document_title_id) REFERENCES required_document_titles(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirement_required_document_titles;
-- +goose StatementEnd
