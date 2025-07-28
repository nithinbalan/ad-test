-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'kpi_policies', 'kpi_associated_documents';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
EXEC sp_rename 'kpi_associated_documents', 'kpi_policies';
-- +goose StatementEnd
