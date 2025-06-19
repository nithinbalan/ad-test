-- +goose Up
-- +goose StatementBegin

ALTER TABLE swot_pestel_categories DROP CONSTRAINT CHK_swot_pestel_categories_type;

ALTER TABLE swot_pestel_categories ADD CONSTRAINT CHK_swot_pestel_categories_type 
CHECK (type IN ('strength', 'weakness','opportunity','threat','political','economic','social','technological','legal','environmental'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE swot_pestel_categories DROP CONSTRAINT CHK_swot_pestel_categories_type;

ALTER TABLE swot_pestel_categories ADD CONSTRAINT CHK_swot_pestel_categories_type 
CHECK (type IN ('strength', 'weakness','opportunity','threat'));
-- +goose StatementEnd
