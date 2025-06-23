-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_kris (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_register_id INT NOT NULL,
    kri_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_asset_kris_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_kris_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_kris;

-- +goose StatementEnd
