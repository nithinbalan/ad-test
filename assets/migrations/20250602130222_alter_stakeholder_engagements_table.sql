-- +goose Up
-- +goose StatementBegin
ALTER TABLE stakeholder_engagements
ADD evaluation_status_id INT NULL;

ALTER TABLE stakeholder_engagements
ADD CONSTRAINT FK_stakeholder_engagements_evaluation_statuses
FOREIGN KEY (evaluation_status_id) REFERENCES evaluation_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stakeholder_engagements
DROP CONSTRAINT FK_stakeholder_engagements_evaluation_statuses;

ALTER TABLE stakeholder_engagements
DROP COLUMN evaluation_status_id;
-- +goose StatementEnd
