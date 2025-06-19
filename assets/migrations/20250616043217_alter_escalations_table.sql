-- +goose Up
-- +goose StatementBegin
ALTER TABLE escalations
ADD description_en NVARCHAR(MAX) NULL;

ALTER TABLE escalations
ADD description_ar NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE escalations
DROP COLUMN description_en;

ALTER TABLE escalations
DROP COLUMN description_ar;
-- +goose StatementEnd
