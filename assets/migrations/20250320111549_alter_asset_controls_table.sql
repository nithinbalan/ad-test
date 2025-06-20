-- +goose Up
-- +goose StatementBegin
ALTER TABLE asset_controls
ADD 
    control_implementation_status_id INT NULL,
    control_efficiency_status_id INT NULL,
    comments NVARCHAR(MAX) NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_asset_controls_control_implementation_statuses FOREIGN KEY (control_implementation_status_id) REFERENCES control_implementation_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_controls_control_efficiency_statuses FOREIGN KEY (control_efficiency_status_id) REFERENCES control_efficiency_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_asset_controls_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE asset_controls 
DROP CONSTRAINT FK_asset_controls_control_implementation_statuses;

ALTER TABLE asset_controls 
DROP CONSTRAINT FK_asset_controls_control_efficiency_statuses;

ALTER TABLE asset_controls 
DROP CONSTRAINT FK_asset_controls_users_2;

ALTER TABLE asset_controls 
DROP COLUMN control_implementation_status_id,
DROP COLUMN control_efficiency_status_id,
DROP COLUMN comments,
DROP COLUMN updated_by_id,
DROP COLUMN updated_at;
-- +goose StatementEnd
