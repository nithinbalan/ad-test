-- +goose Up
-- +goose StatementBegin
CREATE TABLE tprm_settings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    is_export BIT NOT NULL DEFAULT 1,
    is_external_vendor BIT NOT NULL DEFAULT 0,
    updated_at DATETIME,
    updated_by_id INT NULL,
    CONSTRAINT FK_tprm_settings_users_1 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.tprm_settings;
-- +goose StatementEnd