-- +goose Up
-- +goose StatementBegin

DECLARE @ITSupportOrganizationUnitID INT;

SELECT @ITSupportOrganizationUnitID = id
FROM organization_units
WHERE title_en = 'IT Support';

INSERT INTO governance_settings (
    approve_department_id,
    created_by_id,
    created_at,
    updated_by_id,
    updated_at
)
VALUES (
    @ITSupportOrganizationUnitID,
    1,
    GETDATE(),
    NULL,
    NULL
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

TRUNCATE TABLE governance_settings;

-- +goose StatementEnd