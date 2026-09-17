-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_closure_lessons ALTER COLUMN lessons NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_closure_lessons ALTER COLUMN lessons NVARCHAR(500) NOT NULL;
-- +goose StatementEnd
