-- +goose Up
-- +goose StatementBegin
ALTER TABLE third_party_assessments 
ADD next_review_level INT NULL,
submitted_by_id INT NULL,
submitted_at DATETIME NULL,
CONSTRAINT FK_third_party_assessments_users_3 FOREIGN KEY (submitted_by_id) 
REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_assessments
DROP CONSTRAINT FK_third_party_assessments_users_3;

ALTER TABLE third_party_assessments
DROP COLUMN submitted_by_id;

ALTER TABLE third_party_assessments
DROP COLUMN next_review_level;

ALTER TABLE third_party_assessments
DROP COLUMN submitted_at;
-- +goose StatementEnd
