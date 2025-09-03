-- +goose Up
-- +goose StatementBegin
CREATE TABLE  meeting_assets (
    id INT PRIMARY KEY IDENTITY(1, 1),
	meeting_id INT NOT NULL,
    asset_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_meeting_assets_meetings FOREIGN KEY (meeting_id) REFERENCES meetings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_assets_asset_register FOREIGN KEY (asset_register_id) REFERENCES asset_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_meeting_assets_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE  meeting_assets;

-- +goose StatementEnd
