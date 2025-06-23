-- +goose Up
-- +goose StatementBegin

INSERT INTO compliance_libraries (
    icon,
    title_en,
    title_ar,
    description_en,
    description_ar,
    issued_by,
    classification,
    document_classification,
    disclaimer,
    created_at,
    updated_at
)
VALUES (
    'sama-icon',
    'Cyber Security Framework | Saudi Arabian Monetary Authority',
    'إطار عمل الأمن السيبراني | مؤسسة النقد العربي السعودي',
    NULL,
    NULL,
    'Saudi Arabian Monetary Authority (SAMA)',
    NULL,
    NULL,
    NULL,
    GETDATE(),
    GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM compliance_libraries 
WHERE icon = 'sama-icon' 
  AND title_en = 'Cyber Security Framework | Saudi Arabian Monetary Authority';

-- +goose StatementEnd
