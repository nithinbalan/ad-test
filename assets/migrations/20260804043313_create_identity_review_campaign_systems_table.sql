-- +goose Up
-- +goose StatementBegin
CREATE TABLE identity_review_campaign_systems (
    id INT PRIMARY KEY IDENTITY(1, 1),
    identity_review_campaign_id INT NOT NULL,
    system_register_id INT NOT NULL,
    asset_register_id INT NULL,
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,

    CONSTRAINT FK_identity_review_campaign_systems_identity_review_campaigns
        FOREIGN KEY (identity_review_campaign_id)
        REFERENCES identity_review_campaigns(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_systems_system_registers
        FOREIGN KEY (system_register_id)
        REFERENCES system_registers(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_systems_asset_register
        FOREIGN KEY (asset_register_id)
        REFERENCES asset_register(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_identity_review_campaign_systems_users
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE identity_review_campaign_systems;
-- +goose StatementEnd
