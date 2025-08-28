-- +goose Up
-- +goose StatementBegin
CREATE TABLE   bcp_assets (
    id INT PRIMARY KEY IDENTITY(1, 1),
	business_continuity_plan_id INT NOT NULL,
    asset_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_bcp_assets_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_assets_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_bcp_assets_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE   bcp_assets;

-- +goose StatementEnd
