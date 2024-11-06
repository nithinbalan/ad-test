-- +goose Up
-- +goose StatementBegin
CREATE TABLE integrations (
    id INT PRIMARY KEY IDENTITY(1, 1),
    integration_category_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    logo_url NVARCHAR(MAX) NULL,
    api_base_url NVARCHAR(MAX) NULL,
    is_enable TINYINT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL ,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_integrations_integration_categories FOREIGN KEY (integration_category_id) REFERENCES integrations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_integrations_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_integrations_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS integrations;
-- +goose StatementEnd
