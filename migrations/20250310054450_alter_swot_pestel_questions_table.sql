-- +goose Up
-- +goose StatementBegin

ALTER TABLE swot_pestel_questions
ADD answer NVARCHAR(MAX) NULL;

ALTER TABLE swot_pestel_questions
ADD answered_by_id INT NULL;

ALTER TABLE swot_pestel_questions
ADD answered_at DATETIME NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE swot_pestel_questions DROP COLUMN answer;
ALTER TABLE swot_pestel_questions DROP COLUMN answered_by_id;
ALTER TABLE swot_pestel_questions DROP COLUMN answered_at;
-- +goose StatementEnd
