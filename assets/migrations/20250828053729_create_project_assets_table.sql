-- +goose Up
-- +goose StatementBegin
CREATE TABLE  project_assets (
    id INT PRIMARY KEY IDENTITY(1, 1),
	project_id INT NOT NULL,
    asset_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_project_assets_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_assets_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_project_assets_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE  project_assets;

-- +goose StatementEnd
