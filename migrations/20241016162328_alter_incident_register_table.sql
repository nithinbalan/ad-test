-- +goose Up
-- +goose StatementBegin

ALTER TABLE incident_register
    ALTER COLUMN incident_at DATETIME NOT NULL;

ALTER TABLE incident_register
    ALTER COLUMN reported_at DATETIME NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE incident_register
    ALTER COLUMN incident_at DATE NOT NULL;

ALTER TABLE incident_register
    ALTER COLUMN reported_at DATE NOT NULL;

-- +goose StatementEnd
