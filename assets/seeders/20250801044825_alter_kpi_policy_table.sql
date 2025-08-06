-- +goose Up
-- +goose StatementBegin

-- Drop existing constraints
ALTER TABLE dbo.kpi_policies DROP CONSTRAINT FK_kpi_policies_kpi_register;
ALTER TABLE dbo.kpi_policies DROP CONSTRAINT FK_kpi_policies_documents;
ALTER TABLE dbo.kpi_policies DROP CONSTRAINT FK_kpi_policies_users_1;

-- Rename table
EXEC sp_rename 'dbo.kpi_policies', 'kpi_associated_documents';

-- Recreate constraints with new names
ALTER TABLE dbo.kpi_associated_documents ADD CONSTRAINT FK_kpi_associated_documents_kpi_register
    FOREIGN KEY (kpi_register_id) REFERENCES kpi_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE dbo.kpi_associated_documents ADD CONSTRAINT FK_kpi_associated_documents_documents
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE dbo.kpi_associated_documents ADD CONSTRAINT FK_kpi_associated_documents_users_1
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Drop new constraints
ALTER TABLE dbo.kpi_associated_documents DROP CONSTRAINT FK_kpi_associated_documents_kpi_register;
ALTER TABLE dbo.kpi_associated_documents DROP CONSTRAINT FK_kpi_associated_documents_documents;
ALTER TABLE dbo.kpi_associated_documents DROP CONSTRAINT FK_kpi_associated_documents_users_1;

-- Rename table back
EXEC sp_rename 'dbo.kpi_associated_documents', 'kpi_policies';

-- Recreate original constraints
ALTER TABLE dbo.kpi_policies ADD CONSTRAINT FK_kpi_policies_kpi_register
    FOREIGN KEY (kpi_register_id) REFERENCES kpi_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE dbo.kpi_policies ADD CONSTRAINT FK_kpi_policies_documents
    FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE dbo.kpi_policies ADD CONSTRAINT FK_kpi_policies_users_1
    FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
