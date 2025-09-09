-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answers
ADD comment NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_assessment_answers
DROP COLUMN comment;
-- +goose StatementEnd
