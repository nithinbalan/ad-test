-- +goose Up
-- +goose StatementBegin
CREATE TABLE themes (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX),
    description_ar NVARCHAR(MAX),
    primary_color VARCHAR(20) NOT NULL,
    secondary_color VARCHAR(20) NOT NULL,
    link_color VARCHAR(20) NOT NULL,
    sidebar_bg_color VARCHAR(20) NOT NULL,
    is_system BIT NOT NULL DEFAULT 0,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_themes_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_themes_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE themes;
-- +goose StatementEnd
