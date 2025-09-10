-- +goose Up
-- +goose StatementBegin

DECLARE @documentModuleID INT;
SELECT @documentModuleID = id FROM modules WHERE identifier = 'document';

INSERT INTO reminder_contexts (type, title_en, title_ar, module_id)
VALUES (
    'document-review-due',
    'Document Review Date Due',
    'تاريخ استحقاق مراجعة الوثيقة',
    @documentModuleID
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM reminder_contexts
WHERE type = 'document-review-due';

-- +goose StatementEnd
