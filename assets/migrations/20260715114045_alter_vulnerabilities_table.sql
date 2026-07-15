-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'vulnerabilities.deffered_by_id', 'deferred_by_id', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'vulnerabilities.deffered_at', 'deferred_at', 'COLUMN';
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
EXEC sp_rename 'vulnerabilities.deferred_at', 'deffered_at', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'vulnerabilities.deferred_by_id', 'deffered_by_id', 'COLUMN';
-- +goose StatementEnd