-- +goose Up
-- +goose StatementBegin

ALTER TABLE committee_terms
ADD color_code VARCHAR(20) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE committee_terms
DROP COLUMN color_code;

-- +goose StatementEnd