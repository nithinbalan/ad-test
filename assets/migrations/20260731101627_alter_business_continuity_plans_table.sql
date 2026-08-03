-- +goose Up
-- +goose StatementBegin

ALTER TABLE business_continuity_plans
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_business_continuity_plans_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE business_continuity_plans
DROP CONSTRAINT DF_business_continuity_plans_is_cyber;

ALTER TABLE business_continuity_plans
DROP COLUMN is_cyber;

-- +goose StatementEnd