-- +goose Up
-- +goose StatementBegin
UPDATE kpi_statuses
SET title_en = 'Achieved',
    title_ar = N'تم التحقيق',
    color_code = '#28A745',
    type = 'achieved',
    status = 'active',
    label = 'Required Attention',
    is_not_achieved = 0,
    [order] = 2
WHERE type = 'active';

INSERT INTO kpi_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    status,
    label,
    is_not_achieved,
    [order]
)
VALUES
    ('Low', N'منخفض', '#DC3545', 'low', 'active', 'Required Attention', 1, 5),
    ('Needs Improvement', N'يحتاج تحسين', '#FD7E14', 'needs-improvement', 'active', 'Below Target', 1, 4),
    ('Near Target', N'قريب من الهدف', '#D2B48C', 'near-target', 'active', 'Approaching Goal', 0, 3),
    ('Exceeded', N'تجاوز الهدف', '#007BFF', 'exceeded', 'active', 'Outstanding', 0, 1),
    ('Pending', N'قيد الانتظار', '#808080', 'new', 'active', NULL, 0, 6);

DECLARE @lowStatusID INT;

SELECT @lowStatusID = id
FROM kpi_statuses
WHERE type = 'low';

UPDATE kpi_register
SET kpi_status_id = @lowStatusID
WHERE kpi_status_id IN (
    SELECT id
    FROM kpi_statuses
    WHERE type IN ('inactive', 'archived')
);

DELETE FROM kpi_statuses
WHERE type IN ('inactive', 'archived');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO kpi_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    status
)
VALUES
    ('Inactive', N'غير نشط', '#FFA500', 'inactive', 'active'),
    ('Archived', N'مؤرشف', '#DC3545', 'archived', 'active');

DECLARE @activeStatusID INT;
DECLARE @inactiveStatusID INT;
DECLARE @archivedStatusID INT;

SELECT @activeStatusID = id
FROM kpi_statuses
WHERE type = 'achieved';

SELECT @inactiveStatusID = id
FROM kpi_statuses
WHERE type = 'inactive';

SELECT @archivedStatusID = id
FROM kpi_statuses
WHERE type = 'archived';

UPDATE kr
SET kpi_status_id = CASE
    WHEN ks.type IN ('low', 'needs-improvement') THEN @archivedStatusID
    WHEN ks.type = 'new' THEN @inactiveStatusID
    ELSE @activeStatusID
END
FROM kpi_register kr
INNER JOIN kpi_statuses ks
    ON ks.id = kr.kpi_status_id
WHERE ks.type IN (
    'low',
    'needs-improvement',
    'near-target',
    'exceeded',
    'new'
);

DELETE FROM kpi_statuses
WHERE type IN (
    'low',
    'needs-improvement',
    'near-target',
    'exceeded',
    'new'
);

UPDATE kpi_statuses
SET title_en = 'Active',
    title_ar = N'نشيط',
    color_code = '#28A745',
    type = 'active',
    status = 'active',
    label = NULL,
    is_not_achieved = 0,
    [order] = NULL
WHERE type = 'achieved';
-- +goose StatementEnd
