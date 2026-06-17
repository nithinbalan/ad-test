-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_assessment_requirement_required_documents
ADD 
    approved_by_id INT NULL,
    approved_at DATETIME NULL,
    rejected_by_id INT NULL,
    rejected_at DATETIME NULL;

ALTER TABLE compliance_assessment_requirement_required_documents
ADD CONSTRAINT FK_carrd_approved_by
FOREIGN KEY (approved_by_id)
REFERENCES users(id);

ALTER TABLE compliance_assessment_requirement_required_documents
ADD CONSTRAINT FK_carrd_rejected_by
FOREIGN KEY (rejected_by_id)
REFERENCES users(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_assessment_requirement_required_documents
DROP CONSTRAINT FK_carrd_approved_by;

ALTER TABLE compliance_assessment_requirement_required_documents
DROP CONSTRAINT FK_carrd_rejected_by;

ALTER TABLE compliance_assessment_requirement_required_documents
DROP COLUMN 
    approved_by_id,
    approved_at,
    rejected_by_id,
    rejected_at;

-- +goose StatementEnd