-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'swot_question_groups', 'swot_pestel_question_groups';

ALTER TABLE swot_pestel_question_groups
ADD is_swot TINYINT NULL;

ALTER TABLE swot_pestel_question_groups
ADD is_pestel TINYINT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE swot_pestel_question_groups
DROP COLUMN is_swot;

ALTER TABLE swot_pestel_question_groups
DROP COLUMN is_pestel;

EXEC sp_rename 'swot_pestel_question_groups', 'swot_question_groups';
-- +goose StatementEnd
