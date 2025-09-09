-- +goose Up
-- +goose StatementBegin
 
INSERT INTO rag_settings (threshold_days, on_track_color_code, approaching_color_code, overdue_color_code)
VALUES 
(14, '#28A745', '#FFC107', '#DC3545');  
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE rag_settings;
-- +goose StatementEnd