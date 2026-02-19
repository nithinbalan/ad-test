-- +goose Up
-- +goose StatementBegin

ALTER TABLE control_consistencies
ADD
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE control_consistencies
DROP COLUMN description_en,
            description_ar;

-- +goose StatementEnd
