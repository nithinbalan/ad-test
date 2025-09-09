-- +goose Up
-- +goose StatementBegin
UPDATE knowledge_center_resources
SET knowledge_center_resource_status_id = NULL
WHERE title_en IN (
  'Saudi Arabian Monetary Authority (SAMA)',
  'National Cybersecurity Authority (NCA)'
);

DECLARE @complianceID INT;
DECLARE @itemCount INT;

SELECT @complianceID = id FROM modules WHERE identifier = 'compliance';
SELECT @itemCount = COUNT(*) FROM control_library;

INSERT INTO knowledge_center_resources (
    module_id,
    knowledge_center_resource_status_id,
    title_en,
    title_ar,
    icon,
    created_at,
    updated_at,
    type,
    item_count
)
VALUES (
    @complianceID,
    NULL,
    'Control Mapping Library',
    'مكتبة خرائط التحكم',
    'control-mapping-library-icon',
    GETDATE(),
    GETDATE(),
    'control',
    @itemCount
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @statusID INT;
SELECT @statusID = id FROM knowledge_center_resource_statuses WHERE title_en = 'Coming Soon';

UPDATE knowledge_center_resources
SET knowledge_center_resource_status_id = @statusID
WHERE title_en = 'Saudi Arabian Monetary Authority (SAMA)';

UPDATE knowledge_center_resources
SET knowledge_center_resource_status_id = @statusID
WHERE title_en = 'National Cybersecurity Authority (NCA)';

DELETE FROM knowledge_center_resources
WHERE title_en = 'Control Mapping Library'
  AND type = 'control';
-- +goose StatementEnd
