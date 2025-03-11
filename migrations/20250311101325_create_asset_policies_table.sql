-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_policies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_register_id INT NOT NULL,
    document_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_asset_policies_asset_registers FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION ,
    CONSTRAINT FK_asset_policies_controls FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policies_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policies_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_policies;
-- +goose StatementEnd
