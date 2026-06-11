-- +goose Up
-- +goose StatementBegin

IF NOT EXISTS (
    SELECT 1
    FROM roles
    WHERE type = 'committee-requester'
)
BEGIN
    INSERT INTO roles (
        title_en,
        title_ar,
        status,
        is_super_admin,
        license_type_id,
        type
    )
    VALUES (
        'Committee Requester',
        N'طالب اللجنة',
        'active',
        0,
        1,
        'committee-requester'
    );
END

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM roles
WHERE type = 'committee-requester';

-- +goose StatementEnd--