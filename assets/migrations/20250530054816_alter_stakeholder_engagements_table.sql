-- +goose Up
-- +goose StatementBegin
ALTER TABLE stakeholder_engagements
DROP CONSTRAINT FK_stakeholder_engagements_review_frequencies;

ALTER TABLE stakeholder_engagements
ADD CONSTRAINT FK_stakeholder_engagements_review_frequencies
FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE stakeholder_engagements
DROP CONSTRAINT FK_stakeholder_engagements_review_frequencies;

ALTER TABLE stakeholder_engagements
ADD CONSTRAINT FK_stakeholder_engagements_review_frequencies
FOREIGN KEY (review_frequency_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
