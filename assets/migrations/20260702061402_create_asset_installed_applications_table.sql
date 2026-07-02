-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_installed_applications
(
    id INT PRIMARY KEY IDENTITY(1,1),

    asset_register_id INT NOT NULL,
    installed_application_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_asset_installed_application_asset FOREIGN KEY (asset_register_id)
        REFERENCES asset_register(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_installed_application_application FOREIGN KEY (installed_application_id)
        REFERENCES installed_applications(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_installed_application_user_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_installed_applications;
-- +goose StatementEnd