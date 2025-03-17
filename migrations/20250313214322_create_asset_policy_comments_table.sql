-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_policy_comments (
    id INT PRIMARY KEY IDENTITY(1, 1),
    asset_policy_id INT NOT NULL,
    comments NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_asset_policy_comments_asset_policy FOREIGN KEY (asset_policy_id) REFERENCES asset_policies(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_comments_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_policy_comments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE asset_policy_comments;
-- +goose StatementEnd
