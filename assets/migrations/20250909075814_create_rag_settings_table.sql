-- +goose Up
-- +goose StatementBegin
CREATE TABLE rag_settings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    threshold_days INT NOT NULL,
    on_track_color_code VARCHAR(25) NOT NULL,
    approaching_color_code VARCHAR(25) NOT NULL,
    overdue_color_code VARCHAR(25) NOT NULL,
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE rag_settings;
-- +goose StatementEnd