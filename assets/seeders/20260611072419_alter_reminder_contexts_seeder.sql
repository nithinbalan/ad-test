-- +goose Up
-- +goose StatementBegin

DECLARE @OrganizationModuleID INT;

SELECT @OrganizationModuleID = id
FROM modules
WHERE identifier = 'organization';


    INSERT INTO reminder_contexts (
        type,
        title_en,
        title_ar,
        status,
        module_id
    )
    VALUES (
        'committee-management-due',
        'Committee Management Due',
        N'استحقاق إدارة اللجان',
        'active',
        @OrganizationModuleID
    );


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM reminder_contexts
WHERE type = 'committee-management-due';

-- +goose StatementEnd