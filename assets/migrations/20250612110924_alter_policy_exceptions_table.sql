-- +goose Up
-- +goose StatementBegin

ALTER TABLE policy_exceptions
DROP COLUMN title_en,
            title_ar;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE policy_exceptions
ADD title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL;

-- +goose StatementEnd
