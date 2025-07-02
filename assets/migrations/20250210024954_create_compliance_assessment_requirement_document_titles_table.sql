-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirement_document_titles (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_assessment_requirement_id INT NOT NULL,
    required_document_title_id INT NULL,
    applicability_status_id INT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_compliance_assessment_requirement_document_titles_compliance_assessment_requirements FOREIGN KEY (compliance_assessment_requirement_id) REFERENCES compliance_assessment_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_document_titles_required_document_titles FOREIGN KEY (required_document_title_id) REFERENCES required_document_titles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_document_titles_applicability_statuses FOREIGN KEY (applicability_status_id) REFERENCES applicability_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_document_titles_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_assessment_requirement_document_titles;
-- +goose StatementEnd
