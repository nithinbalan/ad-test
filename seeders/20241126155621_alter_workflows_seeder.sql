-- +goose Up
-- +goose StatementBegin
INSERT INTO workflows (document_type_id, type, module_identifier, status, created_by_id, created_at)
SELECT dt.id, 'document', 'document', 'active', 1, GETDATE()
FROM document_types dt
WHERE dt.type IN ('policy', 'procedure', 'others')
AND NOT EXISTS (
    SELECT 1
    FROM workflows wf
    WHERE wf.document_type_id = dt.id
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the inserted workflows for the specific document types
DELETE FROM workflows
WHERE document_type_id IN (
    SELECT id
    FROM document_types
    WHERE type IN ('policy', 'procedure', 'others')
);
-- +goose StatementEnd
