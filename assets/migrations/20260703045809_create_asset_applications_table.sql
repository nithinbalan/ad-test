-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_applications (
    id INT PRIMARY KEY IDENTITY(1,1),

    asset_register_id INT NOT NULL,
    application_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_asset_applications_asset_register FOREIGN KEY (asset_register_id)
      REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_applications_application FOREIGN KEY (application_id)
      REFERENCES applications(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_applications_user_1 FOREIGN KEY (created_by_id)
      REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_applications;
-- +goose StatementEnd