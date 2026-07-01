-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_statuses
DROP CONSTRAINT CHK_asset_statuses_type;

ALTER TABLE asset_statuses
ADD CONSTRAINT CHK_asset_statuses_type
CHECK (
    type IN (
        'active',
        'available',
        'maintenance',
        'disposed',
        'discovered',
        'unavailable',
        'decommissioned',
        'missing',
        'under-review'
    )
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_statuses
DROP CONSTRAINT CHK_asset_statuses_type;

ALTER TABLE asset_statuses
ADD CONSTRAINT CHK_asset_statuses_type
CHECK (
    type IN (
        'active',
        'available',
        'maintenance',
        'disposed'
    )
);
-- +goose StatementEnd