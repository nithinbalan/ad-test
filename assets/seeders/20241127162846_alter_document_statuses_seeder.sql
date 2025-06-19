-- +goose Up
-- +goose StatementBegin
-- Insert the new "Reverted" status
INSERT INTO document_statuses (title_en, title_ar, color_code, type, status) VALUES
('Reverted', 'معاد', '#28A745', 'reverted', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the "Reverted" status
DELETE FROM document_statuses
WHERE type = 'reverted';
-- +goose StatementEnd
