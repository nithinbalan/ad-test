-- +goose Up
-- +goose StatementBegin
CREATE TABLE branding
(
    id INT PRIMARY KEY IDENTITY(1,1),
    login_page_title_en NVARCHAR(MAX) NOT NULL,
    login_page_title_ar NVARCHAR(MAX) NOT NULL,
    login_page_subtitle_en NVARCHAR(MAX) NOT NULL,
    login_page_subtitle_ar NVARCHAR(MAX) NOT NULL,
    login_page_bg_image_path VARCHAR(500) DEFAULT NULL,
    org_logo_image_path VARCHAR(500) DEFAULT NULL,
    favicon_image_path VARCHAR(500) DEFAULT NULL,
    loader_gif_path VARCHAR(500) DEFAULT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_branding_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_branding_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE branding;
-- +goose StatementEnd

