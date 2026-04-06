-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_initiative_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    strategy_initiative_id INT NOT NULL,
    strategic_initiative_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    percentage INT NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_strategy_initiative_status_updates_strategy_initiatives FOREIGN KEY (strategy_initiative_id) REFERENCES strategic_intiatives (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_initiative_status_updates_strategic_initiative_statuses FOREIGN KEY (strategic_initiative_status_id) REFERENCES strategic_initiative_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_initiative_status_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_initiative_status_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE strategic_initiative_status_updates;
-- +goose StatementEnd
