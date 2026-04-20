-- +goose Up
-- +goose StatementBegin

ALTER TABLE controls
ADD 
    control_status_id INT NULL;

    -- Foreign Keys
ALTER TABLE controls
ADD CONSTRAINT FK_controls_control_statuses
FOREIGN KEY (control_status_id)
REFERENCES control_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop constraints first
ALTER TABLE controls DROP CONSTRAINT FK_controls_control_statuses;

-- Then drop columns
ALTER TABLE controls
DROP COLUMN 
    control_status_id;

-- +goose StatementEnd
