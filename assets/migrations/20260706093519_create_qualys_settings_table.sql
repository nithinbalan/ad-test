-- +goose Up
-- +goose StatementBegin
CREATE TABLE qualys_settings (
    id INT IDENTITY(1,1) PRIMARY KEY,

    platform_url VARCHAR(255) NOT NULL,
    username VARCHAR(300) NOT NULL,
    password VARCHAR(300) NOT NULL,
    lookback_days INT NOT NULL DEFAULT 30,
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_qualys_settings_created_by FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_qualys_settings_updated_by FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE qualys_settings;
-- +goose StatementEnd